//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"time"
//	"math"
	"image/color"
//	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/ByteArena/box2d"
)

type box struct {
	pad float32

	score1 int
	score2 int

	output  *widget.RichText
	window  fyne.Window
	game    *fyne.Container
	world   box2d.B2World
}

func (c *box) onTypedKey(ev *fyne.KeyEvent) {

}

func (c *box) moveOneStep() {
}

func (c * box) Refresh() {

}

func (c *box) onKeyUp(key *fyne.KeyEvent) {

}

func (c *box) onKeyDown(key *fyne.KeyEvent) {

}

type GameLayout struct {
	cx float32
	cy float32
}

func (l *GameLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(l.cx, l.cy)
}

func (l *GameLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {

}

func NewGameLayout(cx float32,cy float32) (*GameLayout) {
	return &GameLayout{cx : cx, cy : cy }
}

func (c *box) loadUI(app fyne.App) {

	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	c.world = box2d.MakeB2World(gravity)

	// make score layout
	score1 := canvas.NewText("P1 Score : 0", color.Black)
	score2 := canvas.NewText("P2 Score : 0", color.Black)
	scoreLayout := container.New(layout.NewHBoxLayout(), score1, layout.NewSpacer(), score2)

	// make new window
	var cx float32 = 600
	var cy float32 = 300
	c.window = app.NewWindow("Box")
	c.game = container.New(NewGameLayout(cx,cy))
	c.window.SetContent(container.New(layout.NewVBoxLayout(),scoreLayout,c.game))
	c.window.Content().Refresh()
	c.window.Show()

	// setup colors
//	gray := color.Gray{Y: 0x99}
//	red  := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
//	blue := color.NRGBA{R: 0x33, G: 0x33, B: 0xff, A: 0xff}
//	yellow := color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
//	green := color.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}


	deskCanvas, ok := c.window.Canvas().(desktop.Canvas)
	
	// main loop
	if ok {
		deskCanvas.SetOnKeyDown(c.onKeyDown)
		deskCanvas.SetOnKeyUp(c.onKeyUp)
		
		update := time.NewTicker(5*time.Millisecond)    // about 200 fps physics update rate
		ticker := time.NewTicker(20 * time.Millisecond) // means 50 fps.
    		done := make(chan bool)

		go func() {
		        for {
	            		select {
	            			case <-done:
                				return
            				case <-update.C:
                				c.moveOneStep()
            				case <-ticker.C:
						c.Refresh()
            				}
        		}
    		}()
		
	} else {
		canvas := c.window.Canvas()
		canvas.SetOnTypedKey(c.onTypedKey)
	}
	c.window.Show()
}

func newBox() *box {
	return &box{pad:1}
}
