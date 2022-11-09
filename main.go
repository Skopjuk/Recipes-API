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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id"           bson:"_id"`
	Name         string             `json:"name"         bson:"name"`
	Tags         []string           `json:"tags"         bson:"tags"`
	Ingredients  []string           `json:"ingredients"  bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"  bson:"publishedAt"`
}

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

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

	//var listOfRecipes []interface{}
	//for _, recipe := range recipes {
	//	listOfRecipes = append(listOfRecipes, recipe)
	//}
	//collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	//insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Inserted recipes: ",
	//	len(insertManyResult.InsertedIDs))
}

//swagger:operation POST /recipes recipes newRecipe
//Adds new recipe
//---
//produces:
//- application/json
//responses:
//	'200':
//		description: Successful operation
//	'400':
//		description: Invalid input

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting " +
				"a new recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

//swagger:operation GET /recipes recipes listRecipes
//Returns list of recipes
//---
//produces:
//- application/json
//responses:
//	'200':
//		description: Successful operation

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

//swagger:operation PUT /recipes/{id} recipes updateRecipe
//Update an existind recipe
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
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
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
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes", UpdateRecipeHandler)
	/*	router.DELETE("/recipes", DeleteRecipeHandler)
		router.GET("/recipes/search", SearchRecipesHandler)*/
	router.Run()
}
