//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"time"
	"image/color"
//	"fmt"
	"math"
	"strconv"

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

func MakeRectangleBody(world *box2d.B2World,x float64,y float64,hx float64,hy float64) (*box2d.B2Body) {
	bd := box2d.MakeB2BodyDef()
	bd.Position.Set(x,y)
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.FixedRotation = true
	bd.AllowSleep = false

	body := world.CreateBody(&bd)

	shape := box2d.MakeB2PolygonShape()
	shape.SetAsBox(hx,hy)
	
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
	setType(bodytype uint8)
	setSensor(sensor bool)
	getSensor() bool
	setUserData(data interface{})
	setKind(kind uint8)
	getKind() uint8

	setOnPlayerOn(method func (player iPlayer))
	setOnPlayerOut(method func (player iPlayer))
	getOnPlayerOn() func (iPlayer)
	getOnPlayerOut() func (iPlayer)
}

type Player struct {
	x float64
	y float64

	pad float64

	downKey fyne.KeyName
	upKey	fyne.KeyName
	leftKey	fyne.KeyName
	rightKey fyne.KeyName

	kind uint8

	body *box2d.B2Body

	object fyne.CanvasObject

	onPlayerOn func (pl iPlayer)
	onPlayerOut func (pl iPlayer)
}

func (c *Player) setOnPlayerOn(method func (player iPlayer)) {
	c.onPlayerOn = method
}

func (c *Player) setOnPlayerOut(method func (player iPlayer)) {
	c.onPlayerOut = method
}

func (c *Player) getOnPlayerOn() func (player iPlayer) {
	return c.onPlayerOn
}

func (c *Player) getOnPlayerOut() func (player iPlayer) {
	return c.onPlayerOut
}

func (c *Player) setSensor(sensor bool) {
	c.body.GetFixtureList().SetSensor(sensor)
}

func (c *Player) getSensor() bool {
	return c.body.GetFixtureList().IsSensor()
}

func (c *Player) setType(bodytype uint8) {
	c.body.SetType(bodytype)
}

func (c *Player) setUserData(data interface{}) {
	c.body.SetUserData(data)
}

func (c *Player) setKind(kind uint8) {
	c.kind = kind
}

