package crud

import (
	"os"
	"strings"

	"github.com/bejaneps/csv-webapp/auth"
)

// cleanTmp cleans files created in 'tmp' directory
func cleanTmp(fileName string) {
	_ = os.Remove("/tmp" + fileName)
	_ = os.Remove("/tmp" + fileName + ".gz")
}

// GenerateData is a function for get data button
func GenerateData() error {
	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return err
	}
	defer auth.CloseFTPConnection()

	mgoClient, err := auth.NewMongoClient()
	if err != nil {
		return err
	}

	ftpEntries, err := getFTPEntries(ftpConn)
	if err != nil {
		return err
	}

	mgoEntries, err := getMongoCollections(mgoClient)
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
			fileName, err := createFTPFile(v.Name, "/tmp", ftpConn)
			if err != nil {
				return err
			}

			err = parseCSV(fileName)
			cleanTmp(fileName)
			if err != nil {
				return err
			}

			err = createMongoCollection(noGZName, mgoClient)
			if err != nil {
				return err
			}
		}
	}

	latest, err := getLatestFTPFile(ftpConn)
	if err != nil {
		return err
	}

	err = parseCSV("files/" + strings.TrimSuffix(latest, ".gz"))
	if err != nil {
		return err
	}

	return nil
}
