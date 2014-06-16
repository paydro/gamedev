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

	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s", sdl.GetError())
		os.Exit(1)
	}

	backgroundTexture, err := loadTexture("background.png", renderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load background texture: %s", err)
		os.Exit(1)
	}

	foregroundTexture, err := loadTexture("happy_face_transparent.png", renderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load foreground texture: %s", err)
		os.Exit(1)
	}

	didClear := renderer.Clear()
	if didClear < 0 {
		fmt.Fprintf(os.Stderr, "Failed to clear renderer: %d", didClear)
		os.Exit(1)
	}

	xTiles := screenWidth / tileSize
	yTiles := screenHeight / tileSize

	for i := 0; i < xTiles * yTiles; i++ {
		x := i % xTiles
		y := i / xTiles
		renderTexture(backgroundTexture, renderer, x * tileSize, y * tileSize, tileSize, tileSize)
	}


	middleWidth := screenWidth / 2
	middleHeight := screenHeight / 2

	fmt.Println("Center x/y for window: (%d, %d)", middleWidth, middleHeight)

	var (
		fw int
		fh int
	)
	sdl.QueryTexture(foregroundTexture, nil, nil, &fw, &fh)
	middleFWidth := fw / 2
	middleFHeight := fh / 2

	renderOriginalTexture(foregroundTexture, renderer, middleWidth-middleFWidth, middleHeight-middleFHeight)

	renderer.Present()

	sdl.Delay(2000)

	// Cleanup
	// backgroundTexture.Destroy()
	// foregroundTexture.Destroy()
	// renderer.Destroy()
	window.Destroy()

}
