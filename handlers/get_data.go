package handlers

import (
	"fmt"
	"net/http"

	"github.com/bejaneps/csv-webapp/crud"
	"github.com/gin-gonic/gin"
)

// GetDataHandler handles requests on /dashboard/get_data path
func GetDataHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	q := c.Query("get_data")
	if q == "" || q == "Get Data" {
		c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
		return
	}

	f, err := crud.GetData(q)
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
