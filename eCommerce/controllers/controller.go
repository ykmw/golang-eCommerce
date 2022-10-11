package controllers

import(
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string{
	bcrypt.GenertateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)

}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userPassword))
	valid := true
	msg :=""

	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err!= nil{
			c.JSON{http.StatusBadRequest, gin.H{"error": err.Error()}}
			return
		}

		validationErr := Validata.Struct(user)
		if validationErr != nil {
			c.JSON{http.StatusBadRequest, gin.H{"error": validationErr()}}
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone":user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err})
			return
		}

		if count>0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"this phone no. already in use"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC33339))
		user.Updated_At,
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!" 

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
	}

	var user models.User
	if err := c.BindJSON(&user); err!= nil{
		c.JSON{http.StatusBadRequest, gin.H{"error": err}}
		return
	}

	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
	defer cancel()

	if err !=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
		return
	}

	PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)

	if !PasswordIsValid{
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.PrintIn(msg)
		return
	}

	token, refreshToken, _ := generate.TokenGenerator(*founderuser.Email, *founderuser.First_Name. *founderuser.Last_Name, *founduser.User_ID)
	defer cancel()

	generate.UpdateAllTokens(token, refreshToken, founderuser.User_ID)

	c.JSON(http.StatusFound, founduser)
}

func ProductViewerAdmin() gin.HandlerFunc {

}

func SearchProduct() gin.HandlerFunc {

	return func(c *gin.Context){

		var productlist []models.Product
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err!= nil {
			c.IndentedJSON(http.StatusInteranlServerError, "something went wrong, please try again later")
			return
		}

		err = cursor.All(ctx, &productlist)

		if err!= nil{
			log.PrintIn(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close()

		if err := cursor.err(); err != nil{
			log.PrintIn(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer cancel()
		c.IndentedJSON(200, productlist)
	}

}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context){
		var searchProducts []models.Product
		queryParam := c.Query("name")

		//you want to check if it's empty

		if queryParam == ""{
			log.PrintIn("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})

		if err != nil {
			c.IndentedJSON(404, "something went wrong while fetching the data")
			return
		}

		err = searchquerydb.All(ctx, &searchproducts)
		if err != nil {
			log.PrintIn(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchquerydb.Close(ctx)
		
		if err := searchquerydb.Err(); err != nil
			log.PrintIn(err)
			c.IndentedJSON(400, "invalid request")
			return
	}

	defer cancel()
	c.IndentedJSON(200, searchproducts)

}