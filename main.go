package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	camera "github.com/melonfunction/ebiten-camera"
	"github.com/pablogolobaro/samuraiDuty/animations"
	"github.com/pablogolobaro/samuraiDuty/entity"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/dngn"
	"github.com/solarlune/paths"
	"github.com/solarlune/resolv"
	"image/color"
	"log"
)

const (
	Width  = 440 * 2
	Height = 320 * 2
)

type Game struct {
	Width       int
	Height      int
	Player      *entity.Player
	inputSystem input.System
	back        *ebiten.Image
	cam         *camera.Camera
	Space       *resolv.Space
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	return outsideWidth, outsideHeight
}

func NewGame() *Game {
	back, _, err := ebitenutil.NewImageFromFile("./resources/Background.png")
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{back: back, Width: Width, Height: Height}

	g.cam = camera.NewCamera(g.Width, g.Height, 0, 0, 0, 1)

	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	keymap := input.Keymap{
		entity.ActionMoveLeft:  {input.KeyLeft},
		entity.ActionMoveRight: {input.KeyRight},
		entity.ActionMoveUp:    {input.KeyUp},
		entity.ActionMoveDown:  {input.KeyDown},
		entity.ActionAttack:    {input.KeyF},
	}

	playerAnimations := animations.InitPlayerAnimations()

	g.Player = entity.NewPlayer(
		g.inputSystem.NewHandler(0, keymap),
		playerAnimations,
	)

	g.InitSpace()
	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cam.Surface.Clear()

	//backOpts := &ebiten.DrawImageOptions{}
	//g.cam.Surface.DrawImage(g.back, g.cam.GetTranslation(backOpts, -float64(g.Width/2), -float64(g.Height/2)))

	screen.DrawImage(g.back, &ebiten.DrawImageOptions{})

	for _, o := range g.Space.Objects() {

		if o.HasTags("ramp") {
			drawColor := color.RGBA{255, 50, 100, 255}
			tri := o.Shape.(*resolv.ConvexPolygon)
			animations.DrawPolygon(screen, tri, drawColor)
		} else if o.HasTags("entity") {

		} else {
			drawColor := color.RGBA{60, 60, 60, 255}
			ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
		}

	}

	g.Player.Draw(screen)
	//g.cam.Blit(screen)

	xw, yw := g.cam.GetWorldCoords(g.cam.X, g.cam.Y)
	xm, ym := g.cam.GetCursorCoords()
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Camera:\n  X: %3.3f\n  Y: %3.3f\n  W: %d\n  H: %d\n  Rot: %3.3f\n  Zoom: %3.3f\n"+
				"Tiles:\n  PlayerX: %f\n  PlayerY: %f\n  X %f\n Y %f\n mouseX %f\n mouseY %f\n action %d\n",
			g.cam.X, g.cam.Y, g.cam.Surface.Bounds().Size().X, g.cam.Surface.Bounds().Size().Y, g.cam.Rot, g.cam.Scale,
			g.Player.Physics.Object.X, g.Player.Physics.Object.Y, xw, yw, xm, ym, g.Player.Action,
		))
}

func (g *Game) Update() error {
	g.inputSystem.Update()
	g.cam.SetPosition(g.Player.Physics.Object.X, g.Player.Physics.Object.Y)

	g.Player.Update()

	return nil
}

func main() {
	_ = dngn.Layout{}
	_ = paths.Cell{}

	game := NewGame()
	ebiten.SetWindowSize(game.Width, game.Height)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
