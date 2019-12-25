package handlers

import (
	"net/http"

	"github.com/bejaneps/csv-webapp/crud"
	"github.com/bejaneps/csv-webapp/models"

	"github.com/gin-gonic/gin"
)

const templatesFolder = "templates/"

// IndexHandler handles incoming GET requests on a web-app / path
func IndexHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val != "" {
		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

// LoginHandler handles incoming GET requests on a web-app /login path
func LoginHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val != "" {
		c.Redirect(http.StatusOK, "/dashboard")
		return
	}

	info := models.LoginInfo{
		Email:    c.Query("email"),
		Password: c.Query("password"),
	}

	if info.Email == "" || info.Password == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if ok, err := crud.CheckLoginInfo(info); !ok || err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", err.Error())
	}

	c.SetCookie("auth", "yes", 86400, "/", "127.0.0.1", false, false)
	c.SetCookie("auth", "yes", 86400, "/dashboard", "127.0.0.1", false, false)
	c.SetCookie("auth", "yes", 86400, "/data", "127.0.0.1", false, false)
	c.SetCookie("auth", "yes", 86400, "/report", "127.0.0.1", false, false)
	c.SetCookie("auth", "yes", 86400, "/config", "127.0.0.1", false, false)
	c.SetCookie("auth", "yes", 86400, "/config/submit", "127.0.0.1", false, false)

	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
	return
}
