package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ktkennychow/go-rpg/constants"
)

type Sprite struct {
	Img  *ebiten.Image
	X, Y float64
}

func (s *Sprite) Collider() constants.Rect {
	return constants.Rect{
		X: s.X,
		Y: s.Y,
		W: float64(s.Img.Bounds().Dx()),
		H: float64(s.Img.Bounds().Dy()),
	}
}

func (s *Sprite) CollidesWith(other *Sprite) bool {
	myRect := s.Collider()
	otherRect := other.Collider()

	return myRect.MaxX() > otherRect.X && myRect.X < otherRect.MaxX() && myRect.MaxY() > otherRect.Y && myRect.Y < otherRect.MaxY()
}
