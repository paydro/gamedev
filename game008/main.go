// Game 008
// * Add collision detection with the edges of the screen
// * Refactor window global vars into Window object
// * Change quit keyboard key to GUI key + q (COMMAND+q for Macs)
// * Change movement - protagonist moves automatically. Keys will change direction

package main

import (
	"errors"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"log"
	"os"
	"runtime"
)

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

type Window struct {
	Title    string
	Width    int
	Height   int
	FPS      int
	window   *sdl.Window
	renderer *sdl.Renderer
}

func NewWindow(title string, width, height, fps int) (*Window, error) {
	w := Window{
		Title:  title,
		Width:  width,
		Height: height,
		FPS:    fps,
	}

	ret := sdl.Init(sdl.INIT_EVERYTHING)
	if ret < 0 {
		return nil, errors.New(fmt.Sprintf("Failed to init SDL: %d\n", ret))
	}

	var window *sdl.Window
	window = sdl.CreateWindow(
		w.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		w.Width,
		w.Height,
		sdl.WINDOW_SHOWN)

	if window == nil {
		return nil, errors.New(fmt.Sprintf("Failed to create window: %s", sdl.GetError()))
	}
	w.window = window

	var renderer *sdl.Renderer
	renderer = sdl.CreateRenderer(w.window, -1, 0)
	if renderer == nil {
		return nil, errors.New(fmt.Sprintf("Failed to create renderer: %s", sdl.GetError()))
	}
	w.renderer = renderer

	return &w, nil
}

func (w *Window) Cleanup() {
	if w.renderer != nil {
		w.renderer.Destroy()
	}

	if w.window != nil {
		w.window.Destroy()
	}

	sdl.Quit()
}

type Direction int

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

func (d Direction) String() string {
	var s string
	switch d {
	case UP:
		s = "UP"
	case RIGHT:
		s = "RIGHT"
	case DOWN:
		s = "DOWN"
	case LEFT:
		s = "LEFT"
	}
	return s
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

	Direction

	// texture drawing info
	Width, Height int

	// animation info
	MaxFrames    int
	AnimFPS      float64
	lastTick     int32
	currentFrame int
}

func NewProtagonist(t *sdl.Texture) *Protagonist {
	return &Protagonist{
		DestX: 100,
		DestY: 100,

		MoveSpeed: 200,
		Direction: RIGHT,

		MaxFrames: 8,
		AnimFPS:   16.0,
		Width:     64,
		Height:    64,
		Texture:   t,
	}
}

func (p *Protagonist) Update(dt int, w *Window) {
	// keys := sdl.GetKeyboardState()
	toMove := int32(p.MoveSpeed * dt / 1000)

	switch p.Direction {
	case UP:
		p.DestY -= toMove
	case DOWN:
		p.DestY += toMove
	case RIGHT:
		p.DestX += toMove
	case LEFT:
		p.DestX -= toMove
	}

	// Screen collision -- probably put game info in a game world object

	gameWidth := int32(w.Width)
	gameHeight := int32(w.Height)
	width := int32(p.Width)
	height := int32(p.Height)

	if p.DestX < 0 {
		p.DestX = 0
	}
	if (p.DestX + width) > gameWidth {
		p.DestX = gameWidth - width
	}
	if p.DestY < 0 {
		p.DestY = 0
	}
	if (p.DestY + height) > gameHeight {
		p.DestY = gameHeight - height
	}

}

