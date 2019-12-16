package auth

import (
	"context"
	"time"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/jlaffaye/ftp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connString = "mongodb://127.0.0.1:27017"
)

var (
	ftpConn = &ftp.ServerConn{}

	mgoClient = &mongo.Client{}
	err       error
)

// NewMongoClient returns new MongoDB client instance
func NewMongoClient() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	mgoClient, err = mongo.Connect(ctx, options.Client().SetSocketTimeout(5*time.Hour).ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = mgoClient.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return mgoClient, nil
}

// NewFTPConnection returns a connection to ftp server
func NewFTPConnection() (*ftp.ServerConn, error) {
	ftpConn, err = ftp.Dial(models.FTPURI, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = ftpConn.Login(models.FTPLogin, models.FTPPassword)
	if err != nil {
		return nil, err
	}

	return ftpConn, nil
}

// CloseMongoClient closes a connectio to MongoDB cluster
func CloseMongoClient() error {
	if err := mgoClient.Disconnect(context.TODO()); err != nil {
		return err
	}

	return nil
}

// CloseFTPConnection closes a connection to ftp server
func CloseFTPConnection() error {
	if err := ftpConn.Quit(); err != nil {
		return err
	}

	return nil
}
