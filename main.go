package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main/authentication"
	"main/sessions"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixMilli())
	mongoClient := initMongoDB("127.0.0.1")
	if err := sessions.Init(mongoClient); err != nil {
		return
	}
	log.Println("Session module init")
	if err := authentication.Init(mongoClient); err != nil {
		return
	}
	log.Println("Authentication module init")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/auth/users/:ID", HandlerUser)
	e.POST("/auth/users", HandlerRegistration)
	e.POST("/auth/login", HandlerAuthentication)
	e.PUT("/auth/users/:ID", HandleUpdateUser)

	e.Logger.Fatal(e.Start(":1500"))
}

func HandlerUser(ctx echo.Context) error {
	requestID := ctx.Param("ID")
	//Проверка токена
	token, err := ctx.Cookie("sessionToken")
	if err != nil || token.Value == "" {
		return ctx.JSON(http.StatusUnauthorized, "")
	}
	ID, err := sessions.GetSession(token.Value)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "")
	}
	//Получаем данные о юзере
	userData, err := authentication.GetUserByID(requestID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "")
	}
	if ID == requestID {
		return ctx.JSON(http.StatusOK, userData)
	} else {
		return ctx.JSON(http.StatusOK,
			"First Name: "+userData.FirstName+
				"Second Name: "+userData.LastName)
	}
}

func HandlerAuthentication(ctx echo.Context) error {
	//Берем из тела логин и пароль
	user := authentication.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, "Message: can't decode body")
	}
	if user.Mail == "" || user.Password == "" {
		return ctx.JSON(http.StatusBadRequest, "Message: no email/pass")
	}

	if _, token, err := user.Authentication(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		ctx.SetCookie(sessions.RecreateCookie(token, ctx))
		return ctx.JSON(http.StatusOK, "Message: "+" Вход в аккаунт выполнен!")
	}
}

func HandlerRegistration(ctx echo.Context) error {
	//Берем из тела логин
	user := authentication.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, "Message: can't decode body")
	}

	if user.Mail == "" || user.Password == "" {
		return ctx.JSON(http.StatusBadRequest, "Message: no email/pass")
	}
	if token, err := user.Registration(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		ctx.SetCookie(sessions.RecreateCookie(token, ctx))
		return ctx.JSON(http.StatusOK, "Message: "+"You have been register!")
	}
}

func HandleUpdateUser(ctx echo.Context) error {
	ID := ctx.Param("ID")
	newUserData := authentication.User{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(&newUserData); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusBadRequest, "Message: can't decode body")
	}

	token, err := ctx.Cookie("sessionToken")
	if err != nil || token.Value == "" {
		return ctx.JSON(http.StatusUnauthorized, "")
	}

	if tokenId, err := sessions.GetSession(token.Value); err != nil {
		return ctx.JSON(http.StatusUnauthorized, "")
	} else {
		//У текущего запрос отправителя есть права?
		if tokenId == ID {
			if err = newUserData.Update(); err == nil {
				return ctx.JSON(http.StatusOK, "Message: successful update")
			} else {
				return ctx.JSON(http.StatusNotFound, "Message: "+err.Error())
			}
		}
	}
	return ctx.JSON(http.StatusBadRequest, "Message: wtf?")
}

func initMongoDB(address string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + address + ":27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client
}
