package animations

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
	"github.com/yohamta/ganim8/v2"
	"image/color"
	"log"
	"time"
)

type Animations struct {
	Run    *ganim8.Animation
	Idle   *ganim8.Animation
	Attack *ganim8.Animation
	Jump   *ganim8.Animation
	Fall   *ganim8.Animation
}

func InitPlayerAnimations() Animations {
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
	jump, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(jumpSprite))
	if err != nil {
		log.Fatal(err)
	}
	fall, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(Fall_png))
	if err != nil {
		log.Fatal(err)
	}
	runGrid := ganim8.NewGrid(200, 200, 1600, 200)
	idleGrid := ganim8.NewGrid(200, 200, 1600, 200)
	attackGrid := ganim8.NewGrid(200, 200, 1200, 200)
	jumpGrid := ganim8.NewGrid(200, 200, 400, 200)
	fallGrid := ganim8.NewGrid(200, 200, 400, 200)

	animations := Animations{
		ganim8.New(run, runGrid.Frames("1-8", 1), 100*time.Millisecond),
		ganim8.New(idle, idleGrid.Frames("1-8", 1), 100*time.Millisecond),
		ganim8.New(attack, attackGrid.Frames("1-6", 1), 60*time.Millisecond),
		ganim8.New(jump, jumpGrid.Frames("1-2", 1), 100*time.Millisecond),
		ganim8.New(fall, fallGrid.Frames("1-2", 1), 100*time.Millisecond)}

	return animations
}

func DrawPolygon(screen *ebiten.Image, polygon *resolv.ConvexPolygon, color color.Color) {

	for _, line := range polygon.Lines() {
		ebitenutil.DrawLine(screen, line.Start.X(), line.Start.Y(), line.End.X(), line.End.Y(), color)
	}

}

func (a *Animations) FlipAll() {
	a.Run.Sprite().FlipH()
	a.Idle.Sprite().FlipH()
	a.Attack.Sprite().FlipH()
	a.Jump.Sprite().FlipH()
	a.Fall.Sprite().FlipH()
}
