package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
}

func (g *gameSceneImpl) Start() error {
    for x := -100; x < 100; x++ {
        for y := -100; y < 100; y++ {
            g.tilemap.tiles = append(g.tilemap.tiles, &Tile{
                tileType: "testTile",
                coord: IsometricCoordinate{
                    x: float64(x),
                    y: float64(y),
                    z: getNoise(float64(x), float64(y)),
                },
            })
        }
    }
    g.actionQueue.Add(func() (bool, error) {
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            mouseX, mouseY := ebiten.CursorPosition()
            worldPos := screen2Iso(ScreenCoordinate{
                x: (float64(mouseX)-1920/2)/100,
                y: (float64(mouseY)-1080/2)/100,
            })
            g.cameraPos = IsometricCoordinate{
                x: g.cameraPos.x-worldPos.x,
                y: g.cameraPos.y-worldPos.y,
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
        0: "testTile",
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
