package main

const (
	ImgDir = "go/images"
)


type Item struct {
	ID     int `json:"id"`
	Name     string `json:"name"`
	Category_id int `json:"category"`
	imageFilename string `json:"img"`
}


type Items struct {
	Items []Item `json:"items"`
}


type Category struct {
	ID     int `json:"id"`
	Name     string `json:"name"`
}


type Categories struct {
	Categories []Category `json:"Categories"`
}


type Response struct {
	Message string `json:"message"`
}
