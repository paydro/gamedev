// Implementation http://www.willusher.io/sdl2%20tutorials/2013/08/18/lesson-3-sdl-extension-libraries/
package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
)

var windowTitle = "Game 003"
var screenWidth int = 640
var screenHeight int = 480
var tileSize int = 40

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
	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()

	backgroundTexture, err := loadTexture("background.png", renderer)
	defer backgroundTexture.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load background texture: %s", err)
		os.Exit(1)
	}

	foregroundTexture, err := loadTexture("happy_face_transparent.png", renderer)
	defer foregroundTexture.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load foreground texture: %s", err)
		os.Exit(1)
	}

	var event sdl.Event
	var running bool = true
	for running {


		// Process events
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			case *sdl.MouseButtonEvent:
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
				running = false
			case *sdl.MouseWheelEvent:
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n", t.Timestamp, t.Type, t.Which, t.X, t.Y)
				running = false
			case *sdl.KeyUpEvent:
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n", t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				running = false
			}
		}

		// Update game mechanics if there were any

		// Render

		didClear := renderer.Clear()
		if didClear < 0 {
			fmt.Fprintf(os.Stderr, "Failed to clear renderer: %d", didClear)
			os.Exit(1)
		}

		xTiles := screenWidth / tileSize
		yTiles := screenHeight / tileSize

		for i := 0; i < xTiles*yTiles; i++ {
			x := i % xTiles
			y := i / xTiles
			renderTexture(backgroundTexture, renderer, x*tileSize, y*tileSize, tileSize, tileSize)
		}

		middleWidth := screenWidth / 2
		middleHeight := screenHeight / 2

		var (
			fw int
			fh int
		)
		sdl.QueryTexture(foregroundTexture, nil, nil, &fw, &fh)
		middleFWidth := fw / 2
		middleFHeight := fh / 2

		renderOriginalTexture(foregroundTexture, renderer, middleWidth-middleFWidth, middleHeight-middleFHeight)
		renderer.Present()

		fmt.Printf("-----------------------------\n")
	}


}