func (p *Protagonist) Draw(r *sdl.Renderer) {
	msPerFrame := 1.0 / (p.AnimFPS / 1000.0)

	if int(float64(p.lastTick)+msPerFrame) < int(sdl.GetTicks()) {
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

type ModifierKey uint16

func (m ModifierKey) String() string {
	var s string
	switch m {
	case sdl.KMOD_NONE:
		s = "KMOD_NONE"
	case sdl.KMOD_LSHIFT:
		s = "KMOD_LSHIFT"
	case sdl.KMOD_RSHIFT:
		s = "KMOD_RSHIFT"
	case sdl.KMOD_LCTRL:
		s = "KMOD_LCTRL"
	case sdl.KMOD_RCTRL:
		s = "KMOD_RCTRL"
	case sdl.KMOD_LALT:
		s = "KMOD_LALT"
	case sdl.KMOD_RALT:
		s = "KMOD_RALT"
	case sdl.KMOD_LGUI:
		s = "KMOD_LGUI"
	case sdl.KMOD_RGUI:
		s = "KMOD_RGUI"
	case sdl.KMOD_NUM:
		s = "KMOD_NUM"
	case sdl.KMOD_CAPS:
		s = "KMOD_CAPS"
	case sdl.KMOD_MODE:
		s = "KMOD_MODE"
	case sdl.KMOD_CTRL:
		s = "KMOD_CTRL"
	case sdl.KMOD_SHIFT:
		s = "KMOD_SHIFT"
	case sdl.KMOD_ALT:
		s = "KMOD_ALT"
	case sdl.KMOD_GUI:
		s = "KMOD_GUI"
	case sdl.KMOD_RESERVED:
		s = "KMOD_RESERVED"
	}

	return s
}

func main() {

	w, err := NewWindow("Game 008", 800, 600, 60)
	if err != nil {
		log.Fatalln("Could not create window.", err)
	}
	defer w.Cleanup()

	yoshiTexture, err := loadTexture("yoshi_trans_animation.png", w.renderer)
	defer yoshiTexture.Destroy()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load yoshi texture: %s", err)
		os.Exit(1)
	}

	yoshi := NewProtagonist(yoshiTexture)
	clock := NewClock(w.FPS)
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

			case *sdl.KeyDownEvent:
				// log.Printf("KeyDownEvent: %+v", t)
				// log.Printf(" * Scancode: %s", sdl.GetScancodeName(t.Keysym.Scancode))
				// log.Printf(" * Keycode: %s", sdl.GetKeyName(t.Keysym.Sym))
				// log.Printf(" * Modifier: %s", ModifierKey(t.Keysym.Mod))
				if t.Keysym.Scancode == sdl.SCANCODE_UP || t.Keysym.Scancode == sdl.SCANCODE_W {
					yoshi.Direction = UP
				} else if t.Keysym.Scancode == sdl.SCANCODE_DOWN || t.Keysym.Scancode == sdl.SCANCODE_S {
					yoshi.Direction = DOWN
				} else if t.Keysym.Scancode == sdl.SCANCODE_LEFT || t.Keysym.Scancode == sdl.SCANCODE_A {
					yoshi.Direction = LEFT
				} else if t.Keysym.Scancode == sdl.SCANCODE_RIGHT || t.Keysym.Scancode == sdl.SCANCODE_D {
					yoshi.Direction = RIGHT
				}

			case *sdl.KeyUpEvent:
				// log.Printf("KeyUpEvent: %+v\n", t)
				// log.Printf(" * Scancode: %s", sdl.GetScancodeName(t.Keysym.Scancode))
				// log.Printf(" * Keycode: %s", sdl.GetKeyName(t.Keysym.Sym))
				// log.Printf(" * Modifier: %s", ModifierKey(t.Keysym.Mod))

				if t.Keysym.Sym == sdl.K_q && (t.Keysym.Mod == sdl.KMOD_LGUI || t.Keysym.Mod == sdl.KMOD_RGUI) {
					log.Println("Quitting ...")
					running = false
				}
			}
		}

		// log.Println("Yoshi Direction:", yoshi.Direction)
		// Update entities
		yoshi.Update(dt, w)

		// Render
		w.renderer.SetDrawColor(205, 205, 205, 255)
		w.renderer.Clear()
		yoshi.Draw(w.renderer)
		w.renderer.Present()

		// log.Println("----------------------")
	}
}
