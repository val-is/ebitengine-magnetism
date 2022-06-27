package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
    ebiten.SetWindowSize(960, 540)
    ebiten.SetWindowTitle("ebitengine magnet fishing")
    ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

    g := &Game{}
    g.nextScene, _ = NewTitleScene(g)

    if err := ebiten.RunGame(g); err != nil {
        panic(err)
    }
}
