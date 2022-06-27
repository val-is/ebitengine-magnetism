package main

import (
	"math/rand"
	"sort"
)

const (
    waterLevel = 1
    minIslandSize = 50
    maxIslandSize = 100
)

var (
    idxTileMap = map[int]tileType{
        0: "landTile",
        2: "waterTile",
        3: "sandTile",
    }
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
    if position.z < waterLevel*1.1 {
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
        if len(tiles) > maxIslandSize {
            return nil, false, nil
            // return tiles, true, steps
        }
        if len(tilesNeedNeighbor) == 0 && !firstTile {
            break
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
        alt := getNoise(coordAdding.x+float64(jumped), coordAdding.y+float64(jumped))*3+waterLevel
        if alt > waterLevel {
            finalTileCoord := IsometricCoordinate{
                x: coordAdding.x,
                y: coordAdding.y,
                z: alt-waterLevel*0.2,
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
    if minIslandSize < len(tiles) && len(tiles) < maxIslandSize {
        return tiles, true, steps
    }
    return nil, false, nil
}

func generateIslandFloodFill(size int) ([]*Tile, [][]*Tile) {
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
