package crud

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/tealeg/xlsx"

	"gopkg.in/ini.v1"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/bejaneps/csvutil"
)

// CSVToXLSX generates a csv file, then converts it to xlsx. Takes name of a file as a parameter.
// Returns a file itself
func CSVToXLSX(name string) (*os.File, error) {
	xlsxFile := xlsx.NewFile()

	sheet, err := xlsxFile.AddSheet("Sheet1")
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	//writing headers
	headers := struct {
		Five      string `csv:"0"`
		Nineteen  string `csv:"1"`
		TwentyOne string `csv:"2"`
	}{
		"Connect Datetime",
		"Called Number",
		"Location Pair Category",
	}

	row := sheet.AddRow()
	row.WriteStruct(&headers, -1)

	//writing data
	for _, val := range models.D.Datum {
		//writing report
		report := struct {
			Five      string `csv:"0"`
			Nineteen  int    `csv:"1"`
			TwentyOne string `csv:"2"`
		}{
			val.Five,
			val.Nineteen,
			val.TwentyOne,
		}
		row = sheet.AddRow()
		row.WriteStruct(&report, -1)
	}

	//writing headers of last row of report
	lHeaders := struct {
		FixedToMobile    string `csv:"0"`
		National         string `csv:"1"`
		International    string `csv:"2"`
		IntercapitalCity string `csv:"3"`
	}{
		"Fixed to Mobile",
		"National",
		"International",
		"Intercapital City",
	}
	row = sheet.AddRow()
	row.WriteStruct(&lHeaders, -1)

	//writing last row of report
	row = sheet.AddRow()
	row.WriteStruct(&models.D.TC, -1)

	err = xlsxFile.Save("/tmp/" + name)
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	f, err := os.Open("/tmp/" + name)
	if err != nil {
		return nil, errors.New("CSVToXLSX(): " + err.Error())
	}

	return f, nil
}

// parseHTMLTime parses time that is get from a server
func parseHTMLTime(t string) (start, end time.Time, err error) {
	// initialize variables
	startRange := t[:strings.Index(t, "-")-1]
	endRange := t[strings.Index(t, "-")+2:]

	//make it rfc3339
	startRange = startRange[6:] + "-" + startRange[:2] + "-" + startRange[3:5] + "T00:00:00Z"
	endRange = endRange[6:] + "-" + endRange[:2] + "-" + endRange[3:5] + "T00:00:00Z"

	//convert it to time type
	pStartRange, err := time.Parse(time.RFC3339, startRange)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("CSVToXLSX(): " + err.Error())
	}

	pEndRange, err := time.Parse(time.RFC3339, endRange)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("CSVToXLSX(): " + err.Error())
	}

	return pStartRange, pEndRange, nil
}

// parseCSV parses a csv file and unmarshals all data in slice struct
func parseCSV(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New("CSVToXLSX(): " + err.Error())
	}

	if err = csvutil.Unmarshal(content, &models.D.Datum); err != nil {
		return errors.New("CSVToXLSX(): " + err.Error())
	}

	log.Printf("[INFO]: parsed %s file\n", f.Name())

	return nil
}

// parseCSV parses a csv file and unmarshals all data in slice struct
func parseCSVRange(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.New("parseCSVRange(): " + err.Error())
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New("parseCSVRange(): " + err.Error())
	}

	temp := []models.CDRModified{}
	if err = csvutil.Unmarshal(content, &temp); err != nil {
		return errors.New("parseCSVRange(): " + err.Error())
	}

	for _, v := range temp {
		if strings.Contains(v.TwentyOne, "Fixed") {
			models.D.TC.FixedToMobile += v.Eleven
		} else if strings.Contains(v.TwentyOne, "International") {
			models.D.TC.International += v.Eleven
		} else if strings.Contains(v.TwentyOne, "National") {
			models.D.TC.National += v.Eleven
		} else {
			models.D.TC.IntercapitalCity += v.Eleven
		}

		models.D.Datum = append(models.D.Datum, v)
	}

	log.Printf("[INFO]: parsed %s file\n", f.Name())

	return nil
}

// ParseTemplates parses all templates
func ParseTemplates() error {
	var err error

	models.T, err = template.ParseGlob("templates/*")
	if err != nil {
		return errors.New("ParseTemplates(): " + err.Error())
	}

	return nil
}

// ParseINI parses ini file and umarshalls all data to global variables
func ParseINI(file string) error {
	cfg, err := ini.Load(file)
	if err != nil {
		return errors.New("ParseINI(): " + err.Error())
	}

	models.FTPURI = cfg.Section("common").Key("ftp_uri").String()
	if models.FTPURI == "" {
		return errors.New("empty ftp_uri")
	}

	models.FTPLogin = cfg.Section("common").Key("ftp_login").String()
	if models.FTPLogin == "" {
		return errors.New("empty ftp_login")
	}

	models.FTPPassword = cfg.Section("common").Key("ftp_password").String()
	if models.FTPPassword == "" {
		return errors.New("empty ftp_password")
	}

	models.Port = ":" + cfg.Section("common").Key("port").String()
	if models.Port == "" {
		return errors.New("empty port")
	}

	return nil
}
