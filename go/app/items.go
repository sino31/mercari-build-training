package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"path/filepath"
	"strconv"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "github.com/labstack/gommon/log"
)


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
	idStr := c.Param("id")

	items, err := loadItemsFromDB()
	if err != nil {
			return err
	}

	// Find the item matching id
	for _, item := range items {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Response{Message: "Invalid ID format"})
		}
		if item.ID == id {
					return c.JSON(http.StatusOK, item)
			}
	}
	return c.JSON(http.StatusNotFound, Response{Message: "Item not found"})
}


// Search products containing keywords from db
func searchItems(c echo.Context) error {
	keyword := c.QueryParam("keyword")

	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
			return err
	}
	defer db.Close()

	// LIKE to search for products that contain keywords in the name
	query := "SELECT id, name, category_id, image_name FROM items WHERE name LIKE ?"
	rows, err := db.Query(query, "%"+keyword+"%")
	if err != nil {
			return err
	}
	defer rows.Close()

	// Stores data in the Item structure (if multiple hits are received, they are all grouped together as items)
	var items []Item
	for rows.Next() {
			var db_item DBItem
			if err := rows.Scan(&db_item.ID, &db_item.Name, &db_item.Category_id, &db_item.Image_name); err != nil {
					return err
			}
			category_name, err := getCategoryName(db_item.Category_id)
			if err != nil {
				return err
		}
		item := convertToItem(db_item, category_name)
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, Items{Items: items})
}


// Add item
func addItem(c echo.Context) error {
	name := c.FormValue("name")
	category := c.FormValue("category")
	category_id, err := getCategoryID(category);
	if err != nil {
		return err
	}
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
	hashString, err := GenerateImageHash(src)
	if err != nil {
		return err
	}
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

	// Open the db
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
			return err
	}
	defer db.Close()

	// Add new items to the db
	_, err = db.Exec("INSERT INTO items (name, category_id, image_name) VALUES (?, ?, ?)", name, category_id, img_name)
	if err != nil {
			return err
	}

	message := fmt.Sprintf("item received: %s", name)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}

