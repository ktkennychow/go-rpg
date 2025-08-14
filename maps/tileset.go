package maps

import (
	"encoding/json"
	"image"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ktkennychow/go-rpg/assets"
	"github.com/ktkennychow/go-rpg/constants"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type TilesetJSON struct {
	FirstGid int    `json:"firstgid"`
	Source   string `json:"source"`
}
type TileJSON struct {
	Id          int        `json:"id"`
	ImagePath   string     `json:"image"`
	ImageWidth  int        `json:"imagewidth"`
	ImageHeight int        `json:"imageheight"`
	Columns     int        `json:"columns"`
	Tiles       []TileJSON `json:"tiles"`
}

type UniformTileset struct {
	gid     int
	img     *ebiten.Image
	columns int
}

func (uT *UniformTileset) Img(id int) *ebiten.Image {
	// tile id - tileset offset(g(rid)id) = id relative to tileset(grid)
	id -= uT.gid

	srcX := id % uT.columns
	srcY := id / uT.columns

	srcX *= constants.TileSize
	srcY *= constants.TileSize

	return uT.img.SubImage(image.Rect(srcX, srcY, srcX+constants.TileSize, srcY+constants.TileSize)).(*ebiten.Image)
}

type DynamicTileset struct {
	gid  int
	imgs []*ebiten.Image
}

func (dT *DynamicTileset) Img(id int) *ebiten.Image {
	id -= dT.gid

	return dT.imgs[id]
}

func NewTilesetFromPath(filepath string, gid int) (Tileset, error) {
	f, err := assets.LoadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tileJSON TileJSON
	err = json.Unmarshal(f, &tileJSON)
	if err != nil {
		return nil, err
	}

	// if there are tiles, then it is a dynamic tileset
	if len(tileJSON.Tiles) > 0 {
		imgs := make([]*ebiten.Image, len(tileJSON.Tiles))
		for i, tile := range tileJSON.Tiles {
			// Split the path and take everything after the "../.." parts
			parts := strings.Split(tile.ImagePath, "/")
			// Find where the actual asset path starts (after all ".." entries)
			var assetParts []string
			for _, part := range parts {
				if part != ".." && part != "." && part != "" {
					assetParts = append(assetParts, part)
				}
			}
			// Reconstruct with assets prefix
			imgPath := path.Join("assets", path.Join(assetParts...))
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				return nil, err
			}
			imgs[i] = img
		}

		return &DynamicTileset{
			gid:  gid,
			imgs: imgs,
		}, nil
	}

	// Split the path and take everything after the "../.." parts
	parts := strings.Split(tileJSON.ImagePath, "/")
	// Find where the actual asset path starts (after all ".." entries)
	var assetParts []string
	for _, part := range parts {
		if part != ".." && part != "." && part != "" {
			assetParts = append(assetParts, part)
		}
	}
	// Reconstruct with assets prefix
	imgPath := path.Join("assets", path.Join(assetParts...))

	img, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		return nil, err
	}
	return &UniformTileset{
		gid:     gid,
		img:     img,
		columns: tileJSON.Columns,
	}, nil
}
