package handlers

import (
	"net/http"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/gin-gonic/gin"
)

// ReportHandler handles all requests on /report path
func ReportHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	err := models.T.ExecuteTemplate(c.Writer, "get_report.template", models.D)
	if err != nil {
		err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
		if err != nil {
			panic(err.Error())
		}
		return
	}
}
