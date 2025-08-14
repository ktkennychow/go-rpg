package entities

import "image"

type Collider struct {
	Owner *Sprite
	Rect  image.Rectangle
}
