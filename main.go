package main

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	camera "github.com/melonfunction/ebiten-camera"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/dngn"
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
	"github.com/yohamta/ganim8/v2"
	img "image"
	"log"
	"time"
)

var screenWidth = 440
var screenHeight = 360

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionAttack
)

type Game struct {
	p           *player
	inputSystem input.System
	back        *ebiten.Image
	cam         *camera.Camera
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewGame() *Game {
	back, _, err := ebitenutil.NewImageFromFile("./resources/Background.png")
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{back: back, cam: camera.NewCamera(screenWidth*2, screenHeight*2, 0, 0, 0, 1)}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyLeft},
		ActionMoveRight: {input.KeyRight},
		ActionMoveUp:    {input.KeyUp},
		ActionMoveDown:  {input.KeyDown},
		ActionAttack:    {input.KeyF},
	}

	run, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	idle, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(Idle_png))
	if err != nil {
		log.Fatal(err)
	}
	attack, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(Attack2_png))
	if err != nil {
		log.Fatal(err)
	}
	runGrid := ganim8.NewGrid(200, 200, 1600, 200)
	idleGrid := ganim8.NewGrid(200, 200, 1600, 200)
	attackGrid := ganim8.NewGrid(200, 200, 1200, 200)

	g.p = newPlayer(
		g.inputSystem.NewHandler(0, keymap),
		img.Point{X: 10, Y: 350},
		ganim8.New(run, runGrid.Frames("1-8", 1), 100*time.Millisecond),
		ganim8.New(idle, idleGrid.Frames("1-8", 1), 100*time.Millisecond),
		ganim8.New(attack, attackGrid.Frames("1-6", 1), 100*time.Millisecond),
		ActionMoveRight,
	)

	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.back, &ebiten.DrawImageOptions{})

	g.cam.Surface.Clear()

	backOpts := &ebiten.DrawImageOptions{}

	g.cam.Surface.DrawImage(screen, g.cam.GetTranslation(backOpts, 0, 0))

	g.p.Draw(g.cam.Surface)
	g.cam.Blit(screen)
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Camera:\n  X: %3.3f\n  Y: %3.3f\n  W: %d\n  H: %d\n  Rot: %3.3f\n  Zoom: %3.3f\n"+
				"Tiles:\n  PlayerX: %d\n  PlayerY: %d\n ",
			g.cam.X, g.cam.Y, g.cam.Surface.Bounds().Size().X, g.cam.Surface.Bounds().Size().Y, g.cam.Rot, g.cam.Scale,
			g.p.pos.X, g.p.pos.Y,
		))
}

func (g *Game) Update() error {
	g.inputSystem.Update()
	g.cam.SetPosition(float64(g.p.pos.X)+float64(100)/2, float64(g.p.pos.Y)+float64(100)/2)

	g.p.Update()
	// Note: it assumes that the time delta is 16ms by default
	//       if you need to specify different delta you can use Animation.UpdateWithDelta(delta) instead
	return nil
}

func main() {
	_ = resolv.NewCollision()
	_ = dngn.Layout{}
	_ = paths.Cell{}

	game := NewGame()
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
