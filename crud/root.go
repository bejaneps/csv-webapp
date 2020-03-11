package crud

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bejaneps/csv-webapp/auth"
	"github.com/bejaneps/csv-webapp/models"
)

// GetData is a function for get data button
func GetData(timeRange string) (*os.File, error) {
	now := time.Now()

	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}
	defer auth.CloseFTPConnection(ftpConn)

	mgoEntries, err := getMongoCollections()
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	ftpEntries, err := getFTPEntries(ftpConn)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	w := &sync.WaitGroup{}

	ftpFileName := make(chan string)
	csvFileName := make(chan string)
	mongoCollName := make(chan string)
	errChan := make(chan error)

	//start main goroutines
	go createFTPFile(ftpFileName, ftpConn, csvFileName, w, errChan)
	go parseCSV(csvFileName, mongoCollName, w, errChan)
	go createMongoCollection(mongoCollName, w, errChan)
	w.Add(3)

	//debug errors on separate goroutine
	go func() {
		for err := range errChan {
			fmt.Printf("[ERROR]: %s\n", err)
		}
	}()

	//change working folder to files, for creating them there
	currDir, _ := os.Getwd()
	err = os.Chdir(currDir + "/files")
	if err != nil {
		log.Fatal(err)
	}

	//check if collection already exists in mongodb
	//if no, then download, and create it
	for _, v := range ftpEntries {
		//empty file
		if v.Size == 297 {
			continue
		} //montly files
		if len(v.Name) > 38 {
			continue
		}

		noGZName := strings.TrimSuffix(v.Name, ".gz")

		if ok := hasEntry(noGZName, mgoEntries); !ok {
			ftpFileName <- v.Name
		}
	}
	close(ftpFileName)

	w.Wait()

	elapsed := time.Now().Sub(now)
	log.Printf("elapsed time: %vs", elapsed.Seconds())

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	csvFileName = make(chan string)

	go parseCSV(csvFileName, nil, w, errChan)
	w.Add(1)

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

		csvFileName <- strings.TrimSuffix(v.Name, ".gz")
	}
	close(csvFileName)

	w.Wait()
	close(errChan)

	f, err := generateXLSX("get_data")
	if err != nil {
		return nil, errors.New("GetData(): " + err.Error())
	}

	err = os.Chdir(currDir)
	if err != nil {
		log.Fatal(err)
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
	defer auth.CloseFTPConnection(ftpConn)

	start, end, err := parseHTMLTime(timeRange)
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	rangeEntries, err := getRangeEntries(start, end, ftpConn)
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	//change working folder to files, for creating them there
	currDir, _ := os.Getwd()
	err = os.Chdir(currDir + "/files")
	if err != nil {
		log.Fatal(err)
	}

	w := &sync.WaitGroup{}

	csvFileName := make(chan string)
	errChan := make(chan error)

	go parseCSV(csvFileName, nil, w, errChan)
	w.Add(1)

	//debug errors on separate goroutine
	go func() {
		for err := range errChan {
			fmt.Printf("[ERROR]: %s\n", err)
		}
	}()

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

		csvFileName <- strings.TrimSuffix(v.Name, ".gz")
	}
	close(csvFileName)

	w.Wait()
	close(errChan)

	f, err := generateXLSX("generate_report")
	if err != nil {
		return nil, errors.New("GenerateReport(): " + err.Error())
	}

	err = os.Chdir(currDir)
	if err != nil {
		log.Fatal(err)
	}

	return f, nil
}
