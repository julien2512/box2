//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"time"
	"image/color"
//	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/ByteArena/box2d"
)

func MakeCircleBody(world *box2d.B2World,x float64,y float64,r float64) (*box2d.B2Body) {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(x, y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true
	bd.AllowSleep = false
	
	body := world.CreateBody(&bd)
	
	shape := box2d.MakeB2CircleShape()
	shape.M_radius = r
	
	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &shape
	fd.Density = 20.0
	fd.Restitution = 0.97
	body.CreateFixtureFromDef(&fd)
	
	return body
}

type iPlayer interface {
	getPad() float64
	setPad(pad float64)

	Color(color color.Color, wb float32, wc color.Color)
	Refresh()
	onTypedKey(ev *fyne.KeyEvent)
	onKeyUp(key *fyne.KeyEvent)
	onKeyDown(key *fyne.KeyEvent)
	CanMove(downKey fyne.KeyName, upKey fyne.KeyName, leftKey fyne.KeyName, rightKey fyne.KeyName)
	Move(position fyne.Position)
	getObject() fyne.CanvasObject
	getBody() *box2d.B2Body
}

type Player struct {
	x float64
	y float64

	pad float64

	downKey fyne.KeyName
	upKey	fyne.KeyName
	leftKey	fyne.KeyName
	rightKey fyne.KeyName

	body *box2d.B2Body

	object fyne.CanvasObject
}

func (c *Player) getBody() *box2d.B2Body {
	return c.body
}

func (c *Player) getObject() fyne.CanvasObject {
	return c.object
}

func (c *Player) getPad() float64 {
	return c.pad
}

func (c *Player) setPad(pad float64) {
	c.pad = pad
}

func (c *Player) Refresh() {
	c.x = c.body.GetPosition().X
	c.y = c.body.GetPosition().Y

	c.object.Move(fyne.NewPos(float32(c.x),float32(c.y)))
}

func (c *Player) Move(position fyne.Position) {
	c.object.Move(position)
}


func (c *Player) onTypedKey(ev *fyne.KeyEvent) {

}

func (c *Player) onKeyUp(key *fyne.KeyEvent) {
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

func (c *Player) onKeyDown(key *fyne.KeyEvent) {
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

func (c *Player) CanMove(downKey fyne.KeyName, upKey fyne.KeyName, leftKey fyne.KeyName, rightKey fyne.KeyName) {
	c.downKey = downKey
	c.upKey = upKey
	c.leftKey = leftKey
	c.rightKey = rightKey
}

type Circle struct {
	Player
	x float64
	y float64
	r float64
}

func (c *Circle) Color(color color.Color, wb float32, wc color.Color) {
	var circle *canvas.Circle
	circle = (c.object).(*canvas.Circle)
	circle.FillColor = color
	circle.StrokeWidth = wb
	circle.StrokeColor = wc
}

func NewCircle(world *box2d.B2World,size float64,xb float64, yb float64) *Circle {
	circle := &Circle{}
	circle.pad = 1
	circle.r = size
	circle.x = xb
	circle.y = yb

	circle.object = canvas.NewCircle(color.Black)
	circle.object.Resize(fyne.NewSize(float32(2*size),float32(2*size)))

	body := MakeCircleBody(world,xb,yb,size)
	circle.body = body

	return circle
}

type box struct {
	pad float32

	score1 int
	score2 int

	circles []iPlayer

	output  *widget.RichText
	window  fyne.Window
	game    *fyne.Container
	world   box2d.B2World
}

func (c *box) addPlayer(player iPlayer) {

	c.game.Add(player.getObject())
	if (c.circles == nil) {
		c.circles = make([]iPlayer,0)
	}
	
	c.circles = append(c.circles,player)
}

func (c *box) onTypedKey(ev *fyne.KeyEvent) {
	for i:=0;i<len(c.circles);i++ { 
		c.circles[i].onTypedKey(ev)
	}
}

func (c *box) Refresh() {
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].Refresh()
	}
}

func (c *box) onKeyUp(key *fyne.KeyEvent) {
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].onKeyUp(key)
	}
}

func (c *box) onKeyDown(key *fyne.KeyEvent) {
	for i:=0;i<len(c.circles);i++ {
		c.circles[i].onKeyDown(key)
	}
}

type GameLayout struct {
	cx float64
	cy float64
}

func (l *GameLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(float32(l.cx), float32(l.cy))
}

func (l *GameLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {

}

func NewGameLayout(cx float64,cy float64) (*GameLayout) {
	return &GameLayout{cx : cx, cy : cy }
}

func (c *box) loadUI(app fyne.App) {

	// make new window
	var cx float64 = 600
	var cy float64 = 300
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
		shape.Set(box2d.MakeB2Vec2(0.0, cy), box2d.MakeB2Vec2(cx, cy))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(cx, cy), box2d.MakeB2Vec2(cx, 0.0))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(cx, 0.0), box2d.MakeB2Vec2(0.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
	}
	{
		bd := box2d.MakeB2BodyDef()
		ground := world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(0.0, cy))
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
		circle := NewCircle(&world,10,30.0,50.0)
		circle.CanMove(fyne.KeyDown,fyne.KeyUp,fyne.KeyLeft,fyne.KeyRight)
		circle.Color(red,0,red)
		circle.setPad(3000000.0)

		c.addPlayer(circle)
	}

	// Other orbs
	count:=10
	xb := cx/2
	yb := cy/2
	for i:=0;i<count;i++ {
		alpha := ((float64)(i))*2*math.Pi/(float64)(count)
		dxb := math.Cos(alpha) // length 1
		dyb := math.Sin(alpha) // length 1
		bx := xb+cx*dxb/10
		by := yb+cx*dyb/10

		circle := NewCircle(&world,10,bx,by)
		circle.setPad(3000000.0)

		c.addPlayer(circle)
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
