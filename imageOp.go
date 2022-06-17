package main

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type radius struct {
	p image.Point // 矩形右下角位置
	r int
}

func nearestNeighbor(src *image.RGBA, width, height int) *image.RGBA {
	// simple but the fastest image createItem algorithm
	srcW, srcH := src.Bounds().Dx(), src.Bounds().Dy()
	srcStride := src.Stride

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	dstStride := dst.Stride

	dx := float64(srcW) / float64(width)
	dy := float64(srcH) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := y*dstStride + x*4
			ipos := int((float64(y)+0.5)*dy)*srcStride + int((float64(x)+0.5)*dx)*4

			dst.Pix[pos+0] = src.Pix[ipos+0]
			dst.Pix[pos+1] = src.Pix[ipos+1]
			dst.Pix[pos+2] = src.Pix[ipos+2]
			dst.Pix[pos+3] = src.Pix[ipos+3]
		}
	}

	return dst
}

func imageToRGBA(src image.Image) *image.RGBA {

	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func resizeImg(img image.Image, width, height int) *image.RGBA {
	var dst *image.RGBA
	dst = nearestNeighbor(imageToRGBA(img), width, height)
	return dst
}

func borderRadius(img image.Image, r int) image.Image {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	c := radius{p: image.Point{X: w, Y: h}, r: int(r)}
	radiusImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.DrawMask(radiusImg, radiusImg.Bounds(), img, image.Point{}, &c, image.Point{}, draw.Over)
	return radiusImg
}

func (c *radius) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *radius) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.p.X, c.p.Y)
}

func (c *radius) At(x, y int) color.Color {
	// 对每个像素点进行色值设置，分别处理矩形的四个角，在四个角的内切圆的外侧，色值设置为全透明，其他区域不透明
	var xx, yy, rr float64
	var inArea bool
	// left up
	if x <= c.r && y <= c.r {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-c.r)+0.5, float64(c.r)
		inArea = true
	}
	// right up
	if x >= (c.p.X-c.r) && y <= c.r {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-c.r)+0.5, float64(c.r)
		inArea = true
	}
	// left bottom
	if x <= c.r && y >= (c.p.Y-c.r) {
		xx, yy, rr = float64(c.r-x)+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
		inArea = true
	}
	// right bottom
	if x >= (c.p.X-c.r) && y >= (c.p.Y-c.r) {
		xx, yy, rr = float64(x-(c.p.X-c.r))+0.5, float64(y-(c.p.Y-c.r))+0.5, float64(c.r)
		inArea = true
	}
	if inArea && xx*xx+yy*yy >= rr*rr {
		return color.Alpha{}
	}
	return color.Alpha{A: 255}
}

func saveImg(img image.Image, filePath string) {
	// save image in png format
	outFile, err := os.Create(filePath)
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			panic(err)
		}
	}(outFile)
	if err != nil {
		panic(err)
	}
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		panic(err)
	}
	err = b.Flush()
	if err != nil {
		panic(err)
	}
}
