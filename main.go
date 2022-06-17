package main

import (
	_ "golang.org/x/image/webp"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func createItem(filePath string, targetSize int, radius int) image.Image {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f) // img, formatName, err
	if err != nil {
		panic(err)
	}
	sizeX := img.Bounds().Dx()
	sizeY := img.Bounds().Dy()
	scale := float64(targetSize) / float64(max(sizeX, sizeY))
	tSizeX := int(scale * float64(sizeX))
	tSizeY := int(scale * float64(sizeY))
	resizedImg := resizeImg(img, tSizeX, tSizeY)
	resImg := borderRadius(resizedImg, radius)
	// get the beginning drawing point, 我也不清楚为什么要负号, 这坐标轴不常规
	startingPointX := -(targetSize/2 - tSizeX/2)
	startingPointY := -(targetSize/2 - tSizeY/2)

	dst := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	draw.Draw(dst, dst.Bounds(), resImg, image.Pt(startingPointX, startingPointY), draw.Src)
	return dst
}

func main() {
	a := createItem("test/1.webp", 220, 20)
	saveImg(a, "test/result.png")
}
