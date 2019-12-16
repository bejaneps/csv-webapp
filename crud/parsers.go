package crud

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/ini.v1"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/jszwec/csvutil"
)

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

	if err = csvutil.Unmarshal(content, &models.Datum); err != nil {
		return err
	}

	log.Printf("[INFO]: parsed %s file\n", f.Name())

	//Datum[0].FileName = file

	return nil
}

// ParseTemplates parses all templates
func ParseTemplates() error {
	var err error

	models.T, err = template.ParseGlob("templates/*")
	if err != nil {
		return err
	}

	return nil
}

// ParseINI parses ini file and umarshalls all data to global variables
func ParseINI(file string) error {
	cfg, err := ini.Load(file)
	if err != nil {
		return err
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
