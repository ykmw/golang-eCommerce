package main

import(
	"example.com/m/ecommerce/controllers"
	"example.com/m/ecommerce/database"
	"example.com/m/ecommerce/middleware"
	"example.com/m/ecommerce/routes"
	"example.com/gin-gonic/gin"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

	app := controllers.NetApplication(database.ProductData(database.Client, "Products"), databaseUserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}