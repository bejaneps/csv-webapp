package crud

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/bejaneps/csvutil"
)

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
	//TODO: remove empty connect datetime columns(rows)
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
		//for TC
		if strings.Contains(v.TwentyOne, "Fixed") {
			models.D.TC.FixedToMobile += v.Eleven
		} else if strings.Contains(v.TwentyOne, "International") {
			models.D.TC.International += v.Eleven
		} else if strings.Contains(v.TwentyOne, "National") {
			models.D.TC.National += v.Eleven
		} else if strings.Contains(v.TwentyOne, "Intercapital") {
			models.D.TC.IntercapitalCity += v.Eleven
		} else {
			models.D.TC.Special += v.Eleven
		}

		//for Config
		done := false
		if v.TwentyTwo == 0 {
			if models.D.C.Charge[v.TwentyOne] == "N" || models.D.C.Charge[v.TwentyOne] == "n" {
				v.Sell = 0
				done = true
			} else if models.D.C.Charge[v.TwentyOne] == "Y" || models.D.C.Charge[v.TwentyOne] == "y" {
				if models.D.C.Fixed[v.TwentyOne] != 0 {
					v.Sell = models.D.C.Fixed[v.TwentyOne]
					done = true
				} else if models.D.C.MinSecond[v.TwentyOne] != 0 {
					if v.Ten < models.D.C.MinSecond[v.TwentyOne] {
						v.Ten = models.D.C.MinSecond[v.TwentyOne]
					}
				}
			}
		} else if models.D.C.Fixed[v.TwentyOne] != 0 {
			v.Sell = models.D.C.Fixed[v.TwentyOne]
			done = true
		}

		if !done {
			if models.D.C.MinSecond[v.TwentyOne] != 0 {
				if v.Ten < models.D.C.MinSecond[v.TwentyOne] {
					v.Ten = models.D.C.MinSecond[v.TwentyOne]
				}
			}

			amount := models.D.C.CostSecond[v.TwentyOne] * v.Ten
			if amount < models.D.C.Min[v.TwentyOne] {
				v.Sell = models.D.C.Min[v.TwentyOne]
			} else {
				v.Sell = amount
			}
		}

		models.D.Datum = append(models.D.Datum, v)
	}

	log.Printf("[INFO]: parsed %s file\n", f.Name())

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
