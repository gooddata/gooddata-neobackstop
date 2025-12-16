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

// DimensionDiff holds dimension difference information
type DimensionDiff struct {
	WidthDiff  int // positive = test wider, negative = test narrower
	HeightDiff int // positive = test taller, negative = test shorter
}

// CompareDimensions checks if two images have the same dimensions
func CompareDimensions(img1, img2 image.Image) (same bool, diff DimensionDiff) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	w1, h1 := b1.Dx(), b1.Dy()
	w2, h2 := b2.Dx(), b2.Dy()

	diff = DimensionDiff{
		WidthDiff:  w2 - w1,
		HeightDiff: h2 - h1,
	}

	same = (w1 == w2) && (h1 == h2)
	return
}

// DiffImagesPink compares two images and creates a diff image with pink highlights
// Only compares the overlapping region (intersection of both bounds)
func DiffImagesPink(img1, img2 image.Image) (diff *image.RGBA, mismatch float64) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	// Calculate intersection bounds (only compare overlapping region)
	minX := max(b1.Min.X, b2.Min.X)
	minY := max(b1.Min.Y, b2.Min.Y)
	maxX := min(b1.Max.X, b2.Max.X)
	maxY := min(b1.Max.Y, b2.Max.Y)

	// If no overlap, return empty diff
	if minX >= maxX || minY >= maxY {
		diff = image.NewRGBA(image.Rect(0, 0, 1, 1))
		mismatch = 100
		return
	}

	intersectBounds := image.Rect(minX, minY, maxX, maxY)
	diff = image.NewRGBA(intersectBounds)

	var totalPixels int
	var diffPixels int

	pink := color.RGBA{R: 255, G: 0, B: 255, A: 255} // classic Backstop pink
	tolerance := 5.0                                 // fixed tolerance

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
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

// DiffImagesPinkWithBounds compares two images and highlights out-of-bounds areas
// The diff image covers the union of both image bounds
func DiffImagesPinkWithBounds(img1, img2 image.Image) (diff *image.RGBA, mismatch float64) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	// Calculate union bounds (covers both images)
	minX := min(b1.Min.X, b2.Min.X)
	minY := min(b1.Min.Y, b2.Min.Y)
	maxX := max(b1.Max.X, b2.Max.X)
	maxY := max(b1.Max.Y, b2.Max.Y)

	unionBounds := image.Rect(minX, minY, maxX, maxY)
	diff = image.NewRGBA(unionBounds)

	var totalPixels int
	var diffPixels int

	pink := color.RGBA{R: 255, G: 0, B: 255, A: 255}   // classic Backstop pink for pixel diff
	yellow := color.RGBA{R: 255, G: 255, B: 0, A: 255} // yellow for out-of-bounds areas
	tolerance := 5.0

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			inImg1 := x >= b1.Min.X && x < b1.Max.X && y >= b1.Min.Y && y < b1.Max.Y
			inImg2 := x >= b2.Min.X && x < b2.Max.X && y >= b2.Min.Y && y < b2.Max.Y

			if !inImg1 || !inImg2 {
				// Out of bounds for one of the images - mark as yellow
				diff.Set(x, y, yellow)
				diffPixels++
			} else {
				// Both images have this pixel - compare normally
				c1 := color.RGBAModel.Convert(img1.At(x, y)).(color.RGBA)
				c2 := color.RGBAModel.Convert(img2.At(x, y)).(color.RGBA)

				dr := float64(c1.R) - float64(c2.R)
				dg := float64(c1.G) - float64(c2.G)
				db := float64(c1.B) - float64(c2.B)

				dist := math.Sqrt(dr*dr + dg*dg + db*db)

				if dist > tolerance {
					diff.Set(x, y, pink)
					diffPixels++
				} else {
					diff.Set(x, y, color.RGBA{})
				}
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
