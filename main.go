package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main/handlers"
	"main/systemComponents/deeds"
	"main/systemComponents/sessions"
	"main/systemComponents/users"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixMilli())
	mongoClient := initMongoDB("127.0.0.1")
	if err := sessions.Init(mongoClient); err != nil {
		return
	}
	log.Println("Session module init")
	if err := users.Init(mongoClient); err != nil {
		return
	}
	log.Println("Authentication module init")
	if err := deeds.Init(mongoClient); err != nil {
		return
	}
	log.Println("Authentication module init")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/auth/users/:ID", handlers.GetUser)
	e.POST("/auth/users", handlers.RegistrateUser)
	e.POST("/auth/login", handlers.AuthenticateUser)
	e.PUT("/auth/users/:ID", handlers.UpdateUser)
	e.GET("/deeds/:ID", handlers.GetDeedByID)
	e.GET("/deeds/user/:ID", handlers.GetDeedByUser)
	e.GET("/deeds", handlers.GetAllDeed)
	e.POST("/deeds", handlers.CreateDeed)

	e.Logger.Fatal(e.Start(":1500"))
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
