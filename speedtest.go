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

	i := 0
	starttime := time.Now()

	sdl.Init(WIDTH, HEIGHT)
	defer sdl.Shutdown()

	defer func() {
		duration := time.Now().Sub(starttime)
		fmt.Printf("%d frames in %v\n", i, duration)
	}()

	for {

		i++

		for x := 0; x < WIDTH; x += 1 {

			for y := 0; y < HEIGHT; y += 1 {

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
