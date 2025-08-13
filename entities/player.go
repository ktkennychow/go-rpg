package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ktkennychow/go-rpg/constants"
	"github.com/ktkennychow/go-rpg/utils"
)

type Player struct {
	*Sprite
	Health uint
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.X -= constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.X += constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.Y -= constants.PlayerSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.Y += constants.PlayerSpeed
	}
}

func (p *Player) Draw(screen *ebiten.Image, cam *utils.Camera) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.X, p.Y)
	opts.GeoM.Translate(cam.OffsetX, cam.OffsetY)
	screen.DrawImage(p.Img, &opts)
	opts.GeoM.Reset()

	playerRect := p.Collider()
	vector.StrokeRect(screen, float32(playerRect.X+cam.OffsetX), float32(playerRect.Y+cam.OffsetY), float32(playerRect.W), float32(playerRect.H), 1, color.RGBA{255, 0, 0, 255}, true)
}
