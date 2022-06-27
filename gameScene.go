package main

import (
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type FacingDirection int

const (
    walkSpeed = 5.0/60.0

    FACING_LEFT = 0
    FACING_RIGHT = 1

    playerWidth = tileWidth * 0.5
    playerHeight = tileHeight * 0.5

    playerCameraMaxDist = 2
    playerCameraMoveSpeed = walkSpeed
)

type PlayerCharacter struct {
    pos IsometricCoordinate
    sprites []*ebiten.Image
    facing FacingDirection
}

func (p *PlayerCharacter) Draw(screen *ebiten.Image, cameraPos IsometricCoordinate) {
    img := p.sprites[p.facing]
    w, h := img.Size()
    screenCoord := iso2Screen(IsometricCoordinate{
        x: p.pos.x - cameraPos.x,
        y: p.pos.y - cameraPos.y,
        z: p.pos.z - cameraPos.z,
    })
    drawOpt := ebiten.DrawImageOptions{}
    drawOpt.GeoM.Reset()
    drawOpt.GeoM.Scale(playerWidth/float64(w), playerHeight/float64(h))
    drawOpt.GeoM.Translate(
        screenCoord.x + 1920/2 - playerWidth/2,
        screenCoord.y + 1080/2 - playerHeight/2,
    )
    screen.DrawImage(img, &drawOpt)
}

func (p *PlayerCharacter) ScreenPosition(cameraPos IsometricCoordinate) (ScreenCoordinate) {
    screenCoord := iso2Screen(IsometricCoordinate{
        x: p.pos.x - cameraPos.x,
        y: p.pos.y - cameraPos.y,
        z: p.pos.z - cameraPos.z,
    })
    return ScreenCoordinate{
        screenCoord.x + 1920/2 - playerWidth/2,
        screenCoord.y + 1080/2 - playerHeight/2,
    }
}

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
    mapSteps [][]*Tile
    currentStep int
    mapRotation float64
    
    player *PlayerCharacter
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
    g.player.pos = screen2Iso(centerTileScreen)

    playerSpritesheet, err := LoadTiledSpritemap("./resources/Tiny_Tales_Wild_Beasts_NPC_1.0/RPG_Maker/32/$Fox_1.png", 32, 32, 3, 4, 0, 0)
    if err != nil {
        return err
    }
    g.player.sprites = playerSpritesheet
    
    g.actionQueue.Add(func() (bool, error) {
        if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
            mouseX, mouseY := ebiten.CursorPosition()
            playerScreenPos := g.player.ScreenPosition(g.cameraPos)
            mouseXCentered := float64(mouseX)-playerScreenPos.x
            mouseYCentered := float64(mouseY)-playerScreenPos.y
            clickVec := screen2Iso(ScreenCoordinate{
                x: mouseXCentered,
                y: mouseYCentered,
            })
            clickVecDist := math.Hypot(clickVec.x, clickVec.y)
            moveVec := IsometricCoordinate{
                x: clickVec.x/clickVecDist * walkSpeed,
                y: clickVec.y/clickVecDist * walkSpeed,
            }
            newPlayerPos := IsometricCoordinate{
                x: g.player.pos.x+moveVec.x,
                y: g.player.pos.y+moveVec.y,
            }
            if mouseXCentered < 0 {
                g.player.facing = FACING_LEFT
            } else {
                g.player.facing = FACING_RIGHT
            }
            tilesAtPos := g.tilemap.GetTilesAt(newPlayerPos)
            for _, tile := range tilesAtPos {
                if tile.walkable {
                    g.player.pos = newPlayerPos
                    g.player.pos.z = tile.coord.z + 0.5
                    break
                }
            }
        }
        playerCamDistVec := IsometricCoordinate{
            x: g.player.pos.x - g.cameraPos.x,
            y: g.player.pos.y - g.cameraPos.y,
        }
        playerCamDist := math.Hypot(playerCamDistVec.x, playerCamDistVec.y)
        if playerCamDist > playerCameraMaxDist {
            g.cameraPos = IsometricCoordinate{
                x: g.cameraPos.x + playerCamDistVec.x/playerCamDist*playerCameraMoveSpeed,
                y: g.cameraPos.y + playerCamDistVec.y/playerCamDist*playerCameraMoveSpeed,
                z: g.cameraPos.z,
            }
        }
        // g.cameraPos = g.player.pos
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
    g.player.Draw(screen, g.cameraPos)
}

func NewGameScene(game *Game) (Scene, error) {
    tilemap, err := NewTilemap("./resources/isometric-sandbox-32x32/isometric-sandbox-sheet.png", idxTileMap)
    if err != nil {
        return nil, err
    }
    return &gameSceneImpl{
        baseScene: NewBaseScene(game),                                                              
        tilemap: tilemap,
        cameraPos: IsometricCoordinate{},
        player: &PlayerCharacter{
            pos: IsometricCoordinate{},
        },
    }, nil
}
