package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var backgroundColor = color.RGBA{0xfa, 0xf8, 0xef, 0xff}
var frameColor = color.RGBA{0xbb, 0xad, 0xa0, 0xff}
var tileSize = 80
var tileMargin = 4
var mplusFaceSource *text.GoTextFaceSource

func tileColor(value int) color.Color {
	switch value {
	case 2, 4:
		return color.RGBA{0x77, 0x6e, 0x65, 0xff}
	case 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536:
		return color.RGBA{0xf9, 0xf6, 0xf2, 0xff}
	}
	panic("not reach")
}

func tileBackgroundColor(value int) color.Color {
	switch value {
	case 0:
		return color.NRGBA{0xee, 0xe4, 0xda, 0x59}
	case 2:
		return color.RGBA{0xee, 0xe4, 0xda, 0xff}
	case 4:
		return color.RGBA{0xed, 0xe0, 0xc8, 0xff}
	case 8:
		return color.RGBA{0xf2, 0xb1, 0x79, 0xff}
	case 16:
		return color.RGBA{0xf5, 0x95, 0x63, 0xff}
	case 32:
		return color.RGBA{0xf6, 0x7c, 0x5f, 0xff}
	case 64:
		return color.RGBA{0xf6, 0x5e, 0x3b, 0xff}
	case 128:
		return color.RGBA{0xed, 0xcf, 0x72, 0xff}
	case 256:
		return color.RGBA{0xed, 0xcc, 0x61, 0xff}
	case 512:
		return color.RGBA{0xed, 0xc8, 0x50, 0xff}
	case 1024:
		return color.RGBA{0xed, 0xc5, 0x3f, 0xff}
	case 2048:
		return color.RGBA{0xed, 0xc2, 0x2e, 0xff}
	case 4096:
		return color.NRGBA{0xa3, 0x49, 0xa4, 0x7f}
	case 8192:
		return color.NRGBA{0xa3, 0x49, 0xa4, 0xb2}
	case 16384:
		return color.NRGBA{0xa3, 0x49, 0xa4, 0xcc}
	case 32768:
		return color.NRGBA{0xa3, 0x49, 0xa4, 0xe5}
	case 65536:
		return color.NRGBA{0xa3, 0x49, 0xa4, 0xff}
	}
	panic("not reach")
}

type Tile struct {
	value int
	x     int
	y     int
}

type Board struct {
	grid           [][]Tile
	size           int
	foundCollision bool
}

type Game struct {
	boardImage *ebiten.Image
	board      *Board
	isPressed  bool
	gameOver   bool
}

func (g *Game) addRandomTile() {
	if !g.board.foundCollision {
		return
	}
	var randX, randY int
	for {
		randX = rand.Intn(4)
		randY = rand.Intn(4)
		if g.board.grid[randX][randY].value == 0 {
			g.board.grid[randX][randY].value = 2
			break
		}
	}
}

func (b *Board) isBoardFull() bool {
	for x := range b.grid {
		for y := range b.grid {
			if b.grid[x][y].value == 0 {
				return false
			}
		}
	}
	return true
}

