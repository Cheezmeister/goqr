package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func main() {
	// Define command line flags
	uri := flag.String("uri", "", "URI to encode in the QR code (required)")
	logo := flag.String("logo", "", "Path to logo image (SVG or PNG) to place in center (optional)")
	output := flag.String("output", "qrcode.png", "Output filename")
	size := flag.Int("size", 256, "Size of QR code in pixels")
	flag.Parse()

	// Check for required URI
	if *uri == "" {
		fmt.Println("Error: URI is required")
		flag.Usage()
		os.Exit(1)
	}

	// Generate QR code
	qrCode, err := GenerateQRCode(*uri, *size)
	if err != nil {
		fmt.Printf("Error generating QR code: %v\n", err)
		os.Exit(1)
	}

	// If logo is provided, overlay it
	if *logo != "" {
		qrCode, err = OverlayLogo(qrCode, *logo)
		if err != nil {
			fmt.Printf("Error overlaying logo: %v\n", err)
			os.Exit(1)
		}
	}

	// Save the final QR code
	if err := SaveQRCode(qrCode, *output); err != nil {
		fmt.Printf("Error saving QR code: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("QR code saved to %s\n", *output)
}

// GenerateQRCode creates a QR code from the given URI with specified size
func GenerateQRCode(uri string, size int) (image.Image, error) {
	// Create QR code with high error correction to allow for logo overlay
	qrCode, err := qr.Encode(uri, qr.H, qr.Auto)
	if err != nil {
		return nil, err
	}

	// Scale QR code to requested size
	qrCode, err = barcode.Scale(qrCode, size, size)
	if err != nil {
		return nil, err
	}

	return qrCode, nil
}

// OverlayLogo places a logo in the center of the QR code
func OverlayLogo(qrCode image.Image, logoPath string) (image.Image, error) {
	// Load the logo based on file extension
	var logoImg image.Image
	var err error

	ext := strings.ToLower(filepath.Ext(logoPath))
	switch ext {
	case ".png":
		logoImg, err = loadPNG(logoPath)
	case ".svg":
		logoImg, err = loadSVG(logoPath, qrCode.Bounds().Dx()/4)
	default:
		return nil, errors.New("unsupported logo format: use PNG or SVG")
	}

	if err != nil {
		return nil, err
	}

	// Determine logo size (about 20% of QR code)
	qrSize := qrCode.Bounds().Dx()
	logoSize := qrSize / 5

	// Create a new image to draw on
	result := image.NewRGBA(qrCode.Bounds())
	draw.Draw(result, result.Bounds(), qrCode, image.Point{0, 0}, draw.Src)

	// Calculate position to center the logo
	logoPos := image.Point{
		X: (qrSize - logoSize) / 2,
		Y: (qrSize - logoSize) / 2,
	}

	// Draw the logo
	draw.Draw(result, image.Rect(
		logoPos.X, logoPos.Y,
		logoPos.X+logoSize, logoPos.Y+logoSize),
		logoImg, image.Point{0, 0}, draw.Over)

	return result, nil
}

// loadPNG loads a PNG image from the given path
func loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// loadSVG loads and rasterizes an SVG image from the given path
func loadSVG(path string, size int) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the SVG
	icon, err := oksvg.ReadIconStream(file)
	if err != nil {
		return nil, err
	}

	// Set size
	icon.SetTarget(0, 0, float64(size), float64(size))

	// Create rasterizer
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))
	scanner := rasterx.NewScannerGV(size, size, rgba, rgba.Bounds())
	raster := rasterx.NewDasher(size, size, scanner)

	// Rasterize the SVG
	icon.Draw(raster, 1.0)

	return rgba, nil
}

// SaveQRCode saves the QR code to the specified path
func SaveQRCode(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
