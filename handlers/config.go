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

	if val := c.Query("reset"); val != "" || !models.D.C.Initialized {
		crud.InitConfig()
		c.Redirect(http.StatusTemporaryRedirect, "/config")
	}

	var err error
	//Fixed to Mobile
	models.D.C.CostSecond["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_cost_second"), 64)
	if err != nil {
		if _, ok := models.D.C.CostSecond["Fixed to Mobile"]; !ok {
			models.D.C.CostSecond["Fixed to Mobile"] = 0
		}
	}
	models.D.C.MinSecond["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_min_second"), 64)
	if err != nil {
		if _, ok := models.D.C.MinSecond["Fixed to Mobile"]; !ok {
			models.D.C.MinSecond["Fixed to Mobile"] = 0
		}
	}
	models.D.C.Min["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_min"), 64)
	if err != nil {
		if _, ok := models.D.C.Min["Fixed to Mobile"]; !ok {
			models.D.C.Min["Fixed to Mobile"] = 0
		}
	}
	models.D.C.Fixed["Fixed to Mobile"], err = strconv.ParseFloat(c.Query("fixed_fixed"), 64)
	if err != nil {
		if _, ok := models.D.C.Fixed["Fixed to Mobile"]; !ok {
			models.D.C.Fixed["Fixed to Mobile"] = 0
		}
	}
	if val := c.Query("fixed_charge"); val != "" {
		models.D.C.Charge["Fixed to Mobile"] = val
	}

	//National
	models.D.C.CostSecond["National"], err = strconv.ParseFloat(c.Query("national_cost_second"), 64)
	if err != nil {
		if _, ok := models.D.C.CostSecond["National"]; !ok {
			models.D.C.CostSecond["National"] = 0
		}
	}
	models.D.C.MinSecond["National"], err = strconv.ParseFloat(c.Query("national_min_second"), 64)
	if err != nil {
		if _, ok := models.D.C.MinSecond["National"]; !ok {
			models.D.C.MinSecond["National"] = 0
		}
	}
	models.D.C.Min["National"], err = strconv.ParseFloat(c.Query("national_min"), 64)
	if err != nil {
		if _, ok := models.D.C.Min["National"]; !ok {
			models.D.C.Min["National"] = 0
		}
	}
	models.D.C.Fixed["National"], err = strconv.ParseFloat(c.Query("national_fixed"), 64)
	if err != nil {
		if _, ok := models.D.C.Fixed["National"]; !ok {
			models.D.C.Fixed["National"] = 0
		}
	}
	if val := c.Query("national_charge"); val != "" {
		models.D.C.Charge["National"] = val
	}

	//International
	models.D.C.CostSecond["International"], err = strconv.ParseFloat(c.Query("international_cost_second"), 64)
	if err != nil {
		if _, ok := models.D.C.CostSecond["International"]; !ok {
			models.D.C.CostSecond["International"] = 0
		}
	}
	models.D.C.MinSecond["International"], err = strconv.ParseFloat(c.Query("international_min_second"), 64)
	if err != nil {
		if _, ok := models.D.C.MinSecond["International"]; !ok {
			models.D.C.MinSecond["International"] = 0
		}
	}
	models.D.C.Min["International"], err = strconv.ParseFloat(c.Query("international_min"), 64)
	if err != nil {
		if _, ok := models.D.C.Min["International"]; !ok {
			models.D.C.Min["International"] = 0
		}
	}
	models.D.C.Fixed["International"], err = strconv.ParseFloat(c.Query("international_fixed"), 64)
	if err != nil {
		if _, ok := models.D.C.Fixed["International"]; !ok {
			models.D.C.Fixed["International"] = 0
		}
	}
	if val := c.Query("international_charge"); val != "" {
		models.D.C.Charge["International"] = val
	}

	//Intercapital City
	models.D.C.CostSecond["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_cost_second"), 64)
	if err != nil {
		if _, ok := models.D.C.CostSecond["Intercapital City"]; !ok {
			models.D.C.CostSecond["Intercapital City"] = 0
		}
	}
	models.D.C.MinSecond["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_min_second"), 64)
	if err != nil {
		if _, ok := models.D.C.MinSecond["Intercapital City"]; !ok {
			models.D.C.MinSecond["Intercapital City"] = 0
		}
	}
	models.D.C.Min["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_min"), 64)
	if err != nil {
		if _, ok := models.D.C.Min["Intercapital City"]; !ok {
			models.D.C.Min["Intercapital City"] = 0
		}
	}
	models.D.C.Fixed["Intercapital City"], err = strconv.ParseFloat(c.Query("intercapital_fixed"), 64)
	if err != nil {
		if _, ok := models.D.C.Fixed["Intercapital City"]; !ok {
			models.D.C.Fixed["Intercapital City"] = 0
		}
	}
	if val := c.Query("intercapital_charge"); val != "" {
		models.D.C.Charge["Intercapital City"] = val
	}

	//Special
	models.D.C.CostSecond["Special"], err = strconv.ParseFloat(c.Query("special_cost_second"), 64)
	if err != nil {
		if _, ok := models.D.C.CostSecond["Special"]; !ok {
			models.D.C.CostSecond["Special"] = 0
		}
	}
	models.D.C.MinSecond["Special"], err = strconv.ParseFloat(c.Query("special_min_second"), 64)
	if err != nil {
		if _, ok := models.D.C.MinSecond["Special"]; !ok {
			models.D.C.MinSecond["Special"] = 0
		}
	}
	models.D.C.Min["Special"], err = strconv.ParseFloat(c.Query("special_min"), 64)
	if err != nil {
		if _, ok := models.D.C.Min["Special"]; !ok {
			models.D.C.Min["Special"] = 0
		}
	}
	models.D.C.Fixed["Special"], err = strconv.ParseFloat(c.Query("special_fixed"), 64)
	if err != nil {
		if _, ok := models.D.C.Fixed["Special"]; !ok {
			models.D.C.Fixed["Special"] = 0
		}
	}
	if val := c.Query("special_charge"); val != "" {
		models.D.C.Charge["Special"] = val
	}

	c.Redirect(http.StatusTemporaryRedirect, "/dashboard")
}
