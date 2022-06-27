package main

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
    currentScene Scene
    nextScene Scene
}

func (g *Game) Update() error {
    if g.nextScene != nil {
        if g.currentScene != nil {
            if err := g.currentScene.Stop(); err != nil {
                return err
            }
        }
        g.currentScene = g.nextScene
        g.nextScene = nil
        if err := g.currentScene.Start(); err != nil {
            return err
        }
    }
    return g.currentScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.currentScene.Draw(screen)
}

func (g *Game) Layout(oW, oH int) (sW, sH int) {
	return 1920, 1080
}

type Scene interface {
    Start() error
    Stop() error
    Update() error
    Draw(screen *ebiten.Image)
}

type baseScene struct {
    actionQueue ActionQueue
    game *Game
}
    
func (b *baseScene) Update() error {
    return b.actionQueue.Update()
} 

func NewBaseScene(game *Game) baseScene {
    return baseScene{
        actionQueue: ActionQueue{
            actions: make(map[int64]Action),
        },
        game: game,
    }
}
