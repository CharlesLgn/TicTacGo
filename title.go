package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"github.com/disintegration/gift"
	"image/png"
	"net/http"
	"os"
)

var pictureWidth = 363
var pictureHeight = 82

var titlePicture image.Image
var titleRect = image.Rect(0, 0, 362, 81)

type Picture struct {
	size     image.Rectangle // the sprite size
	Filter   *gift.GIFT      // normal filter used to draw the sprite
}

var pict  = Picture{
	size:titleRect,
	Filter:gift.New(gift.Crop(titleRect)),
}

// get the game frames
func getTitleImage(w http.ResponseWriter, r *http.Request) {
	str := "data:image/png;base64," + title
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(str))
}

func generateTittle()  {
	dst := image.NewRGBA(image.Rect(0, 0, pictureWidth, pictureHeight))
	pict.Filter.DrawAt(dst, titlePicture, image.Pt(0,0), gift.OverOperator)
	createFrame(dst)
}

// create a frame from the image
func createFrame(img image.Image) {
	var buf bytes.Buffer
	png.Encode(&buf, img)
	title = base64.StdEncoding.EncodeToString(buf.Bytes())
}

// get an image from the file
func getImage(filePath string) image.Image {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file:", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file:", err)
	}
	return img
}