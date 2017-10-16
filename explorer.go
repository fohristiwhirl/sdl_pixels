package main

import (
	"math/cmplx"
	"time"

	sdl "./pixels"
)

const (
	WIDTH = 1920
	HEIGHT = 1080
)

type iterator struct {
	z			complex128
	z_slow		complex128
	c			complex128
}

type pixel struct {
	iterator
	x			int
	y			int
	iteration	int
	escape		int			// Negative: converges. Zero: unknown. Positive: escapes.
}

func (self *pixel) iterate(n uint) {

	for i := uint(0); i < n; i++ {

		self.iteration++

		self.z = self.z * self.z + self.c

		if self.z == self.z_slow {
			self.escape = -self.iteration
			return
		}

		if self.iteration % 2 == 0 {
			self.z_slow = self.z_slow * self.z_slow + self.c
		}

		if self.z == self.z_slow {
			self.escape = -self.iteration
			return
		}

		if cmplx.Abs(self.z) > 2 {
			self.escape = self.iteration
			return
		}
	}
}

func locate(x, y int, centre complex128, zoom float64) complex128 {
	dx := float64(x - WIDTH / 2) / zoom
	dy := float64(y - HEIGHT / 2) / zoom
	return complex(real(centre) + dx, imag(centre) + dy)
}

func clear_pixel_list(list []*pixel, centre complex128, zoom float64) {
	var index int
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			p := new(pixel)
			p.x = x
			p.y = y
			p.iterator.z = 0
			p.iterator.c = locate(x, y, centre, zoom)
			list[index] = p
			index++
		}
	}
}

func main() {

	var centre complex128 = complex(-0.5, 0)
	var zoom float64 = 500

	sdl.Init(WIDTH, HEIGHT)
	defer sdl.Shutdown()

	var list []*pixel = make([]*pixel, WIDTH * HEIGHT, WIDTH * HEIGHT)
	clear_pixel_list(list, centre, zoom)

	ticker := time.NewTicker(50 * time.Millisecond)

	var index int
	var next_count int
	var count int = WIDTH * HEIGHT

	for {

		sdl.HandleEvents()

		if sdl.MustQuit() {
			return
		}

		if sdl.GetKeyDown("Escape") {
			return
		}

		click := sdl.GetLastMouseClick()
		if click.OK {
			sdl.Clear(0, 0, 0)
			centre = locate(click.X, click.Y, centre, zoom)

			if click.Button == sdl.LEFT {
				zoom *= 2
			} else {
				zoom /= 2
			}

			clear_pixel_list(list, centre, zoom)
			count = WIDTH * HEIGHT
		}

		for n := 0; n < 1000; n++ {

			if index >= count {
				index = 0
				count = next_count
				next_count = 0
			}

			pixel := list[index]

			if pixel.escape == 0 {

				pixel.iterate(10)

				if pixel.escape == 0 {

					// We continually overwrite the low indices with the pixels that we need to work
					// on next iteration, and keep a count of how many indices there are.

					list[next_count] = pixel
					next_count++

				} else if pixel.escape > 0 {
					sdl.Set(pixel.x, pixel.y, pixel.escape / 4, pixel.escape, 0)
				}
			}

			index++
		}

		select {

		case <- ticker.C:

			sdl.Present()

		default:

			// nothing

		}
	}
}
