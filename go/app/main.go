package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"path"
	"path/filepath"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)


const (
	ImgDir = "go/images"
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


// load of db
func loadItemsFromDB() ([]Item, error) {
	// Open the db
	db, err := sql.Open("sqlite3", "db/mercari.sqlite3")
	if err != nil {
			return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, category, image_name FROM items")
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.imageFilename); err != nil {
					return nil, err
			}
			items = append(items, item)
	}
	return items, nil
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
	items, err := loadItemsFromDB()
	if err != nil {
			return err
	}
	return c.JSON(http.StatusOK, Items{Items: items})
}


// Get item by ID
func getItem(c echo.Context) error {
	// Get id from URL
	id := c.Param("id")

	// Get item list
	items, err := loadItemsFromDB()
	if err != nil {
			return err
	}

	// Find the item matching id
	for _, item := range items {
		if item.ID == id {
					return c.JSON(http.StatusOK, item)
			}
	}

	// If the item is not found
	return c.JSON(http.StatusNotFound, Response{Message: "Item not found"})
}


// Check if the item ID is unique
func isItemIDUnique(id string) (bool, error) {
	items, err := loadItemsFromDB()
	if err != nil {
		return false, err
	}

	for _, item := range items {
		if item.ID == id {
			return false, nil
		}
	}

	return true, nil
}


// Search products containing keywords from db
func searchItems(c echo.Context) error {
	// Get keyword from URL
	keyword := c.QueryParam("keyword")

	db, err := sql.Open("sqlite3", "db/mercari.sqlite3")
	if err != nil {
			return err
	}
	defer db.Close()

	// LIKE to search for products that contain keywords in the name
	query := "SELECT id, name, category, image_name FROM items WHERE name LIKE ?"
	rows, err := db.Query(query, "%"+keyword+"%")
	if err != nil {
			return err
	}
	defer rows.Close()

	// Stores data in the Item structure (if multiple hits are received, they are all grouped together as items)
	var items []Item
	for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.imageFilename); err != nil {
					return err
			}
			items = append(items, item)
	}
	return c.JSON(http.StatusOK, Items{Items: items})
}


func addItem(c echo.Context) error {
	id := c.FormValue("id")
	isUnique, err := isItemIDUnique(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Message: "Internal server error"})
	}
	if !isUnique {
		return c.JSON(http.StatusBadRequest, Response{Message: "Item ID is not unique"})
	}

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

	newItem := Item{ID: id, Name: name, Category: category, imageFilename: img_name}

	// Open the db
	db, err := sql.Open("sqlite3", "db/mercari.sqlite3")
	if err != nil {
			return err
	}
	defer db.Close()

	// Add new items to the db
	_, err = db.Exec("INSERT INTO items (id, name, category, image_name) VALUES (?, ?, ?, ?)",
			newItem.ID, newItem.Name, newItem.Category, newItem.imageFilename)
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
	if err := os.Chdir("../"); err != nil {
		log.Fatalf("Failed to change current directory: %v", err)
	}
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
	e.GET("/image/:imageFilename", getImg)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
