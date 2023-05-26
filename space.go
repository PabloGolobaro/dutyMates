package main

import (
	"github.com/solarlune/resolv"
)

func (g *Game) InitSpace() {
	gw := float64(g.Width)
	gh := float64(g.Height)

	// Define the world's Space. Here, a Space is essentially a grid (the game's width and height, or 640x360), made up of 16x16 cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space to see if the Cells at specific positions contain (or would contain) Objects. This is a broad, simplified approach to collision
	// detection.
	g.Space = resolv.NewSpace(int(gw), int(gh), 16, 16)

	// Construct the solid level geometry. Note that the simple approach of checking cells in a Space for collision works simply when the geometry is aligned with the cells,
	// as it all is in this platformer example.

	g.Space.Add(
		resolv.NewObject(0, 0, 16, gh, "solid"),
		resolv.NewObject(gw-16, 0, 16, gh, "solid"),
		resolv.NewObject(0, 0, gw, 16, "solid"),
		resolv.NewObject(0, gh-24, gw, 32, "solid"),
		resolv.NewObject(160, gh-56, 160, 32, "solid"),
		resolv.NewObject(320, 64, 32, 160, "solid"),
		resolv.NewObject(64, 128, 16, 160, "solid"),
		resolv.NewObject(gw-128, 64, 128, 16, "solid"),
		resolv.NewObject(gw-128, gh-88, 128, 16, "solid"),
	)

	// Create the Player. NewPlayer adds it to the world's Space.
	g.Space.Add(g.Player.Physics.Object)

	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.

	// Non-moving floating Platforms.
	g.Space.Add(
		resolv.NewObject(352, 64, 48, 8, "platform"),
		resolv.NewObject(352, 64+64, 48, 8, "platform"),
		resolv.NewObject(352, 64+128, 48, 8, "platform"),
		resolv.NewObject(352, 64+192, 48, 8, "platform"),
	)
}
