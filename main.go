package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ktkennychow/go-rpg/assets"
	"github.com/ktkennychow/go-rpg/constants"
	"github.com/ktkennychow/go-rpg/entities"
	"github.com/ktkennychow/go-rpg/maps"
	"github.com/ktkennychow/go-rpg/utils"
)

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	colliders   []*entities.Collider
	tilemapJSON *maps.TilemapJSON
	tileSets    []maps.Tileset
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
		if g.colliders[i].Rect.Overlaps(image.Rect(int(g.player.X), int(g.player.Y), int(g.player.X+constants.TileSize), int(g.player.Y+constants.TileSize))) {
			g.player.Health += potion.AmtHeal
			g.potions = append(g.potions[:i], g.potions[i+1:]...)
			g.colliders = append(g.colliders[:i], g.colliders[i+1:]...)
			fmt.Printf("Potion collected: %d\n", g.player.Health)
		}
	}
	// global collider update
	for i := range g.colliders {
		// update the collider's position to the owner's position
		g.colliders[i].Rect = image.Rect(int(g.colliders[i].Owner.X), int(g.colliders[i].Owner.Y), int(g.colliders[i].Owner.X+float64(g.colliders[i].Owner.Img.Bounds().Dx())), int(g.colliders[i].Owner.Y+float64(g.colliders[i].Owner.Img.Bounds().Dy())))
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
	}

	g.player.Draw(screen, g.cam)

	for _, enemy := range g.enemies {
		enemy.Draw(screen, g.cam)
	}

	for _, collider := range g.colliders {
		x := float64(collider.Rect.Min.X) + g.cam.OffsetX
		y := float64(collider.Rect.Min.Y) + g.cam.OffsetY
		w := float64(collider.Rect.Dx())
		h := float64(collider.Rect.Dy())
		opts.GeoM.Translate(x, y)
		vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, color.RGBA{255, 0, 0, 255}, true)
		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

var playerImg = assets.MustLoadImage("Actor/Characters/FighterRed/SpriteSheet.png")

func (g *Game) SpawnPlayer() {
	g.player = &entities.Player{
		Sprite: &entities.Sprite{
			Img: playerImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
			X:   constants.ScreenWidth/2 - constants.TileSize/2,
			Y:   constants.ScreenHeight/2 - constants.TileSize/2,
		},
		Health: constants.PlayerStartingHealth,
	}
}

var potionImg = assets.MustLoadImage("Items/Potion/LifePot.png")

func (g *Game) SpawnPotion() {
	newPotion := &entities.Potion{
		Sprite: &entities.Sprite{
			Img: potionImg,
			X:   200.0,
			Y:   200.0,
		},
		AmtHeal: constants.PotionHealAmount,
	}
	g.potions = append(g.potions, newPotion)
	g.colliders = append(g.colliders, &entities.Collider{
		Owner: newPotion.Sprite,
		Rect:  image.Rect(int(newPotion.X), int(newPotion.Y), int(newPotion.X+float64(newPotion.Sprite.Img.Bounds().Dx())), int(newPotion.Y+float64(newPotion.Sprite.Img.Bounds().Dy()))),
	})
}

var skeletonImg = assets.MustLoadImage("Actor/Characters/Skeleton/SpriteSheet.png")

func (g *Game) SpawnEnemy() {
	newEnemy := &entities.Enemy{
		Sprite: &entities.Sprite{
			Img: skeletonImg.SubImage(image.Rect(0, 0, constants.TileSize, constants.TileSize)).(*ebiten.Image),
			X:   rand.Float64() * float64(constants.ScreenWidth-constants.TileSize),
			Y:   rand.Float64() * float64(constants.ScreenHeight-constants.TileSize),
		},
		FollowsPlayer: true,
	}
	g.enemies = append(g.enemies, newEnemy)
	g.colliders = append(g.colliders, &entities.Collider{
		Owner: newEnemy.Sprite,
		Rect:  image.Rect(int(newEnemy.X), int(newEnemy.Y), int(newEnemy.X+float64(newEnemy.Sprite.Img.Bounds().Dx())), int(newEnemy.Y+float64(newEnemy.Sprite.Img.Bounds().Dy()))),
	})
}

var tilemapImg = assets.MustLoadImage("Backgrounds/Tilesets/TilesetFloor.png")

func main() {
	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	tilemapJSON, err := maps.NewTilemapJSON("maps/spawn_map.json")
	if err != nil {
		log.Fatal(err)
	}

	tileSets, err := tilemapJSON.GetTilesets()
	if err != nil {
		log.Fatal(err)
	}

	game := Game{

		tilemapImg:  tilemapImg,
		tilemapJSON: tilemapJSON,
		tileSets:    tileSets,
		cam:         utils.NewCamera(0, 0),
	}

	game.SpawnPlayer()
	game.SpawnPotion()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()
	game.SpawnEnemy()

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
	select {}
}
