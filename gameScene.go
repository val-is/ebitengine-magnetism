package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
}

const (
    waterLevel = 1
    minIslandSize = 150
    maxIslandSize = 200
)

func generateMap() ([]*Tile, IsometricCoordinate) {
    tiles := generateIslandFloodFill(rand.Intn(maxIslandSize-minIslandSize)+minIslandSize)
    // make map have water
    tileTypes := make(map[IsometricCoordinate]*Tile)
    camCenter := IsometricCoordinate{}
    highestAlt := -100.0
    for _, tile := range tiles {
        tileTypes[IsometricCoordinate{tile.coord.x, tile.coord.y, 0}] = tile
        if tile.coord.z > float64(highestAlt) {
            highestAlt = tile.coord.z
            camCenter = IsometricCoordinate{tile.coord.x, tile.coord.y, 0}
        }
    }
    finalTiles := make([]*Tile, 0)
    for x := -50; x<50; x++ {
        for y := -50; y<50; y++ {
            if tile, p := tileTypes[IsometricCoordinate{float64(x), float64(y), 0}]; p {
                finalTiles = append(finalTiles, tile)
            } else {
                finalTiles = append(finalTiles, &Tile{
                    tileType: "waterTile",
                    coord: IsometricCoordinate{float64(x), float64(y), waterLevel},
                })
            }
        }
    }
    return finalTiles, camCenter
}

func _generateIslandFloodFill(size, jumped int) ([]*Tile, bool) {
    tiles := make([]*Tile, 0)
    tilesNeedNeighbor := make([]*Tile, 0)
    tilesTaken := make(map[IsometricCoordinate]bool)
    for {
        if len(tiles) == size {
            return tiles, true
        }
        if len(tiles) != 0 && len(tilesNeedNeighbor) == 0 {
            return nil, false
        }
        var seedPos IsometricCoordinate
        if len(tiles) == 0 {
            seedPos = IsometricCoordinate{0, 0, 0}
        } else {
            sort.Slice(tilesNeedNeighbor, func(i, j int) bool {
                return tilesNeedNeighbor[i].coord.z > tilesNeedNeighbor[j].coord.z
            })
            seedPos = tilesNeedNeighbor[0].coord
        }
        seedPos = IsometricCoordinate{seedPos.x, seedPos.y, 0}
        adj := getAdjIsometric(seedPos)
        var coordAdding IsometricCoordinate
        adding := false
        for _, coord := range adj {
            if _, p := tilesTaken[coord]; !p {
                coordAdding = coord
                adding = true
                break
            }
        }
        if !adding {
            if len(tilesNeedNeighbor) >= 1 {
                tilesNeedNeighbor = tilesNeedNeighbor[1:]
            }
            continue
        }
        alt := getNoise(coordAdding.x+ float64(jumped), coordAdding.y+float64(jumped)) * 3
        if alt > waterLevel {
            tile := &Tile{
                tileType: "landTile",
                coord: IsometricCoordinate{
                    x: coordAdding.x,
                    y: coordAdding.y,
                    z: alt,
                },
            }
            tiles = append(tiles, tile)
            tilesNeedNeighbor = append(tilesNeedNeighbor, tile)
        }
        if len(tiles) == 0 {
            return nil, false
        }
        tilesTaken[coordAdding] = true
    }
}

func generateIslandFloodFill(size int) []*Tile {
    fmt.Printf("target size %d\n", size)
    jumped := 0
    for {
        if tiles, valid := _generateIslandFloodFill(size, jumped); valid {
            return tiles
        }
        jumped += 100
    }
}

func generateMapRaw() []*Tile {
    tiles := make([]*Tile, 0)
    for x := -100; x < 100; x++ {
        for y := -100; y < 100; y++ {
            alt := getNoise(float64(x), float64(y)) * 3
            if alt > waterLevel {
                tiles = append(tiles, &Tile{
                    tileType: "landTile",
                    coord: IsometricCoordinate{
                        x: float64(x),
                        y: float64(y),
                        z: alt,
                    },
                })
            } else {
                tiles = append(tiles, &Tile{
                    tileType: "waterTile",
                    coord: IsometricCoordinate{
                        x: float64(x),
                        y: float64(y),
                        z: waterLevel,
                    },
                })
            }
        }
    }
    return tiles
}

func (g *gameSceneImpl) Start() error {
    var centerTile IsometricCoordinate
    g.tilemap.tiles, centerTile = generateMap()
    centerTileScreen := iso2Screen(centerTile)
    g.cameraPos = screen2Iso(centerTileScreen)
    g.actionQueue.Add(func() (bool, error) {
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            mouseX, mouseY := ebiten.CursorPosition()
            worldPos := screen2Iso(ScreenCoordinate{
                x: (float64(mouseX)-1920/2)/100,
                y: (float64(mouseY)-1080/2)/100,
            })
            g.cameraPos = IsometricCoordinate{
                x: g.cameraPos.x+worldPos.x,
                y: g.cameraPos.y+worldPos.y,
            }
        }
        return false, nil
    })
    return nil
}

func (g *gameSceneImpl) Stop() error {
    return nil
}

func (g *gameSceneImpl) Update() error {
    return g.baseScene.Update()
}

func (g *gameSceneImpl) Draw(screen *ebiten.Image) {
    g.tilemap.Draw(screen, g.cameraPos)
}

func NewGameScene(game *Game) (Scene, error) {
    tilemap, err := NewTilemap("./resources/isometric-sandbox-32x32/isometric-sandbox-sheet.png", map[int]tileType{
        0: "landTile",
        2: "waterTile",
    })
    if err != nil {
        return nil, err
    }
    return &gameSceneImpl{
        baseScene: NewBaseScene(game),                                                              
        tilemap: tilemap,
        cameraPos: IsometricCoordinate{x: 0, y: 0},
    }, nil
}
