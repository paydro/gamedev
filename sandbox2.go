// Goal of this game is to attach keyboard movement to an image
// allowing the user to move the image around on screen.

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	// "github.com/veandco/go-sdl2/sdl_image"
	"os"
	"time"
)


// #cgo windows LDFLAGS: -lSDL2
// #cgo darwin LDFLAGS: -framework SDL2
// #cgo linux freebsd pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

var windowTitle = "Game 005"
var screenWidth int = 640
var screenHeight int = 480

var dudeX, dudeY int = 0, 0
var dudeSpeed int = 100
var FPS int = 60

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
const (
	INIT_TIMER          = 0x00000001
	INIT_AUDIO          = 0x00000010
	INIT_VIDEO          = 0x00000020
	INIT_JOYSTICK       = 0x00000200
	INIT_HAPTIC         = 0x00001000
	INIT_GAMECONTROLLER = 0x00002000
	INIT_NOPARACHUTE    = 0x00100000
	INIT_EVERYTHING     = INIT_TIMER | INIT_AUDIO | INIT_VIDEO | INIT_JOYSTICK |
		INIT_HAPTIC | INIT_GAMECONTROLLER
)
func main() {


func Init(flags uint32) int {
	return int(C.SDL_Init(C.Uint32(flags)))
}

	var flags uint32 = INIT_EVERYTHING
	C.SDL_Init(C.Uint32(flags))
	defer C.SDL_Quit()

	_title := C.CString(windowTitle)
	defer C.free(unsafe.Pointer(_title))
	var _window = C.SDL_CreateWindow(_title, C.SDL_WINDOWPOS_UNDEFINED, C.SDL_WINDOWPOS_UNDEFINED, C.int(screenWidth), C.int(screenHeight), C.SDL_WINDOW_SHOWN)
	deferC.SDL_DestroyWindow(window.cptr())
	// return (*Window)(unsafe.Pointer(_window))

	// var window *sdl.Window
	// window = sdl.CreateWindow(
		// windowTitle,
		// sdl.WINDOWPOS_UNDEFINED,
		// sdl.WINDOWPOS_UNDEFINED,
		// screenWidth,
		// screenHeight,
		// sdl.WINDOW_SHOWN)

	// if window == nil {
		// fmt.Fprintf(os.Stderr, "Failed to create window: %s", sdl.GetError())
		// os.Exit(1)
	// }
	// defer window.Destroy()

	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()

	// clock := NewClock(FPS)
	// var dt int
	var t0, t1 time.Time
	var event sdl.Event
	var running bool = true
	for running {

		// sdl.Delay(10)
		// _ = clock.tick()
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
		toMove := dudeSpeed / 1000
		if keys[sdl.SCANCODE_UP] == 1 {
			dudeY -= toMove
		}
		if keys[sdl.SCANCODE_DOWN] == 1 {
			dudeY += toMove
		}
		if keys[sdl.SCANCODE_LEFT] == 1 {
			dudeX -= toMove
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 {
			dudeX += toMove
		}

		t1 = time.Now()
		fmt.Printf(" * Processed key presses in %v.\n", t1.Sub(t0))

		// Render
		tt := time.Now()

		renderer.SetDrawColor(205, 205, 205, 255)
		renderer.Clear()
		renderer.Present()

		ttt := time.Now()
		fmt.Printf(" * Rendered in %v to run.\n", ttt.Sub(tt))

		t1 = time.Now()
		fmt.Printf("Frame took %v to run.\n", t1.Sub(t0))

		fmt.Printf("-------------\n")
	}


	window.Destroy()
}



