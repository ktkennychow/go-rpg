package assets

import (
	"embed"
	"image"
	_ "image/png" // Import PNG decoder
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *
var FS embed.FS

// MustLoadImage loads an image from embedded assets and panics on error
func MustLoadImage(name string) *ebiten.Image {
	f, err := FS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// LoadImage loads an image from embedded assets and returns an error
func LoadImage(name string) (*ebiten.Image, error) {
	f, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

func MustLoadFile(name string) []byte {
	f, err := FS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	return data
}

// LoadFile loads a file from embedded assets
func LoadFile(name string) ([]byte, error) {
	f, err := FS.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}
