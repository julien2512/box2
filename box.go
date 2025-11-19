//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"time"
	"image/color"
//	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/ByteArena/box2d"
)

type Circle struct {
	x float64
	y float64
	r float64

	pad float64

	downKey fyne.KeyName
	upKey	fyne.KeyName
	leftKey	fyne.KeyName
	rightKey fyne.KeyName

	body *box2d.B2Body

	circle *canvas.Circle
}

func (c *Circle) Color(color color.Color, wb float32, wc color.Color) {
	if c.circle != nil {
		c.circle.FillColor = color
		c.circle.StrokeWidth = wb
		c.circle.StrokeColor = wc
	}
}

func (c *Circle) Refresh() {
	c.x = c.body.GetPosition().X
	c.y = c.body.GetPosition().Y

	c.circle.Move(fyne.NewPos(float32(c.x),float32(c.y)))
	//c.circle.Refresh()
}

func (c *Circle) onTypedKey(ev *fyne.KeyEvent) {

}

func (c *Circle) onKeyUp(key *fyne.KeyEvent) {
	if key.Name == c.leftKey {
		//c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(c.pad,0),true)
	} else if key.Name == c.rightKey {
		//c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(-c.pad,0),true)
	} else if key.Name == c.upKey {
		//c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0,c.pad),true)
	} else if key.Name == c.downKey {
		//c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0,-c.pad),true)
	}
}

func (c *Circle) onKeyDown(key *fyne.KeyEvent) {
	if key.Name == c.leftKey {
		c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(-c.pad,0),true)
	} else if key.Name == c.rightKey {
		c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(c.pad,0),true)
	} else if key.Name == c.upKey {
		c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0,-c.pad),true)
	} else if key.Name == c.downKey {
		c.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(0,c.pad),true)
	}
}

func (c *Circle) CanMove(downKey fyne.KeyName, upKey fyne.KeyName, leftKey fyne.KeyName, rightKey fyne.KeyName) {
	c.downKey = downKey
	c.upKey = upKey
	c.leftKey = leftKey
	c.rightKey = rightKey
}

func NewCircle(size float64,xb float64, yb float64, body *box2d.B2Body) *Circle {
	circle := &Circle{}
	circle.pad = 1
	circle.r = size
	circle.x = xb
	circle.y = yb

	circle.circle = canvas.NewCircle(color.Black)
	circle.circle.Resize(fyne.NewSize(float32(2*size),float32(2*size)))

	circle.body = body

	return circle
}

type box struct {
	pad float32

	score1 int
	score2 int

	circle *Circle
	circle2 *Circle

	circles []*Circle

	output  *widget.RichText
	window  fyne.Window
	game    *fyne.Container
	world   box2d.B2World
}


func (c *box) onTypedKey(ev *fyne.KeyEvent) {
	c.circle.onTypedKey(ev)
	for i:=0;i<len(c.circles);i++ { 
		c.circles[i].onTypedKey(ev)
	}
}

func (c *box) Refresh() {
	if (c.circle!=nil) {
		c.circle.Refresh()
	}
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].Refresh()
	}
}

func (c *box) onKeyUp(key *fyne.KeyEvent) {
	if (c.circle2!=nil) {
		c.circle2.onKeyUp(key)
	}
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].onKeyUp(key)
	}
}

func (c *box) onKeyDown(key *fyne.KeyEvent) {
	c.circle.onKeyDown(key)
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].onKeyDown(key)
	}
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

	// make new window
	var cx float32 = 600
	var cy float32 = 300
	c.window = app.NewWindow("Box")
	c.game = container.New(NewGameLayout(cx,cy))
	c.window.SetContent(container.New(layout.NewVBoxLayout(),c.game))
	c.window.Content().Refresh()
	c.window.Show()

	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)

	// Ground body
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, 300.0), box2d.MakeB2Vec2(600.0, 300.0))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(600.0, 300.0), box2d.MakeB2Vec2(600.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(600.0, 0.0), box2d.MakeB2Vec2(0.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(0.0, 300.0))
		ground.CreateFixture(&shape, 0.0)
	}

	// setup colors
	//gray := color.Gray{Y: 0x99}
	red  := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	//blue := color.NRGBA{R: 0x33, G: 0x33, B: 0xff, A: 0xff}
	//yellow := color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	//green := color.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}

	// Circle character
	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(30.0, 50.0)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.FixedRotation = true
		bd.AllowSleep = false

		body := world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = 10

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 20.0
		fd.Restitution = 0.97
		body.CreateFixtureFromDef(&fd)

		c.circle = NewCircle(10,30.0,50.0,body)
		c.circle.CanMove(fyne.KeyDown,fyne.KeyUp,fyne.KeyLeft,fyne.KeyRight)
		c.game.Add(c.circle.circle)
		c.circle.Color(red,0,red)
		c.circle.Refresh()
		c.circle.pad = 3000000.0
	}

	// Circles character
	if (c.circles == nil) {
		c.circles = make([]*Circle,0)
	}
	for i:=0;i<10;i++ {
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(50.0+20.0*float64(i), 50.0)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.FixedRotation = true
		bd.AllowSleep = false

		body := world.CreateBody(&bd)

		shape := box2d.MakeB2CircleShape()
		shape.M_radius = 10

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 20.0
		fd.Restitution = 0.97
		body.CreateFixtureFromDef(&fd)

		circle := NewCircle(10,130.0+20*float64(i),50.0,body)
		c.game.Add(circle.circle)
		circle.Refresh()
		circle.pad = 3000000.0

		c.circles = append(c.circles,circle)
	}


	// Prepare for simulation. Typically we use a time step of 1/60 of a
	// second (60Hz) and 10 iterations. This provides a high quality simulation
	// in most game scenarios.
	timeStep := 1.0 / 200.0
	velocityIterations := 8
	positionIterations := 3

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
                				world.Step(timeStep, velocityIterations, positionIterations)

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
