package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
    mapSteps [][]*Tile
    currentStep int
}

const (
    waterLevel = 1
    minIslandSize = 50
    maxIslandSize = 100
)

func generateMap() ([]*Tile, IsometricCoordinate, [][]*Tile) {
    tiles, mapSteps := generateIslandFloodFill(rand.Intn(maxIslandSize-minIslandSize)+minIslandSize)
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
    finalTiles := generateMapFromTiles(tiles)
    return finalTiles, camCenter, mapSteps
}

func generateMapFromTiles(tiles []*Tile) ([]*Tile) {
    // make map have water
    tileTypes := make(map[IsometricCoordinate]*Tile)
    for _, tile := range tiles {
        tileTypes[IsometricCoordinate{tile.coord.x, tile.coord.y, 0}] = tile
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
    return finalTiles
}

func pickTileType(position IsometricCoordinate) tileType {
    if position.z < waterLevel+0.2 {
        return "sandTile"
    }
    return "landTile"
}

func _generateIslandFloodFill(size, jumped int) ([]*Tile, bool, [][]*Tile) {
    tiles := make([]*Tile, 0)
    steps := make([][]*Tile, 0)
    tilesNeedNeighbor := make([]*Tile, 0)
    tilesTaken := make(map[IsometricCoordinate]bool)
    firstTile := true
    for {
        if len(tiles) == size {
            return tiles, true, steps
        }
        if len(tilesNeedNeighbor) == 0 && !firstTile {
            fmt.Println("trying new map...")
            return nil, false, nil
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
        firstTile = false
        if !adding {
            if len(tilesNeedNeighbor) >= 1 {
                tilesNeedNeighbor = tilesNeedNeighbor[1:]
            }
            continue
        }
        alt := getNoise(coordAdding.x+float64(jumped), coordAdding.y+float64(jumped))*2+waterLevel
        if alt > waterLevel {
            finalTileCoord := IsometricCoordinate{
                x: coordAdding.x,
                y: coordAdding.y,
                z: alt-waterLevel*0.6,
            } 
            tile := &Tile{
                tileType: pickTileType(finalTileCoord),
                coord: finalTileCoord,
            }
            tiles = append(tiles, tile)
            steps = append(steps, tiles)
            tilesNeedNeighbor = append(tilesNeedNeighbor, tile)
        }
        tilesTaken[coordAdding] = true
    }
}

func generateIslandFloodFill(size int) ([]*Tile, [][]*Tile) {
    fmt.Printf("target size %d\n", size)
    jumped := 0
    for {
        if tiles, valid, mapSteps := _generateIslandFloodFill(size, jumped); valid {
            return tiles, mapSteps
        }
        jumped += 17/13
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

func (g *gameSceneImpl) _updateMapSteps() error {
    g.tilemap.tiles = generateMapFromTiles(g.mapSteps[g.currentStep])
    g.currentStep++
    if g.currentStep < len(g.mapSteps) {
        g.actionQueue.Add(NewTimerAction(g._updateMapSteps, time.Now().Add(100 * time.Millisecond)))
    }
    return nil
}

func (g *gameSceneImpl) Start() error {
    var centerTile IsometricCoordinate
    g.tilemap.tiles, centerTile, g.mapSteps = generateMap()
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
    // g.actionQueue.Add(NewTimerAction(g._updateMapSteps, time.Now().Add(1500 * time.Millisecond)))
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
        3: "sandTile",
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
