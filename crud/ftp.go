package crud

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// createFTPFile gets file from a server & unzips it in a specified folder. Returns the name of created file
func createFTPFile(name, dir string, conn *ftp.ServerConn) (string, error) {
	resp, err := conn.Retr(name)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	err = os.Chdir(dir)
	if err != nil {
		return "", err
	}

	f, err := os.Create(name)
	if err != nil {
		return "", err
	}

	log.Printf("[INFO]: created %s file\n", f.Name())

	if _, err := io.Copy(f, resp); err != nil {
		return "", err
	}

	if err := resp.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command("gunzip", f.Name())
	if err = cmd.Run(); err != nil {
		return "", err
	}

	err = os.Chdir(currDir)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(dir+"/"+f.Name(), ".gz"), nil
}

// DownloadFTPFiles downloads all ftp files, if they are not downloaded yet
func DownloadFTPFiles(e chan error) {
	ftpConn, err := auth.NewFTPConnection()
	if err != nil {
		e <- err
	}
	defer auth.CloseFTPConnection()

	currDir, err := os.Getwd()
	if err != nil {
		e <- err
	}

	var files []string
	err = filepath.Walk(currDir+"/files", func(path string, info os.FileInfo, err error) error {
		files = append(files, info.Name()+".gz")
		return err
	})
	if err != nil {
		e <- err
	}

	ftpFiles, err := ftpConn.NameList("/")
	if err != nil {
		e <- err
	}

	for _, v := range ftpFiles {
		v = strings.TrimPrefix(v, "/")
		if ok := hasEntry(v, files); !ok {
			name, err := createFTPFile(v, currDir+"/files", ftpConn)
			if err != nil {
				e <- err
			}
			log.Printf("[INFO]: file %s has been downloaded\n", name)
		}
	}

	e <- nil
}
