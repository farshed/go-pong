package main

type ball struct {
	position
	radius int
	xv     int
	yv     int
	color  color
}

func (ball *ball) draw() {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				drawPixel(ball.x+x, ball.y+y, ball.color)
			}
		}
	}
}

func (ball *ball) update(left, right *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	// update score and reset
	if ball.x <= 0 {
		right.score++
		ball.reset()
		return
	} else if ball.x >= winWidth {
		left.score++
		ball.reset()
		return
	}

	// ground and roof collision
	if ball.y-ball.radius < 0 || ball.y+ball.radius > winHeight {
		ball.yv = -ball.yv
	}

	// paddle collision
	if ball.x-(ball.radius/2) <= left.x+(left.width/2) {
		if ball.y > left.y-(left.height/2) && ball.y < left.y+(left.height/2) {
			ball.xv = -ball.xv
		}
	} else if ball.x+(ball.radius/2) >= right.x-(right.width/2) {
		if ball.y > right.y-(right.height/2) && ball.y < right.y+(right.height/2) {
			ball.xv = -ball.xv
		}
	}
}

func (ball *ball) reset() {
	ball.x = winWidth / 2
	ball.y = winHeight / 2
	ball.xv = -ball.xv
}
