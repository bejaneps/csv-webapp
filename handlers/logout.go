package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogoutHandler handles upcoming requests on /logout path
func LogoutHandler(c *gin.Context) {
	c.SetCookie("auth", "", 0, "/", "127.0.0.1", false, false)
	c.SetCookie("auth", "", 0, "/dashboard", "127.0.0.1", false, false)
	c.SetCookie("auth", "", 0, "/report", "127.0.0.1", false, false)

	c.Redirect(http.StatusTemporaryRedirect, "/")
	return
}
