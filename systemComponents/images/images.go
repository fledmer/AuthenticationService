package images

import (
	"context"
	"errors"
	"log"
	"main/systemComponents"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Images systemComponents.Images

var (
	dbClient         *mongo.Client
	imagesCollection *mongo.Collection
)

func Init(client *mongo.Client) error {
	if client == nil {
		return errors.New("nil DbClient")
	}
	dbClient = client
	imagesCollection = dbClient.Database("test").Collection("images")
	return nil
}

func (image *Images) AddNew() error {
	if id, err := systemComponents.GetIDFromDB(imagesCollection); err != nil {
		return err
	} else {
		image.ID = id
	}
	if _, err := imagesCollection.InsertOne(context.TODO(), image); err != nil {
		log.Println("Can't register image, err:", err)
		return err
	} else {
		return nil
	}
}

func GetByID(id string) (image *Images, err error) {
	image = new(Images)
	if err = imagesCollection.FindOne(context.TODO(), bson.M{"ID": id}).Decode(image); err != nil {
		return nil, err
	} else {
		return image, nil
	}
}

func DeleteByID(id string) error {
	if _, err := imagesCollection.DeleteOne(context.TODO(), bson.M{"ID": id}); err != nil {
		return err
	} else {
		return nil
	}
}