func (b *Board) checkCollisons(first, second *Tile) bool {
	//value := first.value
	log.Printf("check %v, %v", first, second)

	if second.value == 0 {
		second.value = first.value
		first.value = 0
		b.foundCollision = true
		return true
	} else if first.value > 0 && first.value == second.value {
		second.value = second.value + first.value
		first.value = 0
		b.foundCollision = true
		return true
	}
	return false
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		log.Printf("right")
		for x := g.board.size - 1; x >= 0; x-- {
			for y := g.board.size - 1; y >= 0; y-- {
				log.Printf("[%v, %v] value %v", x, y, g.board.grid[x][y].value)
				if g.board.grid[x][y].value == 0 {
					continue
				}
				for i := g.board.size - 1; i >= y; i-- {
					log.Printf("x %v, y %v, i %v", x, y, i)
					if y == i {
						continue
					}
					moved := g.board.checkCollisons(&g.board.grid[x][y], &g.board.grid[x][i])
					if moved {
						log.Printf("moved to [%v, %v]", x, i)
						break
					}
				}
			}
		}
		g.addRandomTile()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		for x := 0; x < g.board.size; x++ {
			for y := 0; y < g.board.size; y++ {
				if g.board.grid[x][y].value == 0 {
					continue
				}
				for i := 0; i < g.board.size-1; i++ {
					if i > y {
						break
					}
					if y == i {
						continue
					}
					moved := g.board.checkCollisons(&g.board.grid[x][y], &g.board.grid[x][i])
					if moved {
						break
					}
				}
			}
		}
		g.addRandomTile()

	} else if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		log.Printf("pressed up")
		for x := 0; x < g.board.size; x++ {
			for y := 0; y < g.board.size; y++ {
				log.Printf("[%v, %v] value %v", x, y, g.board.grid[x][y].value)
				if g.board.grid[x][y].value == 0 {
					continue
				}
				for i := 0; i < g.board.size-1; i++ {
					log.Printf("x %v y%x i %v", x, y, i)
					if i > x {
						log.Printf("break")
						break
					}
					if x == i {
						continue
					}
					moved := g.board.checkCollisons(&g.board.grid[x][y], &g.board.grid[i][y])
					if moved {
						log.Printf("moved to [%v, %v]", i, y)
						break
					}
				}
				//				log.Printf("check [%v][%v]", i, j)

			}
		}
		g.addRandomTile()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		log.Printf("down")
		for y := g.board.size - 1; y >= 0; y-- {
			for x := g.board.size - 1; x >= 0; x-- {
				log.Printf("[%v, %v] value %v", x, y, g.board.grid[x][y].value)
				if g.board.grid[x][y].value == 0 {
					continue
				}
				for i := g.board.size - 1; i >= x; i-- {
					log.Printf("x %v y%x i %v", x, y, i)
					if x == i {
						continue
					}
					moved := g.board.checkCollisons(&g.board.grid[x][y], &g.board.grid[i][y])
					if moved {
						log.Printf("moved to [%v, %v]", i, y)
						break
					}
				}
				//				log.Printf("check [%v][%v]", i, j)

			}
		}
		g.addRandomTile()
	}
	if g.board.isBoardFull() {
		if !g.board.foundCollision {
			g.gameOver = true
		}
	}
	g.board.foundCollision = false // reset for next update

	return nil
}

func (t *Tile) Draw(boardImage *ebiten.Image) {
	i, j := t.y, t.x
	v := t.value
	if v == 0 {
		return
	}
	//log.Printf("draw tile [%v, %v] value: %v", j, i, v)
	tileImg := ebiten.NewImage(tileSize, tileSize)

	// Fill it with the background color
	tileImg.Fill(tileBackgroundColor(v))

	op := &ebiten.DrawImageOptions{}
	x := i*tileSize + (i+1)*tileMargin
	y := j*tileSize + (j+1)*tileMargin
	op.GeoM.Translate(float64(x), float64(y))

	// Draw the colored tile
	boardImage.DrawImage(tileImg, op)

	str := strconv.Itoa(v)

	size := 48.0
	switch {
	case 3 < len(str):
		size = 24
	case 2 < len(str):
		size = 32
	}

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(x)+float64(tileSize)/2, float64(y)+float64(tileSize)/2)
	textOp.ColorScale.ScaleWithColor(tileColor(v))
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	text.Draw(boardImage, str, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   size,
	}, textOp)
}

func (b *Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			v := 0
			// Create a new image for this specific tile
			tileImg := ebiten.NewImage(tileSize, tileSize)

			// Fill it with the background color
			tileImg.Fill(tileBackgroundColor(v))

			op := &ebiten.DrawImageOptions{}
			x := i*tileSize + (i+1)*tileMargin
			y := j*tileSize + (j+1)*tileMargin
			op.GeoM.Translate(float64(x), float64(y))
			boardImage.DrawImage(tileImg, op)
		}
	}
	for i := range b.grid {
		for j := range b.grid[i] {
			b.grid[i][j].Draw(boardImage)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		//		size := g.board.size*tileSize + (g.board.size+1)*tileMargin
		g.boardImage = ebiten.NewImage(340, 340)
	}
	if g.gameOver {
		// todo: implement game over screen
		return
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	bw, bh := g.boardImage.Bounds().Dx(), g.boardImage.Bounds().Dy()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
	if g.isPressed {
		ebitenutil.DebugPrint(screen, "Hello, World!")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 420, 600
}

func main() {
	ebiten.SetWindowSize(420, 600)
	ebiten.SetWindowTitle("Hello, World!")

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	rand.Seed(time.Now().UnixNano())
	game := &Game{}
	board := &Board{}
	board.size = 4
	grid := make([][]Tile, 4)
	for i := range grid {
		grid[i] = make([]Tile, 4)
	}
	for i := range grid {
		for j := range grid[i] {
			tile := Tile{
				x:     i,
				y:     j,
				value: 0,
			}
			grid[i][j] = tile
		}
	}
	grid[0][0].value = 2
	grid[1][0].value = 4
	grid[1][1].value = 2

	game.board = board
	board.grid = grid

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
