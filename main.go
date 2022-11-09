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

}

//swagger:operation DELETE /recipes/{id} recipes deleteRecipe
//Delete an existing recipe
//---
//parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
//produces:
// - application/json
//responses:
//	'200':
//		description: Successful operation
//	'400':
//		description: Invalid input
//	'404':
//		description: Invalid recipe ID
/*
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

//swagger:operation GET /recipes/search recipes searchRecipes
//Search recipe by tag
//---
//produces:
//- application/json
//responses:
//  '200':
//		description: Successful operation

func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}
*/
func main() {
	router := gin.Default()
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes", recipesHandler.UpdateRecipeHandler)
	/*	router.DELETE("/recipes", DeleteRecipeHandler)
		router.GET("/recipes/search", SearchRecipesHandler)*/
	router.Run()
}
