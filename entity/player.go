package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pablogolobaro/samuraiDuty/animations"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
	"github.com/yohamta/ganim8/v2"
	"math"
)

const (
	Idle = iota
	Run
	Fall
	Jump
	Attack
)

const (
	Left = iota
	Right
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionAttack
)

type Player struct {
	Physics     Physics
	Action      int
	actionEnded bool
	velocity    int
	jumpPower   int
	input       *input.Handler
	animations  animations.Animations
	direction   int
}

func NewPlayer(input *input.Handler, animations animations.Animations) *Player {

	p := &Player{
		Physics:     Physics{Object: resolv.NewObject(32, 128, 16, 24, "entity")},
		Action:      Idle,
		actionEnded: false,
		input:       input,
		animations:  animations,
		direction:   Right,
	}

	return p
}

func (p *Player) Update() {
	if p.Action != Attack {
		p.Action = Idle
	}

	if p.Physics.OnGround == nil {
		p.Action = p.differActionAnimationImportance(Fall)
	}

	friction := 0.5
	accel := 0.5 + friction
	maxSpeed := 4.0
	jumpSpd := 10.0
	gravity := 0.75

	p.Physics.SpeedY += gravity

	if p.Physics.WallSliding != nil && p.Physics.SpeedY > 1 {
		p.Physics.SpeedY = 1
	}

	// Horizontal movement is only possible when not wallsliding.
	if p.Physics.WallSliding == nil {
		if p.input.ActionIsPressed(ActionMoveRight) {
			p.Physics.SpeedX += accel
			p.Action = p.differActionAnimationImportance(Run)
			if p.direction == Left {
				p.direction = Right
				p.animations.FlipAll()
			}

		}

		if p.input.ActionIsPressed(ActionMoveLeft) {
			p.Physics.SpeedX -= accel
			p.Action = p.differActionAnimationImportance(Run)
			if p.direction == Right {
				p.direction = Left
				p.animations.FlipAll()
			}

		}
	}

	// Apply friction and horizontal speed limiting.
	if p.Physics.SpeedX > friction {
		p.Physics.SpeedX -= friction
	} else if p.Physics.SpeedX < -friction {
		p.Physics.SpeedX += friction
	} else {
		p.Physics.SpeedX = 0
	}

	if p.Physics.SpeedX > maxSpeed {
		p.Physics.SpeedX = maxSpeed
	} else if p.Physics.SpeedX < -maxSpeed {
		p.Physics.SpeedX = -maxSpeed
	}

	if p.input.ActionIsPressed(ActionMoveUp) {
		p.Action = p.differActionAnimationImportance(Jump)
		if p.Physics.OnGround != nil {
			p.Physics.SpeedY = -jumpSpd
		} else if p.Physics.WallSliding != nil {
			// WALLJUMPING
			p.Physics.SpeedY = -jumpSpd

			if p.Physics.WallSliding.X > p.Physics.Object.X {
				p.Physics.SpeedX = -4
			} else {
				p.Physics.SpeedX = 4
			}

			p.Physics.WallSliding = nil

		}
	}

	// dx is the horizontal delta movement variable (which is the Player's horizontal speed). If we come into contact with something, then it will
	// be that movement instead.
	dx := p.Physics.SpeedX

	// Moving horizontally is done fairly simply; we just check to see if something solid is in front of us. If so, we move into contact with it
	// and stop horizontal movement speed. If not, then we can just move forward.

	if check := p.Physics.Object.Check(p.Physics.SpeedX, 0, "solid"); check != nil {

		dx = check.ContactWithCell(check.Cells[0]).X()
		p.Physics.SpeedX = 0

		// If you're in the air, then colliding with a wall object makes you start wall sliding.
		if p.Physics.OnGround == nil {
			p.Physics.WallSliding = check.Objects[0]
		}
	}

	// Then we just apply the horizontal movement to the Player's Object. Easy-peasy.
	p.Physics.Object.X += dx

	// Now for the vertical movement; it's the most complicated because we can land on different types of objects and need
	// to treat them all differently, but overall, it's not bad.

	// First, we set OnGround to be nil, in case we don't end up standing on anything.
	p.Physics.OnGround = nil

	// dy is the delta movement downward, and is the vertical movement by default; similarly to dx, if we come into contact with
	// something, this will be changed to move to contact instead.

	dy := p.Physics.SpeedY

	// We want to be sure to lock vertical movement to a maximum of the size of the Cells within the Space
	// so we don't miss any collisions by tunneling through.

	dy = math.Max(math.Min(dy, 16), -16)

	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	if check := p.Physics.Object.Check(0, checkDistance, "solid", "platform", "ramp"); check != nil {

		// So! Firstly, we want to see if we jumped up into something that we can slide around horizontally to avoid bumping the Player's head.

		// Sliding around a misspaced jump is a small thing that makes jumping a bit more forgiving, and is something different polished platformers
		// (like the 2D Mario games) do to make it a smidge more comfortable to play. For a visual example of this, see this excellent devlog post
		// from the extremely impressive indie game, Leilani's Island: https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138

		// To accomplish this sliding, we simply call Collision.SlideAgainstCell() to see if we can slide.
		// We pass the first cell, and tags that we want to avoid when sliding (i.e. we don't want to slide into cells that contain other solid objects).

		slide := check.SlideAgainstCell(check.Cells[0], "solid")

		// We further ensure that we only slide if:
		// 1) We're jumping up into something (dy < 0),
		// 2) If the cell we're bumping up against contains a solid object,
		// 3) If there was, indeed, a valid slide left or right, and
		// 4) If the proposed slide is less than 8 pixels in horizontal distance. (This is a relatively arbitrary number that just so happens to be half the
		// width of a cell. This is to ensure the player doesn't slide too far horizontally.)

		if dy < 0 && check.Cells[0].ContainsTags("solid") && slide != nil && math.Abs(slide.X()) <= 8 {

			// If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
			p.Physics.Object.X += slide.X()

		} else {

			// If sliding -fails-, that means the Player is jumping directly onto or into something, and we need to do more to see if we need to come into
			// contact with it. Let's press on!

			// First, we check for ramps. For ramps, we can't simply check for collision with Check(), as that's not precise enough. We need to get a bit
			// more information, and so will do so by checking its Shape (a triangular ConvexPolygon, as defined in WorldPlatformer.Init()) against the
			// Player's Shape (which is also a rectangular ConvexPolygon).

			// We get the ramp by simply filtering out Objects with the "ramp" tag out of the objects returned in our broad Check(), and grabbing the first one
			// if there's any at all.
			if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {

				ramp := ramps[0]

				// For simplicity, this code assumes we can only stand on one ramp at a time as there is only one ramp in this example.
				// In actuality, if there was a possibility to have a potential collision with multiple ramps (i.e. a ramp that sits on another ramp, and the player running down
				// one onto the other), the collision testing code should probably go with the ramp with the highest confirmed intersection point out of the two.

				// Next, we see if there's been an intersection between the two Shapes using Shape.Intersection. We pass the ramp's shape, and also the movement
				// we're trying to make horizontally, as this makes Intersection return the next y-position while moving, not the one directly
				// underneath the Player. This would keep the player from getting "stuck" when walking up a ramp into the top of a solid block, if there weren't
				// a landing at the top and bottom of the ramp.

				// We use 8 here for the Y-delta so that we can easily see if you're running down the ramp (in which case you're probably in the air as you
				// move faster than you can fall in this example). This way we can maintain contact so you can always jump while running down a ramp. We only
				// continue with coming into contact with the ramp as long as you're not moving upwards (i.e. jumping).

				if contactSet := p.Physics.Object.Shape.Intersection(dx, 8, ramp.Shape); dy >= 0 && contactSet != nil {

					// If Intersection() is successful, a ContactSet is returned. A ContactSet contains information regarding where
					// two Shapes intersect, like the individual points of contact, the center of the contacts, and the MTV, or
					// Minimum Translation Vector, to move out of contact.

					// Here, we use ContactSet.TopmostPoint() to get the top-most contact point as an indicator of where
					// we want the player's feet to be. Then we just set that position, and we're done.

					dy = contactSet.TopmostPoint()[1] - p.Physics.Object.Bottom() + 0.1
					p.Physics.OnGround = ramp
					p.Physics.SpeedY = 0

				}

			}

			// Finally, we check for simple solid ground. If we haven't had any success in landing previously, or the solid ground
			// is higher than the existing ground (like if the platform passes underneath the ground, or we're walking off of solid ground
			// onto a ramp), we stand on it instead. We don't check for solid collision first because we want any ramps to override solid
			// ground (so that you can walk onto the ramp, rather than sticking to solid ground).

			// We use ContactWithObject() here because otherwise, we might come into contact with the moving platform's cells (which, naturally,
			// would be selected by a Collision.ContactWithCell() call because the cell is closest to the Player).

			if solids := check.ObjectsByTags("solid"); len(solids) > 0 && (p.Physics.OnGround == nil || p.Physics.OnGround.Y >= solids[0].Y) {
				dy = check.ContactWithCell(check.Cells[0]).Y()
				p.Physics.SpeedY = 0

				// We're only on the ground if we land on it (if the object's Y is greater than the player's).
				if solids[0].Y > p.Physics.Object.Y {
					p.Physics.OnGround = solids[0]
				}

			}

			if p.Physics.OnGround != nil {
				p.Physics.WallSliding = nil // Player's on the ground, so no wallsliding anymore.
			}

		}

	}

	// Move the object on dy.
	p.Physics.Object.Y += dy

	wallNext := 1.0
	if p.direction != Right {
		wallNext = -1
	}

	// If the wall next to the Player runs out, stop wall sliding.
	if c := p.Physics.Object.Check(wallNext, 0, "solid"); p.Physics.WallSliding != nil && c == nil {
		p.Physics.WallSliding = nil
	}

	p.Physics.Object.Update() // Update the player's position in the space.

	if p.input.ActionIsJustPressed(ActionAttack) {
		p.Action = Attack
	}

	switch p.Action {
	case Attack:
		p.animations.Attack.Update()

		if p.animations.Attack.Sprite().IsEnd(p.animations.Attack.Position()) {
			p.actionEnded = true
		}
		if p.animations.Attack.Position() == 1 && p.actionEnded {
			p.Action = Idle
			p.actionEnded = false
		}
	case Jump:

		p.animations.Jump.Update()
	case Idle:
		p.animations.Idle.Update()
	case Run:
		p.animations.Run.Update()

	}

}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.Action == Run {
		p.animations.Run.Draw(screen, ganim8.DrawOpts(p.Physics.Object.X, p.Physics.Object.Y, 0, 1, 1, 0.5, 0.5))
	} else if p.Action == Attack {
		p.animations.Attack.Draw(screen, ganim8.DrawOpts(p.Physics.Object.X, p.Physics.Object.Y, 0, 1, 1, 0.5, 0.5))
	} else if p.Action == Jump {
		p.animations.Jump.Draw(screen, ganim8.DrawOpts(p.Physics.Object.X, p.Physics.Object.Y, 0, 1, 1, 0.5, 0.5))
	} else if p.Action == Fall {
		p.animations.Fall.Draw(screen, ganim8.DrawOpts(p.Physics.Object.X, p.Physics.Object.Y, 0, 1, 1, 0.5, 0.5))
	} else {
		p.animations.Idle.Draw(screen, ganim8.DrawOpts(p.Physics.Object.X, p.Physics.Object.Y, 0, 1, 1, 0.5, 0.5))
	}

}

func (p *Player) differActionAnimationImportance(actionIndex int) int {
	if p.Action < actionIndex {
		return actionIndex
	} else {
		return p.Action
	}
}
