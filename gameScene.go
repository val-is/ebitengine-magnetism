package main

import (
	_ "image/png"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type WorldObjectDrawable interface {
    Draw(screen *ebiten.Image, cameraPos IsometricCoordinate)
    ScreenPosition(cameraPos IsometricCoordinate) ScreenCoordinate
    GetBottom(cameraPos IsometricCoordinate) ScreenCoordinate
}

func DrawWorldObjects(screen *ebiten.Image, cameraPos IsometricCoordinate, objects []WorldObjectDrawable) {
    // TODO sorting every frame sus
    sort.Slice(objects, func(i, j int) bool {
        aPos := objects[i].GetBottom(cameraPos)
        bPos := objects[j].GetBottom(cameraPos)
        return aPos.y < bPos.y
    })

    for _, object := range objects {
        object.Draw(screen, cameraPos)
    }
}

type WorldObject struct {
    pos IsometricCoordinate
    width, height float64
}

func (w *WorldObject) ScreenPosition(cameraPos IsometricCoordinate) ScreenCoordinate {
    screenCoord := iso2Screen(IsometricCoordinate{
        x: w.pos.x - cameraPos.x,
        y: w.pos.y - cameraPos.y,
        z: w.pos.z - cameraPos.z,
    })
    return ScreenCoordinate{
        screenCoord.x + 1920/2 - w.width/2,
        screenCoord.y + 1080/2 - w.height/2,
    }
}

func (w *WorldObject) GetBottom(cameraPos IsometricCoordinate) ScreenCoordinate {
    scrPos := w.ScreenPosition(cameraPos)
    return ScreenCoordinate{
        x: scrPos.x,
        y: scrPos.y + w.height,
    }
}

func (wo *WorldObject) DrawWithImg(screen *ebiten.Image, cameraPos IsometricCoordinate, img *ebiten.Image) {
    w, h := img.Size()
    screenCoord := wo.ScreenPosition(cameraPos)
    drawOpt := ebiten.DrawImageOptions{}
    drawOpt.GeoM.Reset()
    drawOpt.GeoM.Scale(wo.width/float64(w), wo.height/float64(h))
    drawOpt.GeoM.Translate(screenCoord.x, screenCoord.y)
    screen.DrawImage(img, &drawOpt)
}

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
    WorldObject
    sprites []*ebiten.Image
    facing FacingDirection
}

func (p *PlayerCharacter) Draw(screen *ebiten.Image, cameraPos IsometricCoordinate) {
    img := p.sprites[p.facing]
    p.DrawWithImg(screen, cameraPos, img)
}

type FoliageType int

const (
    foliageProb = 0.5

    FOLIAGE_GRASS = 0
    FOLIAGE_TREE = 1
)

var (
    foliageSprites map[FoliageType]*ebiten.Image
)

type Foliage struct {
    WorldObject
    foliageType  FoliageType
}

func (f *Foliage) Draw(screen *ebiten.Image, cameraPos IsometricCoordinate) {
    img := foliageSprites[f.foliageType]
    f.DrawWithImg(screen, cameraPos, img)
}

type gameSceneImpl struct {
    baseScene
    tilemap *Tilemap
    cameraPos IsometricCoordinate
    mapSteps [][]*Tile
    currentStep int
    mapRotation float64
    
    drawing []WorldObjectDrawable

    player *PlayerCharacter
    foliage []*Foliage
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
    g.drawing = make([]WorldObjectDrawable, 0)

    playerSpritesheet, err := LoadTiledSpritemap("./resources/Tiny_Tales_Wild_Beasts_NPC_1.0/RPG_Maker/32/$Fox_1.png", 32, 32, 3, 4, 0, 0)
    if err != nil {
        return err
    }
    g.player.sprites = playerSpritesheet
    g.drawing = append(g.drawing, g.player)

    foliageSpritesheetRaw, err := LoadTiledSpritemap("./resources/48x48 & 16x32 Trees/16x32 trees.png", 16, 32, 4, 2, 0, 0)
    if err != nil {
        return err
    }
    foliageSprites = map[FoliageType]*ebiten.Image{
        FOLIAGE_GRASS: foliageSpritesheetRaw[7],
        FOLIAGE_TREE: foliageSpritesheetRaw[1],
    }
    g.foliage = make([]*Foliage, 0)
    for _, tile := range g.tilemap.tiles {
        if tile.tileType == TILE_LAND && rand.Float64() < foliageProb {
            newFoliage := &Foliage{
                WorldObject: WorldObject{
                    pos: IsometricCoordinate{
                        tile.coord.x,
                        tile.coord.y,
                        tile.coord.z+1.5,
                    },
                    width: 0.5 * tileWidth,
                    height: tileHeight,
                },
                foliageType: FoliageType(rand.Intn(2)),
            }
            g.foliage = append(g.foliage, newFoliage)
            g.drawing = append(g.drawing, newFoliage)
        }
    }
    
    g.actionQueue.Add(func() (bool, error) {
        // move player
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

        // lock player to camera view
        playerScreenPos := g.player.ScreenPosition(g.cameraPos)
        playerCamDistVec := screen2Iso(ScreenCoordinate{
            x: playerScreenPos.x - 1920/2,
            y: playerScreenPos.y - 1080/2,
        })
        playerCamDist := math.Hypot(playerCamDistVec.x, playerCamDistVec.y)
        if playerCamDist > playerCameraMaxDist {
            g.cameraPos = IsometricCoordinate{
                x: g.cameraPos.x + playerCamDistVec.x/playerCamDist*playerCameraMoveSpeed,
                y: g.cameraPos.y + playerCamDistVec.y/playerCamDist*playerCameraMoveSpeed,
                z: g.cameraPos.z,
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
    DrawWorldObjects(screen, g.cameraPos, g.drawing)
    // g.player.Draw(screen, g.cameraPos)
    // for _, foliage := range g.foliage {
    //     foliage.Draw(screen, g.cameraPos)
    // }
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
            WorldObject: WorldObject{
                pos: IsometricCoordinate{},
                width: playerWidth,
                height: playerHeight,
            },
        },
    }, nil
}
