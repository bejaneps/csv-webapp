package crud

import (
	"encoding/csv"
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

// CSVToXLSX generates a csv file, then converts it to xlsx.
// Returns a full path to a file
func CSVToXLSX() (string, error) {
	content, err := csvutil.Marshal(&models.D.Datum)
	if err != nil {
		return "", err
	}

	err = os.Chdir("/tmp")
	if err != nil {
		return "", err
	}

	csvFile, err := os.Create("report.csv")
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	xlsxTemp, err := os.Create("report.xlsx")
	if err != nil {
		return "", err
	}
	defer xlsxTemp.Close()

	b, err := csvFile.Write(content)
	if err != nil {
		return "", err
	}
	if b == 0 {
		return "", errors.New("write: no bytes written")
	}

	reader := csv.NewReader(csvFile)
	reader.Comma = rune(',')

	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet(csvFile.Name())
	if err != nil {
		return "", err
	}

	fields, err := reader.Read()
	for err == nil {
		row := sheet.AddRow()
		for _, field := range fields {
			cell := row.AddCell()
			cell.Value = field
		}
		fields, err = reader.Read()
	}
	if err != nil {
		return "", err
	}

	err = xlsxFile.Save(xlsxTemp.Name())
	if err != nil {
		return "", nil
	}

	return "/tmp/" + xlsxTemp.Name(), nil
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
		return time.Time{}, time.Time{}, err
	}

	pEndRange, err := time.Parse(time.RFC3339, endRange)
	if err != nil {
		return time.Time{}, time.Time{}, err
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
		return err
	}

	if err = csvutil.Unmarshal(content, &models.D.Datum); err != nil {
		return err
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
