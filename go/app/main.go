package main

import (
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)


func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}


func main() {
	// if err := os.Chdir("../"); err != nil {
	// 	log.Fatalf("Failed to change current directory: %v", err)
	// }
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Set log level (Only messages above the ERROR level are recorded in the production environment)
	e.Logger.SetLevel(log.DEBUG)

	frontURL := os.Getenv("FRONT_URL")
	if frontURL == "" {
		frontURL = "http://localhost:3000"
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontURL},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.GET("/", root)
	e.GET("/items", getItems)
	e.GET("/items/:id", getItem)
	e.GET("/search", searchItems)
	e.POST("/items", addItem)
	e.POST("/categories", addCategory)
	e.GET("/categories", getCategories)
	e.GET("/categories/:id", getCategory)
	e.GET("/image/:imageFilename", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
