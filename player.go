package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/yohamta/ganim8/v2"
	"image"
)

type player struct {
	input     *input.Handler
	pos       image.Point
	run       *ganim8.Animation
	idle      *ganim8.Animation
	attack    *ganim8.Animation
	direction input.Action
}

func newPlayer(input *input.Handler, pos image.Point, run *ganim8.Animation, idle *ganim8.Animation, attack *ganim8.Animation, direction input.Action) *player {
	return &player{input: input, pos: pos, run: run, idle: idle, attack: attack, direction: direction}
}

func (p *player) Update() {

	if p.input.ActionIsPressed(ActionAttack) {
		p.attack.Update()
	}

	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.run.Update()
		p.pos.X -= 4
		if p.direction == ActionMoveRight {
			p.direction = ActionMoveLeft
			p.run.Sprite().FlipH()
			p.idle.Sprite().FlipH()
		}
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.run.Update()
		p.pos.X += 4
		if p.direction == ActionMoveLeft {
			p.direction = ActionMoveRight
			p.run.Sprite().FlipH()
			p.idle.Sprite().FlipH()
		}

	}
	if p.input.ActionIsPressed(ActionMoveUp) {
		p.run.Update()
		p.pos.Y -= 4
	}
	if p.input.ActionIsPressed(ActionMoveDown) {
		p.run.Update()
		p.pos.Y += 4
	}
	p.idle.Update()
}

func (p *player) Draw(screen *ebiten.Image) {
	if p.input.ActionIsPressed(ActionMoveLeft) || p.input.ActionIsPressed(ActionMoveRight) || p.input.ActionIsPressed(ActionMoveUp) || p.input.ActionIsPressed(ActionMoveDown) {
		p.run.Draw(screen, ganim8.DrawOpts(float64(p.pos.X), float64(p.pos.Y), 0, 1, 1, 0.5, 0.5))
	} else if p.input.ActionIsPressed(ActionAttack) {
		p.attack.Draw(screen, ganim8.DrawOpts(float64(p.pos.X), float64(p.pos.Y), 0, 1, 1, 0.5, 0.5))
	} else {
		p.idle.Draw(screen, ganim8.DrawOpts(float64(p.pos.X), float64(p.pos.Y), 0, 1, 1, 0.5, 0.5))
	}

}
