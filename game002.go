package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"os"
	"errors"
)

var winWidth int = 1200
var winHeight int = 1024

func loadTexture(filepath string, renderer *sdl.Renderer) (*sdl.Texture, error) {
	var texture *sdl.Texture
	var surface *sdl.Surface

	surface = img.Load(filepath)
	if surface == nil {
		return nil, errors.New("Could not load BMP")
	}

	texture = renderer.CreateTextureFromSurface(surface)
	surface.Free()
	if texture == nil {
		return nil, errors.New("Could not create texture")
	}

	return texture, nil
}

func renderTexture(t *sdl.Texture, r *sdl.Renderer, x, y int) {
	rect := sdl.Rect{X: int32(x), Y: int32(y)}
	var w, h int
	sdl.QueryTexture(t, nil, nil, &w, &h)
	rect.W = int32(w)
	rect.H = int32(h)
	r.Copy(t, nil, &rect) // NOTE: This can fail -- need to check for this error
}


func main() {

	ret := sdl.Init(sdl.INIT_EVERYTHING)
	if ret < 0 {
		fmt.Fprintf(os.Stderr, "Failed to init SDL: %d\n", ret)
		os.Exit(1)
	}

	var window *sdl.Window
	window = sdl.CreateWindow(
		"Game 002",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN)

	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s", sdl.GetError())
		os.Exit(1)
	}


	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s", sdl.GetError())
		os.Exit(1)
	}

	backgroundTexture, err := loadTexture("/Users/paydro/Pictures/me_backflip.jpg", renderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load background texture: %s", err)
		os.Exit(1);
	}

	foregroundTexture, err := loadTexture("happy_face.bmp", renderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load foreground texture: %s", err)
		os.Exit(1);
	}


	didClear := renderer.Clear()
	if didClear < 0 {
		fmt.Fprintf(os.Stderr, "Failed to clear renderer: %d", didClear)
		os.Exit(1);
	}

	var bw, bh int
	sdl.QueryTexture(backgroundTexture, nil, nil, &bw, &bh)
	renderTexture(backgroundTexture, renderer, 0, 0)
	renderTexture(backgroundTexture, renderer, bw, 0)

	middleWidth := winWidth / 2
	middleHeight := winHeight / 2

	fmt.Println("Center x/y for window: (%d, %d)", middleWidth, middleHeight)

	var (
		fw int
		fh int
	)
	sdl.QueryTexture(foregroundTexture, nil, nil, &fw, &fh)
	middleFWidth := fw / 2
	middleFHeight := fh / 2

	renderTexture(foregroundTexture, renderer, middleWidth - middleFWidth, middleHeight - middleFHeight)

	// renderer.Copy(texture, nil, nil)
	renderer.Present()


	sdl.Delay(3000)

	// Cleanup
	backgroundTexture.Destroy()
	foregroundTexture.Destroy()
	renderer.Destroy()
	window.Destroy()

}

