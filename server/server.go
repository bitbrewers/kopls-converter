package server

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type loginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// Server struct wraps connection and listening address
type Server struct {
	storage     *Client
	addr        string
	admin       string
	adminPasswd string
	domain      string
}

// NewServer sets values for Server struct and return it
func NewServer(addr string, storage *Client, admin, adminPasswd, domain string) (server *Server) {
	server = &Server{
		storage:     storage,
		addr:        addr,
		admin:       admin,
		adminPasswd: adminPasswd,
		domain:      domain,
	}

	return
}

// Start will start http server
func (s *Server) Start() error {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"byteToString": func(b byte) string { return string(b) },
	})
	router.LoadHTMLGlob("templates/*")

	router.GET("/login", s.renderLogin)
	router.GET("/logout", s.handleLogout)
	router.POST("/login", s.handleLogin)

	authorized := router.Group("/", s.auth)
	{
		authorized.GET("/", s.renderIndex)
		authorized.GET("/variables", s.renderVariables)
		authorized.POST("/convert", s.handleConversion)

		authorized.PUT("/variables/programs", s.setPrograms)
		authorized.PUT("/variables/doormodels", s.setDoorModels)
		authorized.PUT("/variables/hinges", s.setHinges)
		authorized.PUT("/variables/handednesses", s.setHandednesses)
		authorized.PUT("/variables/handles", s.setHandles)
		authorized.PUT("/variables/handlepositions", s.setHandlePositions)

		authorized.POST("/variables/programs", s.addPrograms)
		authorized.POST("/variables/doormodels", s.addDoorModels)
		authorized.POST("/variables/hinges", s.addHinges)
		authorized.POST("/variables/handednesses", s.addHandednesses)
		authorized.POST("/variables/handles", s.addHandles)
		authorized.POST("/variables/handlepositions", s.addHandlePositions)
	}

	return router.Run(s.addr)
}

func (s *Server) renderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func (s *Server) renderVariables(c *gin.Context) {
	cData, err := s.storage.GetAll()
	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "variables.tmpl", cData)
}

func (s *Server) renderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func (s *Server) handleConversion(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		return
	}
	result, err := Convert(file, s.storage)
	if err != nil {
		return
	}

	name := strings.Replace(fileHeader.Filename, ".csv", ".zip", 1)
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Header("Content-Type", "text/plain")

	if _, err := c.Writer.Write(result.Bytes()); err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) setPrograms(c *gin.Context) {
	data := make([]Program, 0)
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateProgram(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) setDoorModels(c *gin.Context) {
	models := make([]DoorModel, 0)
	if err := c.ShouldBindJSON(&models); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateDoorModels(models); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) setHinges(c *gin.Context) {
	data := make([]Hinge, 0)
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateHinges(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) setHandednesses(c *gin.Context) {
	data := make([]Handedness, 0)
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateHandednesses(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) setHandles(c *gin.Context) {
	data := make([]Handle, 0)
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateHandles(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) setHandlePositions(c *gin.Context) {
	data := make([]HandlePosition, 0)
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.UpdateHandlePositions(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addPrograms(c *gin.Context) {
	data := &Program{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddProgram(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addDoorModels(c *gin.Context) {
	data := &DoorModel{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddDoorModels(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addHinges(c *gin.Context) {
	data := &Hinge{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddHinges(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addHandednesses(c *gin.Context) {
	data := &Handedness{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddHandednesses(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addHandles(c *gin.Context) {
	data := &Handle{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddHandles(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) addHandlePositions(c *gin.Context) {
	data := &HandlePosition{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.storage.AddHandlePositions(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// TODO
func (s *Server) handleLogin(c *gin.Context) {
	loginForm := &loginForm{}
	if err := c.ShouldBind(loginForm); err != nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", nil)
		return
	}

	log.Printf("%+v\n", *loginForm)
	if loginForm.User != s.admin || loginForm.Password != s.adminPasswd {
		c.HTML(http.StatusUnauthorized, "login.tmpl", nil)
		return
	}

	c.SetCookie("token", "some-token", 60*60*24*7, "/", s.domain, false, false)
	c.Redirect(http.StatusFound, "/")
}

// TODO
func (s *Server) handleLogout(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*7, "/", s.domain, false, false)
	c.Redirect(http.StatusFound, "/")
}

// TODO
func (s *Server) auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		log.Println("err", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	if token != "some-token" {
		log.Println("invalid token", token)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	c.Next()
}
