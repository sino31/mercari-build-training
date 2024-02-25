package main

import (
	"net/http"
	"os"
	"io"
	"path"
	"strings"
	"crypto/sha256"
	"encoding/hex"

	_ "github.com/mattn/go-sqlite3"
	"github.com/labstack/echo/v4"
)


// Generate image hash
func GenerateImageHash(src io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}


func getImg(c echo.Context) error {
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


func convertToItem(itemDB DBItem, categoryName string) Item {
	return Item{
			ID:       itemDB.ID,
			Name:     itemDB.Name,
			Category: categoryName,
			Image_name: itemDB.Image_name,
	}
}