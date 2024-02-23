package main

import (
	"fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)


// load of items db
func loadItemsFromDB() ([]Item, error) {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
			return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, category_id, image_name FROM items")
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Category_id, &item.imageFilename); err != nil {
					return nil, err
			}
			items = append(items, item)
	}
	return items, nil
}


// load of categories db
func loadCategoriesFromDB() ([]Category, error) {
	db, err := sql.Open("sqlite3", DbPath)
	if err != nil {
			return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var Categories []Category
	for rows.Next() {
			var Category Category
			if err := rows.Scan(&Category.ID, &Category.Name);
			err != nil {
					return nil, err
			}
			Categories = append(Categories, Category)
	}
	return Categories, nil
}


// Get category_id from category name
func getCategoryID(category_name string) (int, error) {
	categories, err := loadCategoriesFromDB()
	if err != nil {
			return 0,err
	}
	for _, category := range categories {
		if category.Name == category_name {
				return category.ID, nil
		}
	}
	return 0, fmt.Errorf("Category not found")
}


