package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
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
func indexHandler(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/api/chat.html")
	c.Abort()

}

func main() {
	room := newRoom()

	r := gin.Default()
	r.LoadHTMLGlob("./public/*")

	gomniauth.SetSecurityKey(os.Getenv("GOOGLE_SECURITY_KEY"))
	gomniauth.WithProviders(

		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"),
			"https://cmkrosp.iptime.org:8080/auth/callback/google"),
	)

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
	r.RunTLS(":8080", "./keys/server.crt", "./keys/server.key")
}
