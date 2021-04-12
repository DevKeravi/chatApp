package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/gomniauth"
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
	default:
		c.JSON(http.StatusNotFound, gin.H{"msg": "Auth action not supported"})
	}
}
