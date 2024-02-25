package main

const (
	DbPath="../db/mercari.sqlite3"
	ImgDir="images"
)

type Item struct {
	ID     int `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Image_name string `json:"image_name"`
}

type DBItem struct {
	ID     int
	Name     string
	Category_id int
	Image_name string
}

type Items struct {
	Items []Item `json:"items"`
}


type Category struct {
	ID     int `json:"id"`
	Name     string `json:"name"`
}


type Categories struct {
	Categories []Category `json:"categories"`
}


type Response struct {
	Message string `json:"message"`
}
