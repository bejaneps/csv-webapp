package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/bejaneps/csv-webapp/crud"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset returns random string from names of a files
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String -
func String(length int) string {
	return StringWithCharset(length, charset)
}

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

// ReportDownloadHandler handles downloads on parsed file
func ReportDownloadHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	fileName := String(7)
	file, err := crud.CSVToXLSX(fileName + ".xlsx")
	if err != nil {
		err = models.T.ExecuteTemplate(c.Writer, "error.template", err.Error())
		if err != nil {
			panic(err.Error())
		}
		return
	}
	fileStats, _ := file.Stat()

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, fileStats.Name()),
	}

	c.DataFromReader(http.StatusOK, fileStats.Size(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file, extraHeaders)
}
