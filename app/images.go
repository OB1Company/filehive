package app

import (
	"encoding/base64"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"os"
	"strings"
)

func saveAvatar(imagePath string, base64ImageData string) error {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64ImageData))
	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		return ErrInvalidImage
	}
	width, height := getImageAttributes(150, 150, img.Bounds().Max.X, img.Bounds().Max.Y)
	newImg := imaging.Resize(img, width, height, imaging.Lanczos)

	f, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer f.Close()

	q := &jpeg.Options{Quality: 100}
	return jpeg.Encode(f, newImg, q)
}

func getImageAttributes(targetWidth, targetHeight, imgWidth, imgHeight int) (width, height int) {
	targetRatio := float32(targetWidth) / float32(targetHeight)
	imageRatio := float32(imgWidth) / float32(imgHeight)
	var h, w float32
	if imageRatio > targetRatio {
		h = float32(targetHeight)
		w = float32(targetHeight) * imageRatio
	} else {
		w = float32(targetWidth)
		h = float32(targetWidth) * (float32(imgHeight) / float32(imgWidth))
	}
	return int(w), int(h)
}
