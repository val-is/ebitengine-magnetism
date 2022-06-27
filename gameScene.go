package main

import (
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
    mapSteps [][]*Tile
    currentStep int
    mapRotation float64
}

func (g *gameSceneImpl) _updateMapSteps() error {
    g.tilemap.tiles = generateMapFromTiles(g.mapSteps[g.currentStep])
    g.currentStep++
    if g.currentStep < len(g.mapSteps) {
        g.actionQueue.Add(NewTimerAction(g._updateMapSteps, time.Now().Add(20 * time.Millisecond)))
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
    tilemap, err := NewTilemap("./resources/isometric-sandbox-32x32/isometric-sandbox-sheet.png", idxTileMap)
    if err != nil {
        return nil, err
    }
    return &gameSceneImpl{
        baseScene: NewBaseScene(game),                                                              
        tilemap: tilemap,
        cameraPos: IsometricCoordinate{x: 0, y: 0},
    }, nil
}
