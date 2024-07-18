package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PanelOptions struct {
	height    int
	width     int
	x         int
	y         int
	lineThick int
	text      string
	tracker   string
}

// func createPanel(screen *ebiten.Image, ops PanelOptions) {
//     panelImage := ebiten.NewImage(ops.width, ops.height)

//     vector.StrokeRect(
//         panelImage,
//         0,
//         0,
//         float32(ops.width),
//         float32(ops.height),
//         float32(2),
//         color.Black,
//         true,
//     )

//     panel := &ebiten.DrawImageOptions{}
//     panel.GeoM.Translate(float64(ops.x), float64(ops.y))

//	    screen.DrawImage(panelImage, panel)
//	}
func createPanel(screen *ebiten.Image, ops PanelOptions) {
	panelImage := ebiten.NewImage(ops.width, ops.height)

	// Draw the panel outline
	vector.StrokeRect(
		panelImage,
		0,
		0,
		float32(ops.width),
		float32(ops.height),
		float32(ops.lineThick),
		color.Black,
		true,
	)

	// Set up text options
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(ops.x)+20, float64(ops.height)/2)
	textOp.ColorScale.ScaleWithColor(color.Black)
	textOp.PrimaryAlign = text.AlignStart
	textOp.SecondaryAlign = text.AlignCenter

	// Draw text on the panel image
	text.Draw(panelImage, ops.text, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   float64(ops.height) / 3, // Adjust size as needed
	}, textOp)

	// Draw the panel image on the screen
	panel := &ebiten.DrawImageOptions{}
	panel.GeoM.Translate(float64(ops.x), float64(ops.y))
	screen.DrawImage(panelImage, panel)
}
