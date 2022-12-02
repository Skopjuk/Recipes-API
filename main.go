//Recipes API

//This is a sample recipe API.

//Schemes: http
//Host: localhost:8080
//BasePath: /
//Version : 1.0.0
//Contact: Kseniia Skopiuk<ksenia.agag@gmail.com>
//
//Consumes:
// - application/json
//
// Produces:
//  - application/json
//swagger:meta
package main

import (
	"context"
	"github.com/Skopjuk/Recipes-API/handlers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var ctx context.Context
var err error
var client *mongo.Client
var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {
	// -- на данный момент просто необходимый параметр для использования подключения к монго, потом может быть чем-то заполнен
	ctx = context.Background()
	//mongo.Connect усанаввливает подключение к серверу и в качестве параметра принимает функцию
	//обратного вызова, которая срабатывает при установке подключения
	//mongo_uri -- переменная окружения которая содержит данные о том куда подключаться
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	//readpref.Primary значит что мы подключаемся к главной копии бд
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipeHandlers(ctx, collection)
	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
}

func main() {
	router := gin.Default()
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signup", authHandler.SignUpHandler)
	authorized := router.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes", recipesHandler.ListRecipesHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/search/:id", recipesHandler.SearchRecipesHandler)
	}
	router.Run()
}
