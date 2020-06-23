package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	playing = iota
	paused
)

var (
	winWidth     = 0
	winHeight    = 0
	centerX      = 0
	centerY      = 0
	scorePos1    = 0
	scorePos2    = 0
	fontSize     = 5
	paddleOffset = 50
	pixels       = []byte{}
	state        = paused
)

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

	winSize, _ := sdl.GetDisplayBounds(0)
	winWidth = int(winSize.W)
	winHeight = int(winSize.H)
	centerX = winWidth / 2
	centerY = winHeight / 2
	pixels = make([]byte, winWidth*winHeight*4)
	scorePos1 = centerX - (winWidth * 15 / 100)
	scorePos2 = centerX + (winWidth * 15 / 100)

	window, _ := sdl.CreateWindow("Pong", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_FULLSCREEN)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	texture, _ := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	defer texture.Destroy()

	p1 := paddle{position{paddleOffset, centerY}, 15, 60, 10, color{255, 255, 255, 1}, 0}
	bot := paddle{position{winWidth - paddleOffset, centerY}, 10, 60, 10, color{255, 255, 255, 1}, 0}
	b := ball{position{centerX, centerY}, 10, 5, 5, color{255, 255, 255, 1}}

	keyState := sdl.GetKeyboardState()

	for {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch d := event.(type) {
			case *sdl.MouseMotionEvent:
				if state == playing {
					if d.YRel < 0 {
						if p1.y-(p1.height/2) > 0 {
							p1.y += int(d.YRel)
						}
					} else if d.YRel > 0 {
						if p1.y+(p1.height/2) < winHeight {
							p1.y += int(d.YRel)
						}
					}
				}
			case *sdl.QuitEvent:
				return
			}
		}

		// quit game
		if keyState[sdl.SCANCODE_ESCAPE] != 0 {
			return
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
