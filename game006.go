// Goal: Sprite animation.
// Create an animated yoshi!

package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
	"runtime"
	// "time"
)

var windowTitle = "Game 006 - Animation"
var screenWidth int = 640
var screenHeight int = 480
var FPS int = 30

var maxFrames = 8  // 8 total frames
var yoshiSize = 64 // both width & height
var yoshiFps float64 = 10.0   // Change this to watch yoshi run!

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

	currentTick := sdl.GetTicks()

	if c.LastTick > 0 {
		processedIn := currentTick - c.LastTick
		delay = 0
		if uint32(msPerFrame) > processedIn {
			delay = uint32(msPerFrame) - processedIn
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

	yoshi, err := loadTexture("yoshi_trans_animation.png", renderer)
	defer yoshi.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load yoshi texture: %s", err)
		os.Exit(1)
	}

	clock := NewClock(FPS)
	var dt int
	var event sdl.Event
	var running bool = true

	var sourceX, sourceY int32
	var lastTime int
	var currentFrame int

	for running {
		dt = clock.tick()
		fmt.Println(dt)

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

		// Render

		// Gray bg
		renderer.SetDrawColor(205, 205, 205, 255)
		renderer.Clear()

		msPerFrame := 1.0 / (yoshiFps / 1000.0)

		if int(float64(lastTime) + msPerFrame) < int(sdl.GetTicks()) {
			fmt.Println("Incrementing frame")
			currentFrame += 1
			lastTime = int(sdl.GetTicks())
		}

		currentFrame = currentFrame % maxFrames

		sourceX = 0 // Image only has one column - X is always 0
		sourceY = int32(currentFrame * yoshiSize)

		// Rect for texture (source rect)
		sourceRect := sdl.Rect{
			X: sourceX,
			Y: sourceY,
			W: int32(yoshiSize),
			H: int32(yoshiSize),
		}

		// Rect for placement on screen (dest rect)
		targetRect := sdl.Rect{
			X: int32((screenWidth / 2) + (yoshiSize / 2)),
			Y: int32((screenHeight / 2) + (yoshiSize / 2)),
			W: int32(yoshiSize),
			H: int32(yoshiSize),
		}

		renderer.Copy(yoshi, &sourceRect, &targetRect)

		renderer.Present()

	}
}
