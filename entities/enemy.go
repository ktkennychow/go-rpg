package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ktkennychow/go-rpg/constants"
	"github.com/ktkennychow/go-rpg/utils"
)

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

func (e *Enemy) Update(player *Player) {
	if e.X > player.X {
		e.X -= constants.EnemySpeed
	}
	if e.X < player.X {
		e.X += constants.EnemySpeed
	}
	if e.Y > player.Y {
		e.Y -= constants.EnemySpeed
	}
	if e.Y < player.Y {
		e.Y += constants.EnemySpeed
	}
}

func (e *Enemy) Draw(screen *ebiten.Image, cam *utils.Camera) {
	opts := ebiten.DrawImageOptions{}

	opts.GeoM.Translate(e.X, e.Y)
	opts.GeoM.Translate(cam.OffsetX, cam.OffsetY)
	screen.DrawImage(e.Img, &opts)
	opts.GeoM.Reset()
}
