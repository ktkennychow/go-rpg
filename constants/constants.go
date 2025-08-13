package constants

// Game constants
const (
	Scale        = 2.0
	TileSize     = 16
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Player constants
const (
	PlayerSpeed          = 2.0
	PlayerStartingHealth = 100
)

// Enemy constants
const (
	EnemySpeed = 0.5
)

// Item constants
const (
	PotionHealAmount = 10
)

// Rect represents a rectangle with position and dimensions
type Rect struct {
	X, Y, W, H float64
}

// MaxX returns the maximum X coordinate of the rectangle
func (r *Rect) MaxX() float64 {
	return r.X + r.W
}

// MaxY returns the maximum Y coordinate of the rectangle
func (r *Rect) MaxY() float64 {
	return r.Y + r.H
}
