package authentication

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"main/sessions"
)

type User struct {
	Login      string `bson:"login"`
	Password   string `bson:"password"`
	Email      string `bson:"email"`
	FirstName  string `bson:"firstName"`
	SecondName string `bson:"secondName"`
	Age        uint32 `bson:"age"`
	IsAdmin    bool   `bson:"isAdmin"`
}

var dbClient *mongo.Client
var usersCollection *mongo.Collection

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
	//Проверка на повтор данных регистрации в БД
	filter := bson.M{"$or": []bson.M{{"login": user.Login}, {"email": user.Email}}}
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
	if token, err := sessions.CreateSessions(user.Login); err != nil {
		return "", err
	} else {
		return token, err
	}

}

// Authentication if User exist in DataBase then return new session
func (user *User) Authentication() (userData User, token string, err error) {
	filter := bson.M{"$and": []bson.M{{"login": user.Login}, {"password": user.Password}}}
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&userData); err != nil {
		log.Println("Can't authenticate error: ", err)
		return User{}, "", err
	}
	if token, err := sessions.CreateSessions(userData.Login); err == nil {
		log.Println("Success authentication, login:", userData.Login, "token: ", token)
		return userData, token, err
	}

	log.Println("Success authentication, unsuccessful session:", userData.Login, "error: ", err)
	return userData, "", err
}

func GetUserByLogin(login string) (user User, err error) {
	filter := bson.M{"login": login}
	if err := usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		log.Println("Can't get user data error: ", err)
		return User{}, err
	}
	return user, err
}
