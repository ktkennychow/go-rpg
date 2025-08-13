package utils

import "github.com/ktkennychow/go-rpg/constants"

type Camera struct {
	OffsetX, OffsetY float64
}

func NewCamera(originX, originY float64) *Camera {
	return &Camera{
		OffsetX: originX,
		OffsetY: originY,
	}
}

func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.OffsetX = screenWidth/2 - (targetX + constants.TileSize/2)
	c.OffsetY = screenHeight/2 - (targetY + constants.TileSize/2)
}

func (c *Camera) Constrain(tilemapWidthPx, tilemapHeightPx, screenWidth, screenHeight float64) {
	if c.OffsetX > 0 {
		c.OffsetX = 0
	}
	if c.OffsetY > 0 {
		c.OffsetY = 0
	}
	// Prevent showing right of world
	if c.OffsetX < screenWidth-tilemapWidthPx {
		c.OffsetX = screenWidth - tilemapWidthPx
	}

	// Prevent showing below world
	if c.OffsetY < screenHeight-tilemapHeightPx {
		c.OffsetY = screenHeight - tilemapHeightPx
	}
}
