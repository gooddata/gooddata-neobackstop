package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
)

// this was written by AI

// copyFile copies a single file from src to dst
func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Copy file permissions
	info, err := os.Stat(srcFile)
	if err != nil {
		return err
	}
	return os.Chmod(dstFile, info.Mode())
}

// CopyDir recursively copies a folder from src to dst
func CopyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute the target path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// LoadImage loads an image from a file path
func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// DecodeImageFromBytes decodes an image from a byte slice
func DecodeImageFromBytes(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	return img, err
}

// DiffImagesPink compares two images and creates a diff image with pink highlights
func DiffImagesPink(img1, img2 image.Image) (diff *image.RGBA, mismatch float64) {
	bounds := img1.Bounds()
	diff = image.NewRGBA(bounds)

	var totalPixels int
	var diffPixels int

	pink := color.RGBA{R: 255, G: 0, B: 255, A: 255} // classic Backstop pink
	tolerance := 5.0                                 // fixed tolerance

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := color.RGBAModel.Convert(img1.At(x, y)).(color.RGBA)
			c2 := color.RGBAModel.Convert(img2.At(x, y)).(color.RGBA)

			dr := float64(c1.R) - float64(c2.R)
			dg := float64(c1.G) - float64(c2.G)
			db := float64(c1.B) - float64(c2.B)

			dist := math.Sqrt(dr*dr + dg*dg + db*db)

			if dist > tolerance {
				diff.Set(x, y, pink) // mark difference
				diffPixels++
			} else {
				diff.Set(x, y, color.RGBA{}) // transparent if no diff
			}
			totalPixels++
		}
	}

	mismatch = float64(diffPixels) / float64(totalPixels) * 100
	return
}

// SaveImage saves an image to a file path
func SaveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
