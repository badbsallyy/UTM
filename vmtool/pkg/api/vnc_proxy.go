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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleVNCProxy(c *gin.Context) {
	name := c.Param("name")
	port := s.manager.GetVNCPort(name)
	if port == 0 {
		port = 5900 // Fallback
	}
	vncAddr := fmt.Sprintf("localhost:%d", port)

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
