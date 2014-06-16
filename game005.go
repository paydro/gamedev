// Goal of this game is to attach keyboard movement to an image
// allowing the user to move the image around on screen.
//
// Thing to note here is that events only happen when a key is pressed or
// unpressed. That means if you attempt to catch keyboard events in the event
// loop, you'll need to keep track of the state of the key. OR you can just use
// SDL_GetKeyboardState().
//
// This also includes the type Clock. This helps maintain a maximum FPS for
// the application. It's very simple and behaves similarily to pygame's clock().
//
// Finally, after a long session of figuring out why the game keeps crashing,
// the SDL calls *must* happen in the main OS thread. To force this to happen,
// we create an init() function that calls runtime.LockOSThread(). This will
// keep out the segfault errors.
//
// TODO The veandco/go-sdl2/sdl* libraries are extremely raw bindings to SDL2.
// We could re-write the functions so that they use go idioms like returninng
// a second value for error.

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
	"runtime"
	"time"
)


var windowTitle = "Game 006 - Animation"
var screenWidth int = 640
var screenHeight int = 480

var dudeX, dudeY int = 0, 0
var dudeSpeed int = 200
var FPS int = 60

func init() {
	// SDL2 render commands are supposed to run on the main thread. This keeps
	// main to run on the same OS thread.
	// See:
	// * https://groups.google.com/forum/#!topic/golang-nuts/2_L7sPzC_6E
	// * https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()
}

// Clock
type Clock struct {
	LastTick uint32
	FPS      float32
}

func NewClock(fps int) *Clock {
	return &Clock{FPS: float32(fps)}
}

func (c *Clock) tick() int {
	var delay uint32
	msPerFrame := 1.0 / c.FPS * 1000.0
	fmt.Println("msPerFrame", msPerFrame)

	currentTick := sdl.GetTicks()
	fmt.Println("currentTick", currentTick)
	fmt.Println("LastTick", c.LastTick)

	if c.LastTick > 0 {
		processedIn := currentTick - c.LastTick
		fmt.Println("processedIn", processedIn)
		fmt.Println("uint32(msPerFrame)", uint32(msPerFrame))

		delay = 0
		if uint32(msPerFrame) > processedIn {
			delay = uint32(msPerFrame) - processedIn
			fmt.Println("Sleeping for", delay)
			sdl.Delay(uint32(delay))
		}
	}

	c.LastTick = sdl.GetTicks()
	return int(delay)
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

	dude, err := loadTexture("dude.png", renderer)
	defer dude.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load dude texture: %s", err)
		os.Exit(1)
	}

	clock := NewClock(FPS)
	var dt int
	var t0, t1 time.Time
	var event sdl.Event
	var running bool = true

	for running {
		dt = clock.tick()
		t0 = time.Now()

		// Handle events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("[%d ms] QuitEvent\n", t.Timestamp)
				running = false
			}
		}
		t1 = time.Now()
		fmt.Printf(" * Processed events in %v.\n", t1.Sub(t0))

		// Figure out what keys were pressed
		keys := sdl.GetKeyboardState()
		toMove := dudeSpeed * dt / 1000
		if keys[sdl.SCANCODE_UP] == 1 || keys[sdl.SCANCODE_W] == 1 {
			dudeY -= toMove
		}
		if keys[sdl.SCANCODE_DOWN] == 1 || keys[sdl.SCANCODE_S] == 1 {
			dudeY += toMove
		}
		if keys[sdl.SCANCODE_LEFT] == 1 || keys[sdl.SCANCODE_A] == 1 {
			dudeX -= toMove
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 || keys[sdl.SCANCODE_D] == 1 {
			dudeX += toMove
		}

		t1 = time.Now()
		fmt.Printf(" * Processed key presses in %v.\n", t1.Sub(t0))

		// Render
		tt := time.Now()

		renderer.SetDrawColor(205, 205, 205, 255)
		renderer.Clear()

		renderOriginalTexture(dude, renderer, dudeX, dudeY)
		renderer.Present()

		ttt := time.Now()
		fmt.Printf(" * Rendered in %v to run.\n", ttt.Sub(tt))

		t1 = time.Now()
		fmt.Printf("Frame took %v to run.\n", t1.Sub(t0))

		fmt.Printf("-------------\n\n")
	}

	window.Destroy()
}
