package main

import (
	"math/rand"
	"time"
	"./pixels"
)

const (
	WIDTH = 640
	HEIGHT = 480
)

func main() {
	pixels.Init(WIDTH, HEIGHT)
	defer pixels.Shutdown()

	render_time := time.Now()

	for {
		pixels.HandleEvents()

		if pixels.MustQuit() {
			return
		}

		if pixels.GetKeyDown("Escape") {
			return
		}

		for {
			x := rand.Intn(WIDTH)
			y := rand.Intn(HEIGHT)
			pixels.Set(x, y, 255, 128, 0)

			if time.Now().Sub(render_time) > 10 * time.Millisecond {
				break
			}
		}

		pixels.Present()
		render_time = time.Now()
	}
}
