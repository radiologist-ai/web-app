package image_processing

import (
	"bytes"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

var BackgroundColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}

func cropRectangle(img image.Image, expectedWidth, expectedHeight int) image.Image {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	expectedCoeff := float32(expectedWidth) / float32(expectedHeight)
	coeff := float32(width) / float32(height)
	if coeff == expectedCoeff {
		return img
	}
	var croppedImage *image.RGBA
	var dp image.Point

	if coeff > expectedCoeff {
		newHeight := int(float32(width) / expectedCoeff)
		croppedImage = image.NewRGBA(image.Rect(0, 0, width, newHeight))
		draw.Draw(croppedImage, croppedImage.Bounds(), &image.Uniform{C: BackgroundColor}, image.Point{}, draw.Src)
		dp = image.Point{
			X: 0,
			Y: (newHeight - height) / 2,
		}
	} else {
		newWidth := int(float32(height) * expectedCoeff)
		croppedImage = image.NewRGBA(image.Rect(0, 0, newWidth, height))
		draw.Draw(croppedImage, croppedImage.Bounds(), &image.Uniform{C: BackgroundColor}, image.Point{}, draw.Src)
		dp = image.Point{
			X: (newWidth - width) / 2,
			Y: 0,
		}
	}

	r := image.Rectangle{Min: dp, Max: dp.Add(img.Bounds().Size())}
	draw.Draw(croppedImage, r, img, img.Bounds().Min, draw.Src)
	return croppedImage
}

func resize(img image.Image, widthPx, heightPx int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, widthPx, heightPx))
	draw.ApproxBiLinear.Scale(resized, resized.Rect, img, img.Bounds(), draw.Over, nil)
	return resized
}

func ProcessImage(src image.Image, widthPx, heightPx int, format string) (*bytes.Reader, error) {
	res := cropRectangle(src, widthPx, heightPx)
	res = resize(res, widthPx, heightPx)

	dst, err := imageToReader(res, format)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func imageToReader(src image.Image, format string) (*bytes.Reader, error) {
	dst := new(bytes.Buffer)
	var err error
	if format == "png" {
		err = png.Encode(dst, src)
	} else {
		err = jpeg.Encode(dst, src, nil)
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(dst.Bytes()), nil
}
