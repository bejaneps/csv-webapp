package crud

import (
	"os"
	"strings"

	"github.com/bejaneps/csv-webapp/models"

	"github.com/bejaneps/csv-webapp/auth"
)

// cleanTmp cleans files created in 'tmp' directory
func cleanTmp(fileName string) {
	_ = os.Remove(fileName)
}

// GenerateData is a function for get data button
func GenerateData(timeRange string) error {
	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return err
	}
	defer auth.CloseFTPConnection()

	mgoClient, err := auth.NewMongoClient()
	if err != nil {
		return err
	}

	mgoEntries, err := getMongoCollections(mgoClient)
	if err != nil {
		return err
	}

	ftpEntries, err := getFTPEntries(ftpConn)
	if err != nil {
		return err
	}

	for _, v := range ftpEntries {
		//empty file
		if v.Size == 297 {
			continue
		}

		noGZName := strings.TrimSuffix(v.Name, ".gz")

		if ok := hasEntry(noGZName, mgoEntries); !ok {
			fileName, err := createFTPFile(v.Name, "files", ftpConn)
			if err != nil {
				return err
			}

			err = parseCSV(fileName)
			if err != nil {
				return err
			}

			err = createMongoCollection(noGZName, mgoClient)
			if err != nil {
				return err
			}
		}
	}

	//for range files
	models.D = models.Data{}

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return err
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return err
	}

	for _, v := range rangeEntries {
		fileName := strings.TrimSuffix("files/"+v.Name, ".gz")
		err = parseCSVRange(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateReport is function for get report button
func GenerateReport(timeRange string) error {
	models.D = models.Data{}

	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return err
	}
	defer auth.CloseFTPConnection()

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return err
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return err
	}

	for _, v := range rangeEntries {
		fileName := strings.TrimSuffix("files/"+v.Name, ".gz")
		err = parseCSVRange(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}
