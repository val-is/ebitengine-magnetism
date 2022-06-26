package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type titleSceneImpl struct {
    baseScene
}

func (t *titleSceneImpl) Start() error {
    t.actionQueue.Add(func() (bool, error) {
        if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
            if next, err := NewGameScene(t.game); err != nil {
                return false, err
            } else {
                t.game.nextScene = next
            }
        }
        return false, nil
    })
    return nil
}

func (t *titleSceneImpl) Stop() error {
    return nil
}

func (t *titleSceneImpl) Update() error {
    return t.baseScene.Update()
}

func (t *titleSceneImpl) Draw(screen *ebiten.Image) {
}

func NewTitleScene(game *Game) (Scene, error) {
    return &titleSceneImpl{
        baseScene: NewBaseScene(game),
    }, nil
}
