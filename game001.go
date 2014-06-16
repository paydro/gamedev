package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)


var winWidth int = 800
var winHeight int = 600

func main() {

	// ret := sdl.Init(sdl.INIT_EVERYTHING)
	// if ret < 0 {
		// fmt.Fprintf(os.Stderr, "Failed to init SDL: %d\n", ret)
		// os.Exit(1)
	// }


	var window *sdl.Window
	window = sdl.CreateWindow(
		"Game 001",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN)

	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", sdl.GetError())
		os.Exit(1)
	}


	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", sdl.GetError())
		os.Exit(1)
	}


	var image *sdl.Surface
	image = sdl.LoadBMP("hello.bmp")
	if image == nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s", sdl.GetError())
		os.Exit(1);
	}

	var tex *sdl.Texture
	tex = renderer.CreateTextureFromSurface(image)
	image.Free()
	if tex == nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture from surface: %s", sdl.GetError())
		os.Exit(1);
	}

	didClear := renderer.Clear()
	if didClear < 0 {
		fmt.Fprintf(os.Stderr, "Failed to clear renderer: %d", didClear)
		os.Exit(1);
	}

	renderer.Copy(tex, nil, nil)
	renderer.Present()


	sdl.Delay(3000)

	window.Destroy()

}
