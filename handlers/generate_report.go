package handlers

import (
	"fmt"
	"net/http"

	"github.com/bejaneps/csv-webapp/crud"
	"github.com/gin-gonic/gin"
)

// GenerateReportHandler handles requests on /dashboard/generate_report path
func GenerateReportHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	q := c.Query("generate_report")
	if q == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
		return
	}

	f, err := crud.GenerateReport(q)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", err.Error())
	}
	defer f.Close()

	fStat, _ := f.Stat()

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fStat.Name()),
	}

	c.DataFromReader(http.StatusOK, fStat.Size(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f, extraHeaders)
}
