package crud

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/bejaneps/csv-webapp/auth"

	"github.com/jlaffaye/ftp"
)

// getFTPEntries returns list of entries in a server.
func getFTPEntries(conn *ftp.ServerConn) ([]*ftp.Entry, error) {
	entries, err := conn.List("/")
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func getLatestFTPFile(conn *ftp.ServerConn) (string, error) {
	entries, err := conn.List("/")
	if err != nil {
		return "", err
	}

	latestFile := entries[0]
	for _, v := range entries {
		if v.Time.UnixNano() > latestFile.Time.UnixNano() {
			latestFile = v
		}
	}

	return latestFile.Name, nil
}

// getRangeEntries returns list of files between start and end dates from ftp
func getRangeEntries(start, end time.Time, conn *ftp.ServerConn) ([]*ftp.Entry, error) {
	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		return nil, err
	}
	defer auth.CloseFTPConnection()

	entries, err := getFTPEntries(ftpConn)
	if err != nil {
		return nil, err
	}

	var rangeFiles []*ftp.Entry
	for _, v := range entries {
		if ok := v.Time.After(start); ok {
			if ok := v.Time.Before(end); ok {
				rangeFiles = append(rangeFiles, v)
			}
		}
	}

	return rangeFiles, nil
}

// createFTPFile gets file from a server & unzips it in a specified folder. Returns the name of created file
func createFTPFile(ftpFileName <-chan string, conn *ftp.ServerConn, csvFileName chan<- string, w *sync.WaitGroup, errChan chan<- error) {
	for m := range ftpFileName {
		resp, err := conn.Retr(m)
		if err != nil {
			errChan <- errors.New("createFTPFile(): " + err.Error())
			return
		}

		f, err := os.Create(m)
		if err != nil {
			errChan <- errors.New("createFTPFile(): " + err.Error())
			return
		}
		defer f.Close()

		log.Printf("[INFO]: downloaded %s file\n", f.Name())

		if _, err := io.Copy(f, resp); err != nil {
			errChan <- errors.New("createFTPFile(): " + err.Error())
			return
		}
		resp.Close()

		cmd := exec.Command("gunzip", f.Name())
		if err = cmd.Run(); err != nil {
			os.Remove(f.Name())
			errChan <- errors.New("createFTPFile(): " + err.Error())
			return
		}

		csvFileName <- strings.TrimSuffix(f.Name(), ".gz")
	}
	w.Done()
	close(csvFileName)
}
