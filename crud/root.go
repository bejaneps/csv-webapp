package crud

import (
	"errors"
	"os"
	"strings"

	"github.com/bejaneps/csv-webapp/auth"
	"github.com/bejaneps/csv-webapp/models"
)

// GetData is a function for get data button
func GetData(timeRange string) (*os.File, error) {
	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}
	defer auth.CloseFTPConnection()

	mgoClient, err := auth.NewMongoClient()
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}
	defer auth.CloseMongoClient()

	mgoEntries, err := getMongoCollections(mgoClient)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	ftpEntries, err := getFTPEntries(ftpConn)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	for _, v := range ftpEntries {
		//empty file
		if v.Size == 297 {
			continue
		} //montly files
		if len(v.Name) > 38 {
			continue
		}

		noGZName := strings.TrimSuffix(v.Name, ".gz")

		currDir, _ := os.Getwd()
		if ok := hasEntry(noGZName, mgoEntries); !ok {
			fileName, err := createFTPFile(v.Name, currDir+"/"+"files", ftpConn)
			if err != nil {
				return nil, errors.New("GetData(): " + err.Error())
			}

			err = parseCSV(fileName)
			if err != nil {
				return nil, errors.New("GetData(): " + err.Error())
			}

			err = createMongoCollection(noGZName, mgoClient)
			if err != nil {
				return nil, errors.New("GetData(): " + err.Error())
			}
		}
	}

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	models.D.Datum = []models.CDRModified{}
	for _, v := range rangeEntries {
		//empty file
		if v.Size == 297 {
			continue
		}
		//monthly files
		if len(v.Name) > 38 {
			continue
		}
		fileName := strings.TrimSuffix("files/"+v.Name, ".gz")
		err = parseCSV(fileName)
		if err != nil {
			return nil, errors.New("GetData(): " + err.Error())
		}
	}

	f, err := generateXLSX("get_data")
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	return f, nil
}

// GenerateReport is function for get report button
func GenerateReport(timeRange string) (*os.File, error) {
	if timeRange == "Generate Report" {
		f, err := generateXLSX("generate_report")
		if err != nil {
			return nil, errors.New("GenerateReport(): " + err.Error())
		}

		return f, nil
	}

	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}
	defer auth.CloseFTPConnection()

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	models.D.Datum = []models.CDRModified{}
	for _, v := range rangeEntries {
		//empty file
		if v.Size == 297 {
			continue
		}
		//monthly files
		if len(v.Name) > 38 {
			continue
		}

		fileName := strings.TrimSuffix("files/"+v.Name, ".gz")
		err = parseCSV(fileName)
		if err != nil {
			return nil, errors.New("GenerateReport(): " + err.Error())
		}
	}

	f, err := generateXLSX("generate_report")
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	return f, nil
}
