package main

import "github.com/veandco/go-sdl2/sdl"

type paddle struct {
	position
	velocity int
	height   int
	width    int
	color    color
	score    int
}

func (paddle *paddle) draw() {
	posX := paddle.x - paddle.width/2
	posY := paddle.y - paddle.height/2

	for y := 0; y < paddle.height; y++ {
		for x := 0; x < paddle.width; x++ {
			drawPixel(posX+x, posY+y, paddle.color)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		if paddle.y-(paddle.height/2) > 0 {
			paddle.y -= paddle.velocity
		}
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		if paddle.y+(paddle.height/2) < winHeight {
			paddle.y += paddle.velocity
		}
	}
}

func (paddle *paddle) autoPlay(ball *ball) {
	paddle.y = ball.y
}
