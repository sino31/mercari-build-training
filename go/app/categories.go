package main

import (
	"fmt"
	"net/http"
	"strconv"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	// "github.com/labstack/gommon/log"
)


// Get Categories List
func getCategories(c echo.Context) error {
	categories, err := loadCategoriesFromDB()
	if err != nil {
			return err
	}
	return c.JSON(http.StatusOK, Categories{Categories: categories})
}


// Get Category by ID
func getCategory(c echo.Context) error {
	// Get id from URL
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
			return c.JSON(http.StatusBadRequest, Response{Message: "Invalid ID format"})
	}

	// Get Category list
	Categories, err := loadCategoriesFromDB()
	if err != nil {
			return err
	}

	// Find the Category matching id
	for _, Category := range Categories {
		if Category.ID == id {
					return c.JSON(http.StatusOK, Category)
			}
	}

	// If the Category is not found
	return c.JSON(http.StatusNotFound, Response{Message: "Category not found"})
}


// Create a categories
func addCategory(c echo.Context) error {
	idStr := c.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
			return c.JSON(http.StatusBadRequest, Response{Message: "Invalid ID format"})
	}

	name := c.FormValue("name")

	newItem := Item{ID: id, Name: name}

	// Open the db
	db, err := sql.Open("sqlite3", "db/mercari.sqlite3")
	if err != nil {
			return err
	}
	defer db.Close()

	// Add new category to the db
	_, err = db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", newItem.ID, newItem.Name)
	if err != nil {
			return err
	}

	message := fmt.Sprintf("category received: %s", name)
	res := Response{Message: message}

	return c.JSON(http.StatusOK, res)
}
