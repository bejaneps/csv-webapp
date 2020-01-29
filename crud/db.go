package crud

import (
	"context"
	"errors"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bejaneps/csv-webapp/auth"
	"github.com/bejaneps/csv-webapp/models"
	"go.mongodb.org/mongo-driver/bson"
)

// CheckLoginInfo returns a boolean, that shows if a user exists in database or no
func CheckLoginInfo(info models.LoginInfo) (bool, error) {
	client, err := auth.NewMongoClient()
	if err != nil {
		return false, errors.New("CheckLoginInfo(): " + err.Error())
	}
	defer auth.CloseMongoClient()

	collection := client.Database("cdr").Collection("users")

	temp := models.LoginInfo{}
	user := bson.M{"email": info.Email, "password": info.Password}
	err = collection.FindOne(context.TODO(), user).Decode(&temp)
	if err != nil {
		return false, errors.New("CheckLoginInfo(): " + err.Error())
	}

	if temp.Email != "" && temp.Password != "" {
		return true, nil
	}

	return false, errors.New("CheckLoginInfo(): " + err.Error())
}

// getMongoCollections returns names of collections in MongoDB
func getMongoCollections(client *mongo.Client) ([]string, error) {
	names, err := client.Database("cdr").ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return nil, errors.New("getMongoCollections(): " + err.Error())
	}

	return names, nil
}

func hasEntry(entry string, entries []string) bool {
	for _, v := range entries {
		if v == entry {
			return true
		}
	}

	return false
}

// createMongoCollection creates a collection in a Mongo DB
func createMongoCollection(mongoColl <-chan string, mgoClient *mongo.Client, w *sync.WaitGroup, errChan chan<- error) {
	for m := range mongoColl {
		if len(models.D.Datum) == 0 {
			return
		}
		collection := mgoClient.Database("cdr").Collection(m)

		//can't use []Datum as type []interface{}
		temp := make([]interface{}, len(models.D.Datum))
		for i, v := range models.D.Datum {
			temp[i] = v
		}

		_, err := collection.InsertMany(context.TODO(), temp)
		if err != nil {
			errChan <- errors.New("createMongoCollection(): " + err.Error())
			return
		}

		log.Printf("[INFO]: created %s mongo collection\n", collection.Name())
	}
	w.Done()
}
