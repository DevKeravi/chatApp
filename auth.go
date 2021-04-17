package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

func loginHandler(c *gin.Context) {
	action := c.Param("action")
	log.Println(action)
	provider := c.Param("provider")
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error when trying to get provider"})
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when trying to GetBeginAuthURL"})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, loginUrl)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error when trying to get provider"})
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(c.Request.URL.RawQuery))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error when trying complete Auth"})
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error when to get user"})
			return
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()

		unescapeCookie, err := url.QueryUnescape(authCookieValue)
		if err != nil {
			log.Println("QueryUnescape is failed:", err)
		}

		c.SetCookie("auth", unescapeCookie, 3600, "/", "", false, false)
		c.Redirect(http.StatusTemporaryRedirect, "/api/chat.html")
		return

	default:
		c.JSON(http.StatusNotFound, gin.H{"msg": "Auth action not supported"})
		return
	}
}
