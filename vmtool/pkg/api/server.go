package api

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utmapp/vmtool/pkg/config"
	"github.com/utmapp/vmtool/pkg/vm"
	"github.com/utmapp/vmtool/pkg/web"
)

type Server struct {
	config  *config.AppConfig
	manager *vm.Manager
	router  *gin.Engine
}

func NewServer(manager *vm.Manager, cfg *config.AppConfig) *Server {
	router := gin.Default()
	s := &Server{
		config:  cfg,
		manager: manager,
		router:  router,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Public routes
	s.router.GET("/ui/*any", func(c *gin.Context) {
		staticFS, _ := fs.Sub(web.StaticFiles, "static")
		http.StripPrefix("/ui", http.FileServer(http.FS(staticFS))).ServeHTTP(c.Writer, c.Request)
	})
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ui/index.html")
	})

	// Protected routes
	protected := s.router.Group("/")
	protected.Use(s.authMiddleware())

	protected.GET("/vms", s.handleListVMs)
	protected.POST("/vms/:name/start", s.handleStartVM)
	protected.POST("/vms/:name/stop", s.handleStopVM)
	protected.POST("/vms/:name/pause", s.handlePauseVM)
	protected.POST("/vms/:name/resume", s.handleResumeVM)
	protected.GET("/vms/:name/status", s.handleStatusVM)
	protected.POST("/vms/:name/snapshot/create", s.handleCreateSnapshot)
	
	// VNC WebSocket endpoint handles auth internally (since WebSocket can't use headers)
	s.router.GET("/vms/:name/vnc", s.handleVNCProxy)
}

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No token required if not set in config
		if s.config.Security.APIToken == "" {
			c.Next()
			return
		}

		// Get token from header only
		token := c.GetHeader("X-VMTool-Token")

		if token != s.config.Security.APIToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) handleListVMs(c *gin.Context) {
	vms := s.manager.ListVMs()
	var resp []gin.H
	for _, v := range vms {
		resp = append(resp, gin.H{
			"name":   v.Name,
			"status": s.manager.GetStatus(v.Name),
			"arch":   v.System.Architecture,
		})
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) handleStartVM(c *gin.Context) {
	name := c.Param("name")
	if err := s.manager.StartVM(c.Request.Context(), name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "started"})
}

func (s *Server) handleStopVM(c *gin.Context) {
	name := c.Param("name")
	if err := s.manager.StopVM(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

func (s *Server) handlePauseVM(c *gin.Context) {
	name := c.Param("name")
	if err := s.manager.PauseVM(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "paused"})
}

func (s *Server) handleResumeVM(c *gin.Context) {
	name := c.Param("name")
	if err := s.manager.ResumeVM(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "running"})
}

func (s *Server) handleStatusVM(c *gin.Context) {
	name := c.Param("name")
	status := s.manager.GetStatus(name)
	c.JSON(http.StatusOK, gin.H{"name": name, "status": status})
}

func (s *Server) handleCreateSnapshot(c *gin.Context) {
	vmName := c.Param("name")
	snapName := c.Query("name")
	if snapName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "snapshot name required"})
		return
	}
	if err := s.manager.CreateSnapshot(vmName, snapName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "snapshot created"})
}
