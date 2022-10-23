package main

import (
	"image/color"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH  = 800
	HEIGHT = 600
)

const (
	MENU = iota
	PLAY
)

type player struct {
	width  int
	height int
	x, y   float64
}

func (p *player) DrawPlayer(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, float64(p.width), float64(p.height), color.White)
}

func (p *player) UpdatePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.y > 0 {
		p.y -= 5
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && int(p.y) < HEIGHT-p.height {
		p.y += 5
	}
}

type musuh struct {
	width  int
	height int
	x, y   float64
}

func (m *musuh) DrawMusuh(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, m.x, m.y, float64(m.width), float64(m.height), color.White)
}

func (m *musuh) UpdateMusuh(b ball) {
	if m.y+float64(m.height)/2 < b.y && m.y < float64(HEIGHT)-float64(m.height) {
		m.y += 5
	}

	if m.y+float64(m.height)/2 > b.y && m.y > 0 {
		m.y += -5
	}
}

type ball struct {
	radius float64
	x, y   float64
	sx, sy int
}

func (b *ball) DrawBall(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, b.x, b.y, b.radius, color.White)
}

func (b *ball) UpdateBall(g *Game) {
	if b.x > float64(WIDTH) {
		g.musuh_score += 1
		b.x, b.y = float64(WIDTH)/2, float64(HEIGHT)/2
		b.sx *= -1
	} else if b.x < 0 {
		g.player_score += 1
		b.x, b.y = float64(WIDTH)/2, float64(HEIGHT)/2
		b.sx *= -1
	}

	if b.y > float64(HEIGHT)-b.radius || b.y < b.radius {
		b.sy *= -1
	}

	b.x += float64(b.sx)
	b.y += float64(b.sy)
}

type Game struct {
	player
	musuh
	ball
	player_score int
	musuh_score  int
	f            font.Face
	menuFont     font.Face
	subMenuFont  font.Face
	scene        int
}

func collisionBall(b ball, p player) bool {
	cx := b.x
	cy := b.y

	if b.x < p.x {
		cx = p.x
	} else if b.x > p.x+float64(p.width) {
		cx = p.x + float64(p.width)
	}

	if b.y < p.y {
		cy = p.y
	} else if b.y > p.y+float64(p.height) {
		cy = p.y + float64(p.height)
	}

	dx := b.x - cx
	dy := b.y - cy

	jarak := math.Sqrt(dx*dx + dy*dy)

	return jarak <= b.radius
}

func (g *Game) Update() error {
	// change scene to play scene
	if ebiten.IsKeyPressed(ebiten.KeySpace) && g.scene == MENU {
		g.scene = PLAY
	}

	if g.scene == PLAY {
		g.player.UpdatePlayer()
		g.musuh.UpdateMusuh(g.ball)
		g.ball.UpdateBall(g)
		if collisionBall(g.ball, g.player) || collisionBall(g.ball, player(g.musuh)) {
			g.ball.sx *= -1
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.scene {
	case MENU:
		text.Draw(screen, "PongChamp", g.menuFont, WIDTH/2-140, HEIGHT/4, color.White)
		text.Draw(screen, "Tekan \"Space\" untuk play", g.subMenuFont, WIDTH/2-155, HEIGHT/2, color.White)
		text.Draw(screen, "created by aji mustofa @pepega90", g.subMenuFont, 10, HEIGHT-20, color.White)
	case PLAY:
		g.player.DrawPlayer(screen)
		g.ball.DrawBall(screen)
		g.musuh.DrawMusuh(screen)
		for i := 0; i < HEIGHT; i++ {
			ebitenutil.DrawRect(screen, float64(WIDTH)/2, float64(i)*15, 5, 5, color.White)
		}
		// draw score
		text.Draw(screen, strconv.Itoa(g.player_score), g.f, WIDTH/2+50, 35, color.White)
		text.Draw(screen, strconv.Itoa(g.musuh_score), g.f, WIDTH/2-80, 35, color.White)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("PongChamp by Aji Mustofa")

	// load font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{}

	// player
	game.player.x = float64(WIDTH - 60)
	game.player.y = float64(HEIGHT / 2)
	game.player.width = 20
	game.player.height = 100

	// musuh
	game.musuh.x = float64(30)
	game.musuh.y = float64(HEIGHT / 2)
	game.musuh.width = 20
	game.musuh.height = 100

	// ball
	game.ball.x = float64(WIDTH) / 2
	game.ball.y = float64(HEIGHT) / 2
	game.ball.radius = 10
	game.ball.sx, game.ball.sy = 6, 6

	// other
	game.player_score, game.musuh_score = 0, 0
	game.f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    35,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	game.menuFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    50,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	game.subMenuFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    25,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
