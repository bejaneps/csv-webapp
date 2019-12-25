package handlers

import (
	"net/http"
	"strconv"

	"github.com/bejaneps/csv-webapp/crud"

	"github.com/bejaneps/csv-webapp/models"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles upcoming requests on /config path
func ConfigHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.HTML(http.StatusOK, "config.html", nil)
}

// ConfigSubmitHandler handles upcoming requests on /config/submit path
func ConfigSubmitHandler(c *gin.Context) {
	val, _ := c.Cookie("auth")
	if val == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	var err error

	crud.InitConfig()

	//Fixed to Mobile
	models.D.C.CostSecond["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_cost_second"), 64)
	if err != nil {
		models.D.C.CostSecond["Fixed to Mobile"] = 0
	}
	models.D.C.MinSecond["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_min_second"), 64)
	if err != nil {
		models.D.C.MinSecond["Fixed to Mobile"] = 0
	}
	models.D.C.Min["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_min"), 64)
	if err != nil {
		models.D.C.Min["Fixed to Mobile"] = 0
	}
	models.D.C.Fixed["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_fixed"), 64)
	if err != nil {
		models.D.C.Fixed["Fixed to Mobile"] = 0
	}
	models.D.C.Charge["Fixed to Mobile"] = c.Query("fixed_charge")

	//National
	models.D.C.CostSecond["National"], err = strconv.ParseFloat(c.Query("national_cost_second"), 64)
	if err != nil {
		models.D.C.CostSecond["National"] = 0
	}
	models.D.C.MinSecond["National"], err = strconv.ParseFloat(c.Query("national_min_second"), 64)
	if err != nil {
		models.D.C.MinSecond["National"] = 0
	}
	models.D.C.Min["National"], err = strconv.ParseFloat(c.Query("national_min"), 64)
	if err != nil {
		models.D.C.Min["National"] = 0
	}
	models.D.C.Fixed["National"], err = strconv.ParseFloat(c.Query("national_fixed"), 64)
	if err != nil {
		models.D.C.Fixed["National"] = 0
	}
	models.D.C.Charge["National"] = c.Query("national_charge")

	//International
	models.D.C.CostSecond["International"], err = strconv.ParseFloat(c.Query("international_cost_second"), 64)
	if err != nil {
		models.D.C.CostSecond["International"] = 0
	}
	models.D.C.MinSecond["International"], err = strconv.ParseFloat(c.Query("international_min_second"), 64)
	if err != nil {
		models.D.C.MinSecond["International"] = 0
	}
	models.D.C.Min["International"], err = strconv.ParseFloat(c.Query("international_min"), 64)
	if err != nil {
		models.D.C.Min["International"] = 0
	}
	models.D.C.Fixed["International"], err = strconv.ParseFloat(c.Query("international_fixed"), 64)
	if err != nil {
		models.D.C.Fixed["International"] = 0
	}
	models.D.C.Charge["International"] = c.Query("international_charge")

	//Intercapital City
	models.D.C.CostSecond["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_cost_second"), 64)
	if err != nil {
		models.D.C.CostSecond["Intercapital City"] = 0
	}
	models.D.C.MinSecond["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_min_second"), 64)
	if err != nil {
		models.D.C.MinSecond["Intercapital City"] = 0
	}
	models.D.C.Min["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_min"), 64)
	if err != nil {
		models.D.C.Min["Intercapital City"] = 0
	}
	models.D.C.Fixed["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_fixed"), 64)
	if err != nil {
		models.D.C.Fixed["Intercapital City"] = 0
	}
	models.D.C.Charge["Intercapital City"] = c.Query("intercapital_charge")

	//Special
	models.D.C.CostSecond["Special"], err = strconv.ParseFloat(c.Query("special_cost_second"), 64)
	if err != nil {
		models.D.C.CostSecond["Special"] = 0
	}
	models.D.C.MinSecond["Special"], err = strconv.ParseFloat(c.Query("special_min_second"), 64)
	if err != nil {
		models.D.C.MinSecond["Special"] = 0
	}
	models.D.C.Min["Special"], err = strconv.ParseFloat(c.Query("special_min"), 64)
	if err != nil {
		models.D.C.Min["Special"] = 0
	}
	models.D.C.Fixed["Special"], err = strconv.ParseFloat(c.Query("special_fixed"), 64)
	if err != nil {
		models.D.C.Fixed["Special"] = 0
	}
	models.D.C.Charge["Special"] = c.Query("special_charge")

	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}
