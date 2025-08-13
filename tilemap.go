package main

import (
	"encoding/json"
	"os"
	"path"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

type TilemapJSON struct {
	Layers   []TilemapLayerJSON `json:"layers"`
	Tilesets []TilesetJSON      `json:"tilesets"`
	Width    int                `json:"width"`
	Height   int                `json:"height"`
}

func (t *TilemapJSON) GetTilesets() ([]Tileset, error) {
	tileSets := make([]Tileset, len(t.Tilesets))
	for i, tileset := range t.Tilesets {
		tilesetPath := path.Join("assets/maps/", tileset.Source)
		tileSet, err := NewTilesetFromPath(tilesetPath, tileset.FirstGid)
		if err != nil {
			return nil, err
		}
		tileSets[i] = tileSet
	}

	return tileSets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	jsonFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJson TilemapJSON
	err = json.Unmarshal(jsonFile, &tilemapJson)
	if err != nil {
		return nil, err
	}

	return &tilemapJson, nil
}
