package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ktkennychow/go-rpg/constants"
	"github.com/ktkennychow/go-rpg/entities"
	"github.com/ktkennychow/go-rpg/utils"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tileSets    []Tileset
	tilemapImg  *ebiten.Image
	cam         *utils.Camera
}

func (g *Game) Update() error {
	g.player.Update()

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			enemy.Update(g.player)
		}
	}

	for i, potion := range g.potions {
		if potion.CollidesWith(g.player.Sprite) {
			g.player.Health += potion.AmtHeal
			g.potions = append(g.potions[:i], g.potions[i+1:]...)
			fmt.Printf("Potion collected: %d\n", g.player.Health)
		}
	}

	g.cam.FollowTarget(g.player.X, g.player.Y, constants.ScreenWidth, constants.ScreenHeight)
	g.cam.Constrain(float64(g.tilemapJSON.Width*constants.TileSize), float64(g.tilemapJSON.Height*constants.TileSize), constants.ScreenWidth, constants.ScreenHeight)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// sky blue
	screen.Fill(color.RGBA{135, 206, 235, 255})

	opts := ebiten.DrawImageOptions{}

	for layerIdx, layer := range g.tilemapJSON.Layers {
		for idx, tileid := range layer.Data {
			// skip empty tiles
			if tileid == 0 {
				continue
			}

			img := g.tileSets[layerIdx].Img(tileid)
			X := idx % layer.Width * constants.TileSize
			//int division loses decimal precision = flooring
			Y := idx / layer.Width * constants.TileSize

			opts.GeoM.Translate(float64(X), float64(Y))
			// if tile is larger than the tile size, we need to translate the tile up by the difference
			opts.GeoM.Translate(0.0, -(float64(img.Bounds().Dy()))+constants.TileSize)
			opts.GeoM.Translate(g.cam.OffsetX, g.cam.OffsetY)
			screen.DrawImage(img, &opts)
			opts.GeoM.Reset()

			// srcX := (tileid-1)%22*constants.TileSize + 1
			// //int division loses decimal precision = flooring
			// srcY := (tileid-1)/22*constants.TileSize + 1
			// screen.DrawImage(g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+constants.TileSize, srcY+constants.TileSize)).(*ebiten.Image), &opts)
			// opts.GeoM.Reset()
		}
	}

	for _, potion := range g.potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(g.cam.OffsetX, g.cam.OffsetY)
		screen.DrawImage(potion.Img, &opts)
		opts.GeoM.Reset()

		potionRect := potion.Collider()
		vector.StrokeRect(screen, float32(potionRect.X+g.cam.OffsetX), float32(potionRect.Y+g.cam.OffsetY), float32(potionRect.W), float32(potionRect.H), 1, color.RGBA{255, 0, 0, 255}, true)
	}

	g.player.Draw(screen, g.cam)

	for _, enemy := range g.enemies {
		enemy.Draw(screen, g.cam)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/Actor/Characters/FighterRed/SpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/Actor/Characters/Skeleton/SpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	potionImg, _, err := ebitenutil.NewImageFromFile("assets/Items/Potion/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/Backgrounds/Tilesets/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/spawn_map.json")
	if err != nil {
		log.Fatal(err)
	}

	tileSets, err := tilemapJSON.GetTilesets()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
				X:   constants.ScreenWidth/2 - constants.TileSize/2,
				Y:   constants.ScreenHeight/2 - constants.TileSize/2,
			},
			Health: constants.PlayerStartingHealth,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
					X:   44.0,
					Y:   34.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
					X:   544.0,
					Y:   34.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
					X:   44.0,
					Y:   234.0,
				},
				FollowsPlayer: false,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   200.0,
					Y:   200.0,
				},
				AmtHeal: constants.PotionHealAmount,
			},
		},
		tilemapImg:  tilemapImg,
		tilemapJSON: tilemapJSON,
		tileSets:    tileSets,
		cam:         utils.NewCamera(0, 0),
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
