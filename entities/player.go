package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ktkennychow/go-rpg/constants"
	"github.com/ktkennychow/go-rpg/utils"
)

type Player struct {
	*Sprite
	Health uint
}

func (p *Player) Update() {
	// reset velocity
	p.VelocityX = 0
	p.VelocityY = 0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.VelocityX = -constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.VelocityX = constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.VelocityY = -constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.VelocityY = constants.PlayerSpeed
	}

	p.X += p.VelocityX
	p.Y += p.VelocityY
}

func (p *Player) Draw(screen *ebiten.Image, cam *utils.Camera) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.OffsetX, cam.OffsetY)
	screen.DrawImage(p.Img, &opts)
	opts.GeoM.Reset()
}
