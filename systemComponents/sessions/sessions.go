package sessions

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type cookieSession struct {
	ID    string `bson:"ID"`
	Token string `bson:"token"`
}

var (
	dbClient          *mongo.Client
	sessionCollection *mongo.Collection
)

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

func CreateSessions(ID string) (string, error) {
	token := getNewSessionToken()
	_, err := sessionCollection.InsertOne(context.TODO(), cookieSession{token, ID})
	if err == nil {
		log.Println("Create new session, ID: ", ID, " token: ", token)
	} else {
		log.Println("Error to create new session, ID: ", ID, " error: ", err)
	}
	return token, err
}

func GetSession(token string) (login string, err error) {
	var obj cookieSession
	err = sessionCollection.FindOne(context.TODO(), bson.D{{"token", token}}).Decode(&obj)
	if err == nil {
		log.Println("Get session, token: ", token, " id: ", obj.ID)
	} else {
		log.Println("Error to get session token:,", token, " error: ", err)
	}
	return obj.ID, err
}

func getNewSessionToken() string {
	id := ""
	s := rand.Int31()%15 + 15
	for x := int32(0); x < s; x++ {
		id += string(rand.Int31()%25 + 65)
	}
	return id
}

func RecreateCookie(newToken string, ctx echo.Context) *http.Cookie {
	oldToken, err := ctx.Cookie("sessionToken")
	if err == nil && oldToken != nil && oldToken.Value != "" {
		_, _ = DeleteSession(oldToken.Value)
	}
	cookie := new(http.Cookie)
	cookie.Name = "sessionToken"
	cookie.Value = newToken
	cookie.Expires = time.Now().Add(time.Hour)
	return cookie
}
