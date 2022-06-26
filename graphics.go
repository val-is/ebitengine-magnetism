package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// <tileset version="1.5" tiledversion="1.6.0" name="isometric-sandbox-sheet" tilewidth="32" tileheight="32" tilecount="2" columns="2" objectalignment="topleft">

func LoadTiledSpritemap(filepath string, imW, imH, tilesW, tilesH int) ([]*ebiten.Image, error) {
    spritemapRaw, _, err := ebitenutil.NewImageFromFile(filepath)
    if err != nil {
        return nil, err
    }
    tileImages := make([]*ebiten.Image, 0)
    for x := 0; x<tilesW; x++ {
        for y := 0; y<tilesH; y++ {
            subImg := spritemapRaw.SubImage(image.Rectangle{
                Min: image.Point{
                    X: x*imW,
                    Y: y*imH,
                },
                Max: image.Point{
                    X: (x+1)*imW,
                    Y: (y+1)*imH,
                },
            })
            ebiImg := ebiten.NewImageFromImage(subImg)
            tileImages = append(tileImages, ebiImg)
        }
    }
    return tileImages, nil
}
