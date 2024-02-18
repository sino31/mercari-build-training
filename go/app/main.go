package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"path"
	"path/filepath"
	"strings"
	"encoding/json"
	"io/ioutil"
	"crypto/sha256"
	"encoding/hex"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)


const (
	ImgDir = "images"
	ItemsFile  = "app/items.json"
)


type Item struct {
	ID     string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	imageFilename string `json:"img"`
}


type Items struct {
	Items []Item `json:"items"`
}


// load of items.json
func loadItemsFromFile() ([]Item, error) {
	var items Items
	data, err := ioutil.ReadFile(ItemsFile)
	if err != nil {
			return nil, err
	}
	// Convert data from json to go
	err = json.Unmarshal(data, &items)
	if err != nil {
			return nil, err
	}
	return items.Items, nil
}


type Response struct {
	Message string `json:"message"`
}


func root(c echo.Context) error {
	res := Response{Message: "Hello, world!"}
	return c.JSON(http.StatusOK, res)
}


// Get item List
func getItems(c echo.Context) error {
	items, err := loadItemsFromFile()
	if err != nil {
		c.Logger().Errorf("Failed to load items from file: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to load items."})
	}
	return c.JSON(http.StatusOK, Items{Items: items})
}


// Get item by ID
func getItem(c echo.Context) error {
	// Get item_id from URL
	id := c.Param("id")

	// Get item list
	items, err := loadItemsFromFile()
	if err != nil {
		c.Logger().Errorf("Failed to load items from file: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to load items."})
	}

	// Find the item matching id
	for _, item := range items {
		if item.ID == id {
					return c.JSON(http.StatusOK, item)
			}
	}

	// If the item is not found
	return c.JSON(http.StatusNotFound, map[string]string{"message": "Item not found"})
}


func addItem(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	category := c.FormValue("category")

	// Receive image files
	file, err := c.FormFile("image")
	if err != nil {
		c.Logger().Errorf("Failed to receive the file: %v", err)
    return c.JSON(http.StatusBadRequest, Response{Message: "Failed to receive the file"})
	}
	c.Logger().Infof("Receive item: %s", name)

	// Open file
	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("Failed to open the file: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to open the file"})
	}
	defer src.Close()

	// Read file and calculate hash value
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		c.Logger().Errorf("Failed to calculate the file hash: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to calculate the file hash"})
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	// Generate file names from hash values
	img_name := hashString + ".jpg"

	// Save images in the images directory
	dst, err := os.Create(filepath.Join(ImgDir, img_name))
	if err != nil {
		c.Logger().Errorf("Failed to save the image: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to save the image"})
	}
	defer dst.Close()

	// move the file pointer back to the beginning
	src.Seek(0, io.SeekStart)
	if _, err := io.Copy(dst, src); err != nil {
		c.Logger().Errorf("Failed to save the image: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to save the image"})
	}

	newItem := Item{ID:id, Name: name, Category: category, imageFilename:img_name}

	// Read the current item list from items.json
	var items Items
	data, err := ioutil.ReadFile(ItemsFile)
	if err == nil {
		json.Unmarshal(data, &items)
	}

	//　Add new item to list
	items.Items = append(items.Items, newItem)
	updatedData, err := json.Marshal(items)
	if err != nil {
		c.Logger().Errorf("Failed to update the item list: %v", err)
  	return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to update the item list"})
	}

	// Encode updated item list to JSON
	err = ioutil.WriteFile(ItemsFile, updatedData, 0644)
	if err != nil {
		c.Logger().Errorf("Failed to save the item list: %v", err)
    return c.JSON(http.StatusInternalServerError, Response{Message: "Failed to save the item list"})
	}

	message := fmt.Sprintf("item received: %s", name)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}


func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		return c.JSON(http.StatusBadRequest, Response{Message: "Image path does not end with .jpg"})
	}
	if _, err := os.Stat(imgPath); err != nil {
		c.Logger().Debugf("Image not found: %s", imgPath)
		imgPath = path.Join(ImgDir, "default.jpg")
	}
	return c.File(imgPath)
}

func main() {
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
	e.POST("/items", addItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
