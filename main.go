package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 600
	fontSize  = 4
	scorePos1 = (winWidth / 2) - (winWidth * 0.15)
	scorePos2 = (winWidth / 2) + (winWidth * 0.15)
)

const (
	playing = iota
	paused
)

var pixels = make([]byte, winWidth*winHeight*4)
var state = paused

type color struct {
	r, g, b, a byte
}

type position struct {
	x, y int
}

func drawPixel(x, y int, c color) {
	index := (x + winWidth*y) * 4
	if index < len(pixels) && index > 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
		pixels[index+3] = c.a
	}
}

func clearScreen() {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, _ := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	defer texture.Destroy()

	p1 := paddle{position{25, 300}, 15, 60, 10, color{255, 255, 255, 1}, 0}
	bot := paddle{position{775, 300}, 10, 60, 10, color{255, 255, 255, 1}, 0}
	b := ball{position{400, 300}, 10, 5, 5, color{255, 255, 255, 1}}

	keyState := sdl.GetKeyboardState()

	for {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.MouseMotionEvent:
				if t.YRel < 0 {
					if p1.y-(p1.height/2) > 0 {
						p1.y += int(t.YRel)
					}
				} else if t.YRel > 0 {
					if p1.y+(p1.height/2) < winHeight {
						p1.y += int(t.YRel)
					}
				}
			case *sdl.QuitEvent:
				return
			}
		}

		// play/pause
		if keyState[sdl.SCANCODE_SPACE] != 0 {
			sdl.Delay(100)
			if state == playing {
				state = paused
			} else {
				state = playing
			}
		}

		if state == playing {
			p1.update(keyState)
			bot.autoPlay(&b)
			b.update(&p1, &bot)
		}

		// reset game
		if p1.score == 9 || bot.score == 9 {
			state = paused
			p1.score = 0
			bot.score = 0
			b.reset()
		}

		clearScreen()
		drawNumber(p1.score, position{scorePos1, 70})
		drawNumber(bot.score, position{scorePos2, 70})
		p1.draw()
		bot.draw()
		b.draw()
		texture.Update(nil, pixels, winWidth*4)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		// cap framerate to 100fps
		elapsedTime := uint32(time.Since(frameStart).Milliseconds())
		if elapsedTime < 10 {
			sdl.Delay(10 - elapsedTime)
		}
	}
}
