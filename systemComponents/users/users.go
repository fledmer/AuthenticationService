package users

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"main/systemComponents"
	"main/systemComponents/sessions"
	"strconv"
)

var dbClient *mongo.Client
var usersCollection *mongo.Collection

type User systemComponents.User

func Init(client *mongo.Client) error {
	if client == nil {
		return errors.New("nil DbClient")
	}
	dbClient = client
	usersCollection = dbClient.Database("test").Collection("users")
	return nil
}

// Registration return (true,nil) if user has been added to DataBase;
// if user has not added to DataBase, return false, error
func (user *User) Registration() (token string, err error) {
	//Создаем ID если пользователь не решил его взять.
	if user.ID == "" {
		if numID, err := usersCollection.CountDocuments(context.TODO(), bson.M{}); err != nil {
			log.Println("Can't get new ID: ", err)
			return "", err
		} else {
			user.ID = strconv.Itoa(int(numID))
		}
	}
	//Проверка на повтор данных регистрации в БД
	filter := bson.M{"$or": []bson.M{{"mail": user.Mail}, {"ID": user.ID}}}
	if count, err := usersCollection.CountDocuments(context.TODO(), filter); count > 0 {
		return "", errors.New("already registered")
	} else if err != nil {
		return "", err
	}
	//Вносим
	if _, err := usersCollection.InsertOne(context.TODO(), user); err != nil {
		return "", err
	}
	//Создаем сессию
	if token, err := sessions.CreateSessions(user.ID); err != nil {
		return "", err
	} else {
		return token, err
	}

}

// Authentication if User exist in DataBase then return new session
func (user *User) Authentication() (userData User, token string, err error) {
	filter := bson.M{"$and": []bson.M{{"mail": user.Mail}, {"password": user.Password}}}
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&userData); err != nil {
		log.Println("Can't authenticate error: ", err)
		return User{}, "", err
	}
	if token, err := sessions.CreateSessions(userData.ID); err == nil {
		log.Println("Success users. ID:", userData.ID, "token: ", token)
		return userData, token, err
	}

	log.Println("Success users, unsuccessful session. ID:", userData.ID, "error: ", err)
	return userData, "", err
}

func GetUserByID(ID string) (user User, err error) {
	filter := bson.M{"ID": ID}
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		log.Println("Can't get userData ", "ID: ", ID, " error: ", err)
		return User{}, err
	}
	return user, err
}

func (user *User) Update() error {
	filter := bson.M{"ID": user.ID}
	if result, err := usersCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": user}); err != nil {
		log.Println("Can't update user. ID: ", user.ID, " error: ", err)
		return err
	} else if result.MatchedCount == 0 {
		log.Println("Can't update user. ID: ", user.ID, " error: ", "user not found")
		return errors.New("user not found")
	}
	return nil
}
