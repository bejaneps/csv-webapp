package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bejaneps/csv-webapp/models"
	"github.com/jlaffaye/ftp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//connString = "mongodb+srv://bejanhtc:246150Sab!@cluster0-6tflo.mongodb.net/cdr"
	//connString = "mongodb://mongodb:27017"
	connString = "mongodb://127.0.0.1:27017"
)

var (
	err error
)

// NewMongoClient returns new MongoDB client instance
func NewMongoClient() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	mgoClient, err := mongo.Connect(ctx, options.Client().SetSocketTimeout(5*time.Hour).ApplyURI(connString).SetMaxPoolSize(10).SetMaxConnIdleTime(2*time.Minute).SetSocketTimeout(2*time.Minute))
	if err != nil {
		return nil, errors.New("NewMongoClient(): " + err.Error())
	}

	// Check the connection
	err = mgoClient.Ping(context.TODO(), nil)
	if err != nil {
		return nil, errors.New("NewMongoClient(): " + err.Error())
	}

	return mgoClient, nil
}

// NewFTPConnection returns a connection to ftp server
func NewFTPConnection() (*ftp.ServerConn, error) {
	ftpConn, err := ftp.Dial(models.FTPURI, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, errors.New("NewFTPConnection(): " + err.Error())
	}

	err = ftpConn.Login(models.FTPLogin, models.FTPPassword)
	if err != nil {
		return nil, errors.New("NewFTPConnection(): " + err.Error())
	}

	return ftpConn, nil
}

// CloseMongoClient closes a connection to MongoDB client
func CloseMongoClient(mgoClient *mongo.Client) error {
	if err := mgoClient.Disconnect(context.TODO()); err != nil {
		return errors.New("CloseMongoClient(): " + err.Error())
	}

	return nil
}

// CloseFTPConnection closes a connection to ftp server
func CloseFTPConnection(ftpConn *ftp.ServerConn) error {
	if err := ftpConn.Quit(); err != nil {
		return errors.New("CloseFTPConnection(): " + err.Error())
	}

	return nil
}
