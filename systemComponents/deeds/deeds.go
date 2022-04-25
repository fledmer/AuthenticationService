package deeds

import (
	"context"
	"errors"
	"log"
	"main/systemComponents"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Deed systemComponents.Deed

var (
	dbClient        *mongo.Client
	deedsCollection *mongo.Collection
)

func Init(client *mongo.Client) error {
	if client == nil {
		return errors.New("nil DbClient")
	}
	dbClient = client
	deedsCollection = dbClient.Database("test").Collection("deeds")
	return nil
}

func (deed Deed) Registration() error {
	if id, err := systemComponents.GetIDFromDB(deedsCollection); err != nil {
		return err
	} else {
		deed.ID = id
	}
	if _, err := deedsCollection.InsertOne(context.TODO(), deed); err != nil {
		log.Println("Can't register deed, err:", err)
		return err
	} else {
		return nil
	}
}

func GetAllDeeds() (deeds []Deed, err error) {
	return getDeedsByFilter(bson.M{})
}

func GetDeedsByUserID(ID string) (deeds []Deed, err error) {
	return getDeedsByFilter(bson.M{"creatorID": ID})
}

func GetDeedsByID(ID string) (deed Deed, err error) {
	deeds, err := getDeedsByFilter(bson.M{"ID": ID})
	if err == nil && len(deeds) > 0 {
		return deeds[0], err
	}
	return deed, errors.New("Not found")
}

func getDeedsByFilter(filter interface{}) (deeds []Deed, err error) {
	if cursor, err := deedsCollection.Find(context.TODO(), filter); err != nil {
		log.Println("Can't take all deeds, err:", err)
		return nil, err
	} else {
		err = cursor.All(context.TODO(), &deeds)
		if err != nil {
			log.Println("Can't decode deeds, err: ", err)
		}
		return deeds, err
	}
}
