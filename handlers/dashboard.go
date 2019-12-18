package handlers

import (
	"net/http"

	"github.com/bejaneps/csv-webapp/crud"
	"github.com/bejaneps/csv-webapp/models"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles upcoming GET requests on a web-app /dashboard path
func DashboardHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	//get data button pressed
	if getData := c.Query("get_data"); getData != "" {
		if getData == "Get Data" {
			err := models.T.ExecuteTemplate(c.Writer, "get_data.template", nil)
			if err != nil {
				err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
				if err != nil {
					panic(err.Error())
				}
				return
			}

			return
		}
		err := crud.GenerateData(getData)
		if err != nil {
			err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
			if err != nil {
				panic(err.Error())
			}
			return
		}

		err = models.T.ExecuteTemplate(c.Writer, "get_data.template", models.D)
		if err != nil {
			err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
			if err != nil {
				panic(err.Error())
			}
			return
		}

		return
	} else if getReport := c.Query("get_report"); getReport != "" { //get report button pressed
		if getReport == "Generate Report" {
			err := models.T.ExecuteTemplate(c.Writer, "get_data.template", nil)
			if err != nil {
				err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
				if err != nil {
					panic(err.Error())
				}
				return
			}
		}
		err := crud.GenerateReport(getReport)
		if err != nil {
			err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
			if err != nil {
				panic(err.Error())
			}
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/report")
	}

	err := models.T.ExecuteTemplate(c.Writer, "get_data.template", nil)
	if err != nil {
		err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
		if err != nil {
			panic(err.Error())
		}
		return
	}
	return
}
