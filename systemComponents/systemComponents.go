package systemComponents

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
)

func GetIDFromDB(collection *mongo.Collection) (string, error) {
	if numID, err := collection.CountDocuments(context.TODO(), bson.M{}); err != nil {
		log.Println("Can't get new ID: ", err)
		return "", err
	} else {
		ID := strconv.Itoa(int(numID))
		return ID, nil
	}
}
