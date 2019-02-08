package main

import (
	"fmt"
	"time"

	sdl "./pixels"
)

const (
	WIDTH = 1920
	HEIGHT = 1080
)

func main() {

	sdl.Init(WIDTH, HEIGHT)
	defer sdl.Shutdown()

	i := 0
	starttime := time.Now()

	defer func() {
		duration := time.Now().Sub(starttime)
		fmt.Printf("%d frames in %v\n", i, duration)
	}()

	for {

		i++

		for x := 0; x < WIDTH; x += 4 {

			for y := 0; y < HEIGHT; y += 4 {

				r := (x + y) % i
				g := (x * y) % i
				b := (x - y) % i

				sdl.Set(x, y, r, g, b)
			}
		}

		sdl.Present()

		//

		sdl.HandleEvents()

		if sdl.MustQuit() {
			return
		}

		if sdl.GetKeyDown("Escape") {
			return
		}
	}
}
