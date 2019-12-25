package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles upcoming GET requests on a web-app /dashboard path
func DashboardHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", nil)
}
