package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileWidth  = 1920 / 7
	tileHeight = 1920 / 7

	tileSidePx = 9.2376 / 32 * tileHeight
)

type tileType string

type Tile struct {
	tileType tileType
	coord    IsometricCoordinate
	walkable bool
}

func (t *Tile) CollidesWith(c IsometricCoordinate) bool {
	return t.coord.x-0.5 <= c.x+0.5 && c.x+0.5 <= t.coord.x+0.5 &&
		t.coord.y-0.5 <= c.y+0.5 && c.y+0.5 <= t.coord.y+0.5
}

type Tilemap struct {
	spritemap map[tileType]*ebiten.Image // TODO animated/custom sprite class support
	tiles     []*Tile
}

func NewTilemap(filepath string, indexToType map[int]tileType) (*Tilemap, error) {
	loadedTiles, err := LoadTiledSpritemap(filepath, 32, 32, 6, 9, 0, 0) // TODO hardcoded lol
	if err != nil {
		return nil, err
	}
	t := &Tilemap{
		spritemap: make(map[tileType]*ebiten.Image),
		tiles:     make([]*Tile, 0),
	}
	for idx, tType := range indexToType {
		if idx > len(loadedTiles) {
			return nil, fmt.Errorf("%d out of range of loaded tiles (%d)", idx, len(loadedTiles))
		}
		t.spritemap[tType] = loadedTiles[idx]
	}
	return t, nil
}

func (t *Tilemap) Draw(screen *ebiten.Image, cameraPos IsometricCoordinate) {
	drawOpt := ebiten.DrawImageOptions{}
	for _, tile := range t.tiles {
		img, present := t.spritemap[tile.tileType]
		if !present {
			panic("sprite " + tile.tileType + " not set up!!!")
		}
		screenCoord := iso2Screen(IsometricCoordinate{
			x: tile.coord.x - cameraPos.x,
			y: tile.coord.y - cameraPos.y,
            z: tile.coord.z - cameraPos.z,
		})
		w, h := img.Size()
		drawOpt.GeoM.Reset()
		drawOpt.GeoM.Scale(tileWidth/float64(w), tileWidth/float64(h))
		drawOpt.GeoM.Translate(
			screenCoord.x+1920/2-tileWidth/2,
			screenCoord.y+1080/2-tileHeight/2,
		)
		screen.DrawImage(img, &drawOpt)
	}
}

func (t *Tilemap) GetTilesAt(coordinate IsometricCoordinate) []*Tile {
	tiles := make([]*Tile, 0)
	for _, tile := range t.tiles {
		if tile.CollidesWith(coordinate) {
			tiles = append(tiles, tile)
		}
	}
	return tiles
}

func (t *Tilemap) GetClickedCoordinate(screenCoord ScreenCoordinate) IsometricCoordinate {
	// ignores z val
	return IsometricCoordinate{}
}
