package entity

import "github.com/solarlune/resolv"

type Physics struct {
	Object      *resolv.Object
	SpeedX      float64
	SpeedY      float64
	OnGround    *resolv.Object
	WallSliding *resolv.Object
}
