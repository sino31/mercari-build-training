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
	Name     string `json:"name"`
	Category string `json:"category"`
	ImageName string `json:"img"`
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
			return err
	}
	return c.JSON(http.StatusOK, Items{Items: items})
}

func addItem(c echo.Context) error {
	name := c.FormValue("name")
	category := c.FormValue("category")

	// Receive image files
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	c.Logger().Infof("Receive item: %s", name)

	// Open file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read file and calculate hash value
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return err
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	// Generate file names from hash values
	img_name := hashString + ".jpg"

	// Save images in the images directory
	dst, err := os.Create(filepath.Join(ImgDir, img_name))
	if err != nil {
		return err
	}
	defer dst.Close()

	// move the file pointer back to the beginning
	src.Seek(0, io.SeekStart)
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	newItem := Item{Name: name, Category: category, ImageName:img_name}

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
		return err
	}

	// Encode updated item list to JSON
	err = ioutil.WriteFile(ItemsFile, updatedData, 0644)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("item received: %s", name)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}

func getImg(c echo.Context) error {
	// Create image path
	imgPath := path.Join(ImgDir, c.Param("imageFilename"))

	if !strings.HasSuffix(imgPath, ".jpg") {
		res := Response{Message: "Image path does not end with .jpg"}
		return c.JSON(http.StatusBadRequest, res)
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
	e.Logger.SetLevel(log.INFO)

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
	e.POST("/items", addItem)
	e.GET("/image/:imageFilename", getImg)


	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