func (c *Player) getKind() uint8 {
	return c.kind
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
	c.object.Refresh()
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

func (c *Circle) Refresh() {
	c.x = c.body.GetPosition().X-c.r
	c.y = c.body.GetPosition().Y-c.r

	c.object.Move(fyne.NewPos(float32(c.x),float32(c.y)))
	c.object.Refresh()
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

type Rectangle struct {
	Player
	x float64
	y float64
	hx float64
	hy float64
}

func (c *Rectangle) Refresh() {
	c.x = c.body.GetPosition().X-c.hx/2
	c.y = c.body.GetPosition().Y-c.hy/2

	c.object.Move(fyne.NewPos(float32(c.x),float32(c.y)))
	c.object.Refresh()
}

func (c *Rectangle) Color(color color.Color, wb float32, wc color.Color) {
	var rectangle *canvas.Rectangle
	rectangle = (c.object).(*canvas.Rectangle)
	rectangle.FillColor = color
	rectangle.StrokeWidth = wb
	rectangle.StrokeColor = wc
}

func NewRectangle(world *box2d.B2World,hx float64,hy float64,x float64, y float64) *Rectangle {
	rectangle := &Rectangle{}
	rectangle.pad = 1
	rectangle.hx = hx
	rectangle.hy = hy
	rectangle.x = x
	rectangle.y = y
	
	rectangle.object = canvas.NewRectangle(color.Black)
	rectangle.object.Resize(fyne.NewSize(float32(hx),float32(hy)))
	
	body := MakeRectangleBody(world,x,y,hx/2,hy/2)
	rectangle.body = body
	
	return rectangle
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

func (c *box) BeginContact(contact box2d.B2ContactInterface) {
	fix1 := contact.GetFixtureA()
	fix2 := contact.GetFixtureB()

	if (fix1.GetBody().GetUserData() == nil || fix2.GetBody().GetUserData() == nil) {
		return
	}

	player1 := (fix1.GetBody().GetUserData()).(iPlayer)
	player2 := (fix2.GetBody().GetUserData()).(iPlayer)

	if (player1.getSensor()) {
		if (player2.getKind()==BALL) {
			if (player1.getOnPlayerOn() != nil) {
				player1.getOnPlayerOn()(player2)
			}
		}
	}
}

func (c *box) EndContact(contact box2d.B2ContactInterface) {
	fix1 := contact.GetFixtureA()
	fix2 := contact.GetFixtureB()

	if (fix1.GetBody().GetUserData() == nil || fix2.GetBody().GetUserData() == nil) {
		return
	}

	player1 := (fix1.GetBody().GetUserData()).(iPlayer)
	player2 := (fix2.GetBody().GetUserData()).(iPlayer)

	if (player1.getSensor()) {
		if (player2.getKind()==BALL) {
			if (player1.getOnPlayerOut() != nil) {
				player1.getOnPlayerOut()(player2)
			}
		}
	}
}

func (c *box) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {

}

func (c *box) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {

}

const BALL = 1

func (c *box) loadUI(app fyne.App) {

	// make new window
	var cx float64 = 800
	var cy float64 = 600
	c.window = app.NewWindow("Box")

	// make score layout
	score1 := canvas.NewText("P1 Score : 0", color.Black)
	score2 := canvas.NewText("P2 Score : 0", color.Black)
	scoreLayout := container.New(layout.NewHBoxLayout(), score1, layout.NewSpacer(), score2)

	// make help layout
	help1 := canvas.NewText("P1 : Z Q S D", color.Black)
	help2 := canvas.NewText("P2 : up left down right", color.Black)
	helpLayout := container.New(layout.NewHBoxLayout(), help1, layout.NewSpacer(), help2)

	// Make game layout
	c.game = container.New(NewGameLayout(cx,cy))
	c.window.SetContent(container.New(layout.NewVBoxLayout(),scoreLayout,c.game,helpLayout))
	c.window.Content().Refresh()
	c.window.Show()
	c.window.SetFixedSize(true)

	// Define the gravity vector.
	gravity := box2d.MakeB2Vec2(0.0, 0.0)

	// Construct a world object, which will hold and simulate the rigid bodies.
	world := box2d.MakeB2World(gravity)
	world.SetContactListener(c)

	// setup colors
	//gray := color.Gray{Y: 0x99}
	red  := color.NRGBA{R: 0xff, G: 0x33, B: 0x33, A: 0xff}
	blue := color.NRGBA{R: 0x33, G: 0x33, B: 0xff, A: 0xff}
	yellow := color.NRGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	green := color.NRGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
	purple := color.NRGBA{R: 0x9d, G: 0x00, B: 0xff, A: 0xff}

	// Goal character Player 1
	{
		rectangle := NewRectangle(&world,30.0,cy-10,20.0,cy/2)
		rectangle.Color(yellow,0,yellow)
		rectangle.setSensor(true)
		rectangle.setUserData(rectangle)

		c.addPlayer(rectangle)

		rectangle.onPlayerOn = func(pl iPlayer){ 
			pl.Color(purple,0,purple)
			pl.Refresh()
			c.score2 = c.score2+1
			score2.Text = "P2 Score :"+strconv.Itoa(c.score2)
			score2.Refresh()
		}
		rectangle.onPlayerOut = func(pl iPlayer){ 
			pl.Color(color.Black,0,color.Black)
			pl.Refresh()
		}
	}

	// Goal character Player 2
	{
		rectangle := NewRectangle(&world,30.0,cy-10,cx-20.0,cy/2)
		rectangle.Color(yellow,0,yellow)
		rectangle.setSensor(true)
		rectangle.setUserData(rectangle)

		c.addPlayer(rectangle)

		rectangle.onPlayerOn = func(pl iPlayer){ 
			pl.Color(purple,0,purple)
			pl.Refresh()
			c.score1 = c.score1+1
			score1.Text = "P1 Score :"+strconv.Itoa(c.score1)
			score1.Refresh()
		}
		rectangle.onPlayerOut = func(pl iPlayer){ 
			pl.Color(color.Black,0,color.Black)
			pl.Refresh()
		}
	}

	// Ground character
	{
		rectangle := NewRectangle(&world,1,cy,0,cy/2)
		rectangle.Color(blue,0,blue)
		rectangle.setType(box2d.B2BodyType.B2_staticBody)
		c.addPlayer(rectangle)
	}
	{
		rectangle := NewRectangle(&world,1,cy,cx,cy/2)
		rectangle.Color(blue,0,blue)
		rectangle.setType(box2d.B2BodyType.B2_staticBody)
		c.addPlayer(rectangle)
	}
	{
		rectangle := NewRectangle(&world,cx,1,cx/2,0)
		rectangle.Color(blue,0,blue)
		rectangle.setType(box2d.B2BodyType.B2_staticBody)
		c.addPlayer(rectangle)
	}
	{
		rectangle := NewRectangle(&world,cx,1,cx/2,cy)
		rectangle.Color(blue,0,blue)
		rectangle.setType(box2d.B2BodyType.B2_staticBody)
		c.addPlayer(rectangle)
	}

	// Static Center
	{
		circle := NewCircle(&world,30,cx/2,cy/2)
		circle.Color(blue,0,blue)
		circle.setType(box2d.B2BodyType.B2_staticBody)
		c.addPlayer(circle)
	}

	// Circle character Player 1
	{
		circle := NewCircle(&world,10,30.0,cy/2)
		circle.CanMove(fyne.KeyS,fyne.KeyZ,fyne.KeyQ,fyne.KeyD)
		circle.Color(green,0,green)
		circle.setPad(3000000.0)

		c.addPlayer(circle)
	}

	// Rectangle character Player 1
	{
		rectangle := NewRectangle(&world,50.0,50.0,100,cy/2)
		rectangle.CanMove(fyne.KeyS,fyne.KeyZ,fyne.KeyQ,fyne.KeyD)
		rectangle.Color(green,0,green)
		rectangle.setPad(300000000.0)

		c.addPlayer(rectangle)
	}


	// Circle character Player 2
	{
		circle := NewCircle(&world,10,cx-30,cy/2)
		circle.CanMove(fyne.KeyDown,fyne.KeyUp,fyne.KeyLeft,fyne.KeyRight)
		circle.Color(red,0,red)
		circle.setPad(3000000.0)

		c.addPlayer(circle)
	}

	// Rectangle character Player 2
	{
		rectangle := NewRectangle(&world,50.0,50.0,cx-100,cy/2)
		rectangle.CanMove(fyne.KeyDown,fyne.KeyUp,fyne.KeyLeft,fyne.KeyRight)
		rectangle.Color(red,0,red)
		rectangle.setPad(300000000.0)

		c.addPlayer(rectangle)
	}

	// Other orbs
	count:=20
	xb := cx/2
	yb := cy/2
	for i:=0;i<count;i++ {
		alpha := ((float64)(i))*2*math.Pi/(float64)(count)
		dxb := math.Cos(alpha) // length 1
		dyb := math.Sin(alpha) // length 1
		bx := xb+cx*dxb/8
		by := yb+cx*dyb/8

		circle := NewCircle(&world,10,bx,by)
		circle.setPad(3000000.0)
		circle.setUserData(circle)
		circle.setKind(BALL)

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
