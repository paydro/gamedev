// Goal refactor game006 and game005.
// Consolidate yoshi vars into it's own struct. No more globals for Yoshi.
// Add ability to move yoshi around.

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
	"runtime"
)

var windowTitle = "Game 007 - Animation"
var screenWidth int = 640
var screenHeight int = 480
var FPS int = 60

func init() {
	// SDL2 render commands are supposed to run on the main thread. This keeps
	// main to run on the same OS thread.
	// See:
	// * https://groups.google.com/forum/#!topic/golang-nuts/2_L7sPzC_6E
	// * https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()
}


func loadTexture(filepath string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	texture := img.LoadTexture(renderer, filepath)
	if texture == nil {
		return nil, sdl.GetError()
	}
	return texture, nil
}

func renderTexture(t *sdl.Texture, r *sdl.Renderer, x, y, w, h int) {
	rect := sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(w),
		H: int32(h),
	}
	r.Copy(t, nil, &rect) // NOTE: This can fail -- need to check for this error
}

func renderOriginalTexture(t *sdl.Texture, r *sdl.Renderer, x, y int) {
	// Get the original width / height of texture
	var w, h int
	sdl.QueryTexture(t, nil, nil, &w, &h)
	renderTexture(t, r, x, y, w, h)
}

// Bloated struct representing Yoshi, our main character.
// Lots of refactorings can happen here.
type Protagonist struct {

	// Preloaded texture
	Texture *sdl.Texture

	// Where to draw protagonist
	DestX, DestY int32

	// Movement speed for protagonist
	MoveSpeed int

	// texture drawing info
	Width, Height int

	// animation info
	MaxFrames int
	AnimFPS float64
	lastTick int32
	currentFrame int
}

func NewProtagonist(t *sdl.Texture) *Protagonist {
	return &Protagonist{
		DestX: 100,
		DestY: 100,

		MoveSpeed: 200,

		MaxFrames: 8,
		AnimFPS: 16.0,
		Width: 64,
		Height: 64,
		Texture: t,
	}
}

func (p *Protagonist) Update(dt int) {
	keys := sdl.GetKeyboardState()
	toMove := int32(p.MoveSpeed * dt / 1000)
	if keys[sdl.SCANCODE_UP] == 1 || keys[sdl.SCANCODE_W] == 1 {
		p.DestY -= toMove
	}
	if keys[sdl.SCANCODE_DOWN] == 1 || keys[sdl.SCANCODE_S] == 1 {
		p.DestY += toMove
	}
	if keys[sdl.SCANCODE_LEFT] == 1 || keys[sdl.SCANCODE_A] == 1 {
		p.DestX -= toMove
	}
	if keys[sdl.SCANCODE_RIGHT] == 1 || keys[sdl.SCANCODE_D] == 1 {
		p.DestX += toMove
	}
}

func (p *Protagonist) Draw(r *sdl.Renderer) {
	msPerFrame := 1.0 / (p.AnimFPS / 1000.0)

	if int(float64(p.lastTick) + msPerFrame) < int(sdl.GetTicks()) {
		p.currentFrame += 1
		p.lastTick = int32(sdl.GetTicks())
	}

	p.currentFrame = p.currentFrame % p.MaxFrames

	sourceRect := sdl.Rect{
		X: 0,
		Y: int32(p.currentFrame * p.Height),
		W: int32(p.Width),
		H: int32(p.Height),
	}

	// Rect for placement on screen (dest rect)
	targetRect := sdl.Rect{
		X: p.DestX,
		Y: p.DestY,
		W: int32(p.Width),
		H: int32(p.Height),
	}

	r.Copy(p.Texture, &sourceRect, &targetRect)
}


func main() {
	ret := sdl.Init(sdl.INIT_EVERYTHING)
	if ret < 0 {
		fmt.Fprintf(os.Stderr, "Failed to init SDL: %d\n", ret)
		os.Exit(1)
	}
	defer sdl.Quit()

	var window *sdl.Window
	window = sdl.CreateWindow(
		windowTitle,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		screenWidth,
		screenHeight,
		sdl.WINDOW_SHOWN)

	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s", sdl.GetError())
		os.Exit(1)
	}
	defer window.Destroy()

	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(window, -1, 0)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()


	yoshiTexture, err := loadTexture("yoshi_trans_animation.png", renderer)
	defer yoshiTexture.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load yoshi texture: %s", err)
		os.Exit(1)
	}

	yoshi := NewProtagonist(yoshiTexture)

	clock := NewClock(FPS)
	var dt int
	var event sdl.Event
	var running bool = true



	// This is a variable game loop -- drawing/frames depend on dt
	for running {
		dt = clock.tick()

		// Handle events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("[%d ms] QuitEvent\n", t.Timestamp)
				running = false
			case *sdl.KeyUpEvent:
				if t.Keysym.Sym == sdl.K_q {
					fmt.Println("Quitting ...")
					running = false
				}
			}
		}

		// Update entities
		yoshi.Update(dt)

		// Render
		renderer.SetDrawColor(205, 205, 205, 255)
		renderer.Clear()
		yoshi.Draw(renderer)
		renderer.Present()

	}
}
