//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"time"
	"image/color"
//	"fmt"
	"sort"

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
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(c.pad,0),true)
	} else if key.Name == c.rightKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(-c.pad,0),true)
	} else if key.Name == c.upKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(0,c.pad),true)
	} else if key.Name == c.downKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(0,-c.pad),true)
	}
}

func (c *Circle) onKeyDown(key *fyne.KeyEvent) {
	if key.Name == c.leftKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(-c.pad,0),true)
	} else if key.Name == c.rightKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(c.pad,0),true)
	} else if key.Name == c.upKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(0,-c.pad),true)
	} else if key.Name == c.downKey {
		c.body.ApplyForceToCenter(box2d.MakeB2Vec2(0,c.pad),true)
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

	output  *widget.RichText
	window  fyne.Window
	game    *fyne.Container
	world   box2d.B2World
}


func (c *box) onTypedKey(ev *fyne.KeyEvent) {
	c.circle.onTypedKey(ev)
}

func (c * box) Refresh() {
	if (c.circle!=nil) {
		c.circle.Refresh()
	}
}

func (c *box) onKeyUp(key *fyne.KeyEvent) {
	c.circle.onKeyUp(key)
}

func (c *box) onKeyDown(key *fyne.KeyEvent) {
	c.circle.onKeyDown(key)
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

	characters := make(map[string]*box2d.B2Body)

	// Ground body
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, 300.0), box2d.MakeB2Vec2(600.0, 300.0))
		ground.CreateFixture(&shape, 0.0)
		characters["00_ground"] = ground
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(600.0, 300.0), box2d.MakeB2Vec2(600.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
		characters["01_ground"] = ground
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(600.0, 0.0), box2d.MakeB2Vec2(0.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
		characters["02_ground"] = ground
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(0.0, 300.0))
		ground.CreateFixture(&shape, 0.0)
		characters["03_ground"] = ground
	}

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
		body.CreateFixtureFromDef(&fd)

		c.circle = NewCircle(10,30.0,50.0,body)
		c.circle.CanMove(fyne.KeyDown,fyne.KeyUp,fyne.KeyLeft,fyne.KeyRight)
		c.game.Add(c.circle.circle)
		c.circle.Refresh()
		characters["01_circle"] = c.circle.body
		c.circle.pad = 3000000.0
	}

	// Prepare for simulation. Typically we use a time step of 1/60 of a
	// second (60Hz) and 10 iterations. This provides a high quality simulation
	// in most game scenarios.
	timeStep := 1.0 / 60.0
	velocityIterations := 8
	positionIterations := 3

	characterNames := make([]string, 0)
	for k, _ := range characters {
		characterNames = append(characterNames, k)
	}

	sort.Strings(characterNames)

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
