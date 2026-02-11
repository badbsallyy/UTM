package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (s *Server) handleVNCProxy(c *gin.Context) {
	name := c.Param("name")
	
	// Check authentication if token is configured
	// Since WebSocket upgrades from browser can't send custom headers,
	// we need to check the token from query string for this endpoint
	if s.config.Security.APIToken != "" {
		token := c.Query("token")
		if token != s.config.Security.APIToken {
			log.Printf("Unauthorized VNC connection attempt for VM %s", name)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}
	
	port := s.manager.GetVNCPort(name)
	if port == 0 {
		log.Printf("No VNC port configured for VM %s", name)
		c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("no VNC port configured for VM %s", name),
		})
		return
	}
	vncAddr := fmt.Sprintf("localhost:%d", port)

	// Create upgrader with origin check specific to this server
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			// Allow same-origin requests and localhost origins
			if origin == "" {
				return true // Non-browser clients
			}
			// Check if origin matches the configured server address (HTTP or HTTPS)
			expectedOriginHTTP := fmt.Sprintf("http://%s:%d", s.config.Server.Host, s.config.Server.Port)
			expectedOriginHTTPS := fmt.Sprintf("https://%s:%d", s.config.Server.Host, s.config.Server.Port)
			localhostOriginHTTP := fmt.Sprintf("http://localhost:%d", s.config.Server.Port)
			localhostOriginHTTPS := fmt.Sprintf("https://localhost:%d", s.config.Server.Port)
			loopbackOriginHTTP := fmt.Sprintf("http://127.0.0.1:%d", s.config.Server.Port)
			loopbackOriginHTTPS := fmt.Sprintf("https://127.0.0.1:%d", s.config.Server.Port)
			return origin == expectedOriginHTTP || origin == expectedOriginHTTPS ||
				origin == localhostOriginHTTP || origin == localhostOriginHTTPS ||
				origin == loopbackOriginHTTP || origin == loopbackOriginHTTPS
		},
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer ws.Close()

	conn, err := net.DialTimeout("tcp", vncAddr, 2*time.Second)
	if err != nil {
		log.Printf("Failed to connect to VNC at %s: %v", vncAddr, err)
		return
	}
	defer conn.Close()

	errChan := make(chan error, 2)

	// WS -> TCP
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			_, err = conn.Write(msg)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	// TCP -> WS
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}
			err = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	<-errChan
}
