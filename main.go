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
	mongoClient := initMongoDB("79.120.10.217")
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

	e.GET("/:user", HandlerUser)
	e.POST("/registration", HandlerRegistration)
	e.POST("/authentication", HandlerAuthentication)

	e.Logger.Fatal(e.Start(":1500"))
}

func HandlerUser(ctx echo.Context) error {
	requestLogin := ctx.Param("user")

	token, err := ctx.Cookie("sessionToken")
	if err != nil || token.Value == "" {
		return ctx.JSON(http.StatusOK, "Message: Авторизуйся, пес")
	}
	login, err := sessions.GetSession(token.Value)
	if err != nil {
		return ctx.JSON(http.StatusOK, "Message: Авторизуйся, пес")
	}

	userData, err := authentication.GetUserByLogin(requestLogin)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Message: Я такого не знаю, ди в хуй")
	}
	if login == requestLogin {
		return ctx.JSON(http.StatusOK, userData)
	} else {
		return ctx.JSON(http.StatusOK, "First Name: "+userData.FirstName)
	}
}

func HandlerAuthentication(ctx echo.Context) error {
	//Берем из тела логин
	bodyData := make(map[string]string)
	if err := json.NewDecoder(ctx.Request().Body).Decode(&bodyData); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Message: can't decode body")
	}
	login, found := bodyData["login"]
	if found != true {
		return ctx.JSON(http.StatusBadRequest, "Message: no Login")
	}
	password, found := bodyData["password"]
	if found != true {
		return ctx.JSON(http.StatusBadRequest, "Message: no Login")
	}

	user := authentication.User{
		Login:    login,
		Password: password,
	}

	if _, token, err := user.Authentication(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		//удаляем старый токен,
		oldToken, _ := ctx.Cookie("sessionToken")
		sessions.DeleteSession(oldToken.Value)

		cookie := new(http.Cookie)
		cookie.Name = "sessionToken"
		cookie.Value = token
		cookie.Expires = time.Now().Add(time.Hour)
		ctx.SetCookie(cookie)
		return ctx.JSON(http.StatusOK, "Message: "+" Вход в аккаунт выполнен!")
	}
}

func HandlerRegistration(ctx echo.Context) error {
	//Берем из тела логин
	bodyData := make(map[string]string)
	if err := json.NewDecoder(ctx.Request().Body).Decode(&bodyData); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Message: can't decode body")
	}
	login, found := bodyData["login"]
	if found != true {
		return ctx.JSON(http.StatusBadRequest, "Message: no Login")
	}
	password, found := bodyData["password"]
	if found != true {
		return ctx.JSON(http.StatusBadRequest, "Message: no Login")
	}

	user := authentication.User{
		Login:      login,
		Password:   password,
		Email:      bodyData["email"],
		FirstName:  bodyData["firstName"],
		SecondName: bodyData["secondName"],
		Age:        0,
		IsAdmin:    false,
	}

	if token, err := user.Registration(); err != nil {
		log.Println("Registration error: ", err)
		return ctx.JSON(http.StatusBadRequest, "Message: "+err.Error())
	} else {
		oldToken, _ := ctx.Cookie("sessionToken")
		sessions.DeleteSession(oldToken.Value)
		cookie := new(http.Cookie)
		cookie.Name = "sessionToken"
		cookie.Value = token
		cookie.Expires = time.Now().Add(time.Hour)
		ctx.SetCookie(cookie)
		return ctx.JSON(http.StatusOK, "Message: "+"You have been register!")
	}
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
