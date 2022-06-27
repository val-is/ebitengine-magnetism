package main

import (
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
)

type ScreenCoordinate struct {
    x, y float64
}

type IsometricCoordinate struct {
    x, y, z float64
}

const (
    i_x = 1
    i_y = 0.5
    j_x = -1
    j_y = 0.5
)

func invMatrix(a, b, c, d float64) (i, j, k, l float64) {
    det := (1 / (a*d - b*c))
    return det*d, det*-b, det*-c, det*a
}

// TODO eff++
func screen2Iso(s ScreenCoordinate) IsometricCoordinate {
    a := i_x * 0.5 * tileWidth
    b := j_x * 0.5 * tileWidth
    c := i_y * 0.5 * tileHeight
    d := j_y * 0.5 * tileHeight

    inv_a, inv_b, inv_c, inv_d := invMatrix(a, b, c, d)
    
    return IsometricCoordinate{
        x: s.x * inv_a + s.y * inv_b,
        y: s.x * inv_c + s.y * inv_d,
    }
}

func iso2Screen(i IsometricCoordinate) ScreenCoordinate {
    return ScreenCoordinate{
        x: i.x * i_x * 0.5 * tileWidth + i.y * j_x * 0.5 * tileWidth,
        y: i.x * i_y * 0.5 * tileHeight + i.y * j_y * 0.5 * tileHeight,
    }
}

func getAdjIsometric(i IsometricCoordinate) []IsometricCoordinate {
    // assumes int passed in, i.e. applies offset of 1
    return []IsometricCoordinate{
        {i.x+1, i.y, i.z},
        {i.x-1, i.y, i.z},
        {i.x, i.y+1, i.z},
        {i.x, i.y-1, i.z},
    }
}

var perlinGen *perlin.Perlin

func init() {
    perlinGen = perlin.NewPerlin(2, 2, 3, time.Now().Unix())
    rand.Seed(time.Now().Unix())
}

func getNoise(x, y float64) float64 {
    return perlinGen.Noise2D(x/13, y/13)
}
