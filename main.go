package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"qvickly/database/postgres"
	_ "qvickly/docs"
	"qvickly/router"
)

// @title			Qvickly APIs
// @version		1.0
// @description
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host 3.110.183.54
func main() {
	err := godotenv.Load("pg.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = postgres.InitConn()
	if err != nil {
		panic(err.Error() + ": \nError starting postgres\n")
	}
	defer postgres.CloseConn()

	app := gin.Default()

	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)

	router.Router(app)

	if err = app.Run(":8080"); err != nil {
		panic(err.Error())
	}
}
