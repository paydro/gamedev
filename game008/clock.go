package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

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

