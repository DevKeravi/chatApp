package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func wrapHandler(h func(context *gin.Context)) gin.HandlerFunc {

	return func(c *gin.Context) {
		_, err := c.Cookie("auth")
		if err == http.ErrNoCookie {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		h(c)
	}
}

func loginHandler(c *gin.Context) {
	action := c.Param("action")
	log.Println(action)
	provider := c.Param("provider")
	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		c.JSON(http.StatusNotFound, gin.H{"msg": "Auth action not supported"})
	}
}

func indexHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/api/chat.html")
	c.Abort()

}

func main() {
	room := newRoom()

	r := gin.Default()
	r.LoadHTMLGlob("./public/*")

	r.Use(static.Serve("/api", static.LocalFile("./public", true)))

	r.Group("/api")
	{
		r.GET("/", wrapHandler(indexHandler))
	}

	r.GET("/room", func(c *gin.Context) {
		room.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/auth/:action/:provider", loginHandler)

	go room.run()
	r.Run(":8080")
}
