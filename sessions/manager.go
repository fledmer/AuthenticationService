package sessions

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
)

type cookieSession struct {
	Token string `bson:"token"`
	Login string `bson:"login"`
}

var dbClient *mongo.Client
var sessionCollection *mongo.Collection

func Init(client *mongo.Client) error {
	if client == nil {
		return errors.New("nil DbClient")
	}
	dbClient = client
	sessionCollection = dbClient.Database("test").Collection("sessions")
	return nil
}

func DeleteSession(token string) (deletedCount int64, err error) {
	if deleted, err := sessionCollection.DeleteOne(context.TODO(), bson.D{{"token", token}}); err == nil {
		log.Println("Success to delete session, token: ", token, " count: ", deleted.DeletedCount)
		return deleted.DeletedCount, err
	} else {
		log.Println("Error to delete session, toke: ", token, " error: ", err)
		return 0, err
	}
}

func CreateSessions(login string) (string, error) {
	token := getNewSessionToken()
	_, err := sessionCollection.InsertOne(context.TODO(), cookieSession{token, login})
	if err == nil {
		log.Println("Create new session, login: ", login, " token: ", token)
	} else {
		log.Println("Error to create new session, login: ", login, " error: ", err)
	}
	return token, err
}

func GetSession(token string) (login string, err error) {
	var obj cookieSession
	err = sessionCollection.FindOne(context.TODO(), bson.D{{"token", token}}).Decode(&obj)
	if err == nil {
		log.Println("Get session, token: ", token, " login: ", obj.Login)
	} else {
		log.Println("Error to get session login: ", login, " error: ", err)
	}
	return obj.Login, err
}

func getNewSessionToken() string {
	id := ""
	s := rand.Int31()%15 + 15
	for x := int32(0); x < s; x++ {
		id += string(rand.Int31()%25 + 65)
	}
	return id
}
