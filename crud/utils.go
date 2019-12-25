package crud

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/tealeg/xlsx"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// randomStringWithCharset returns random string from names of a files
func randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// randomString ...
func randomString(length int) string {
	return randomStringWithCharset(length, charset)
}

// cleanTmp cleans files created in 'tmp' directory
func cleanTmp(fileName string) {
	_ = os.Remove(fileName)
}

// generateXLSX generates an xlsx file for report button
func generateXLSX(button string) (*os.File, error) {
	xlsxFile := xlsx.NewFile()

	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	if button == "get_data" {
		headers := []string{"Connect Datetime", "Disconnect Datetime", "Charged Duration (Seconds)", "Charged Duration (Minutes)", "Calling Number", "Called Number", "Location Pair Category", "Charged Amount", "Currency Code", "Completion Code ID", "Completion Code Name", "Sell"}

		row := sheet.AddRow()
		row.WriteSlice(&headers, -1)

		for _, val := range models.D.Datum {
			if val.Five == "" {
				continue
			}

			row = sheet.AddRow()
			row.WriteStruct(&val, -1)
		}
	} else {
		headers := []string{"Connect Datetime", "Called Number", "Location Pair Category"}

		row := sheet.AddRow()
		row.WriteSlice(&headers, -1)

		//writing data
		for _, val := range models.D.Datum {
			//writing report
			report := struct {
				Five      string  `csv:"0"`
				Nineteen  int     `csv:"1"`
				TwentyOne string  `csv:"2"`
				Sell      float64 `csv:"3"`
			}{
				val.Five,
				val.Nineteen,
				val.TwentyOne,
				val.Sell,
			}

			if val.Five == "" {
				continue
			}
			row = sheet.AddRow()
			row.WriteStruct(&report, -1)
		}

		//blank rows
		sheet.AddRow()
		sheet.AddRow()

		//writing headers of last row of report
		lHeaders := struct {
			FixedToMobile    string `csv:"0"`
			National         string `csv:"1"`
			International    string `csv:"2"`
			IntercapitalCity string `csv:"3"`
			Special          string `csv:"4"`
		}{
			"Fixed to Mobile",
			"National",
			"International",
			"Intercapital City",
			"Special",
		}
		row = sheet.AddRow()
		row.WriteStruct(&lHeaders, -1)

		//writing last row of report
		row = sheet.AddRow()
		row.WriteStruct(&models.D.TC, -1)
	}

	fileName := randomString(10) + ".xlsx"
	err = xlsxFile.Save("/tmp/" + fileName)
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	f, err := os.Open("/tmp/" + fileName)
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	log.Printf("[INFO]: xlsx file created.\n")

	return f, nil
}

// InitConfig initialies struct of config
func InitConfig() {
	models.D.C.CostSecond = make(map[string]float64)
	models.D.C.MinSecond = make(map[string]float64)
	models.D.C.Min = make(map[string]float64)
	models.D.C.Fixed = make(map[string]float64)
	models.D.C.Charge = make(map[string]string)
}
