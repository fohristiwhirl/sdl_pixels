package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WIDTH = 1920
	HEIGHT = 1080
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_FULLSCREEN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	i := 0
	starttime := time.Now()

	defer func() {
		duration := time.Now().Sub(starttime)
		fmt.Printf("%d frames in %v\n", i, duration)
	}()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				return
			}
		}

		i++

		for x := 0; x < WIDTH; x += 4 {
			for y := 0; y < HEIGHT; y += 4 {
				r := (x + y) % i
				g := (x * y) % i
				b := (x - y) % i

				if r > 255 { r = 255 }
				if g > 255 { g = 255 }
				if b > 255 { b = 255 }

				if b < 0 { b = 0 }

				surface.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}

		window.UpdateSurface()
	}
}
