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

type display struct {
	pixels		[]*pixel
	count		int			// How many of the pixels in the array are still live. (We arrange for them to be at the start.)
	next_count	int			// How many live pixels there will be next loop.
	index		int			// Our index for this iteration.
	zoom		float64
	centre		complex128
	offset		int
	divisor		int
}

func (self *display) init() {
	self.pixels = make([]*pixel, WIDTH * HEIGHT, WIDTH * HEIGHT)
	self.zoom = 500
	self.centre = complex(-0.5, 0)
	self.divisor = 1
}

func (self *display) clear() {
	var i int
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			p := new(pixel)
			p.x = x
			p.y = y
			p.iterator.z = 0
			p.iterator.c = self.locate(x, y)
			self.pixels[i] = p
			i++
		}
	}
	self.count = WIDTH * HEIGHT
	self.next_count = 0
	self.index = 0
}

func (self *display) draw(p *pixel) {
	if p.escape > 0 {
		sdl.Set(p.x, p.y, (p.escape / self.divisor + self.offset) / 4, p.escape / self.divisor + self.offset, 0)
	}
}

func (self *display) redraw() {
	for i := 0; i < len(self.pixels); i++ {
		pixel := self.pixels[i]
		self.draw(pixel)
	}
}

func (self *display) progress(i int) {

	// Use some local variables to avoid constant indirection...

	index := self.index
	count := self.count
	next_count := self.next_count
	pixels := self.pixels

	for n := 0; n < i; n++ {

		if index >= count {
			index = 0
			count = next_count
			next_count = 0
		}

		pixel := pixels[index]
		pixel.iterate(10)

		if pixel.escape == 0 {

			// We want the "live" (undecided) pixels to be at the start of the slice.
			// So swap this pixel into the first available slot.

			tmp := pixels[next_count]
			pixels[next_count] = pixel
			pixels[index] = tmp
			next_count++

		} else if pixel.escape > 0 {
			self.draw(pixel)
		}

		index++
	}

	self.index = index
	self.count = count
	self.next_count = next_count
}

func (self *display) zoom_click(x, y int, multiplier float64) {
	self.centre = self.locate(x, y)
	self.zoom *= multiplier
	self.clear()
}

func (self *display) locate(x, y int) complex128 {
	dx := float64(x - WIDTH / 2) / self.zoom
	dy := float64(y - HEIGHT / 2) / self.zoom
	return complex(real(self.centre) + dx, imag(self.centre) + dy)
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

func main() {

	sdl.Init(WIDTH, HEIGHT)
	defer sdl.Shutdown()

	var state display
	state.init()
	state.clear()

	ticker := time.NewTicker(50 * time.Millisecond)

	for {

		sdl.HandleEvents()

		if sdl.MustQuit() {
			return
		}

		if sdl.GetKeyDown("Escape") {
			return
		}

		if sdl.GetKeyDownClear("r") {
			state.redraw()
		}

		if sdl.GetKeyDownClear("Keypad -") {
			state.offset -= 5
			state.redraw()
		}

		if sdl.GetKeyDownClear("Keypad +") {
			state.offset += 5
			state.redraw()
		}

		if sdl.GetKeyDownClear("Keypad /") {
			state.divisor += 1
			state.redraw()
		}

		if sdl.GetKeyDownClear("Keypad *") {
			state.divisor -= 1
			if state.divisor < 1 {
				state.divisor = 1
			}
			state.redraw()
		}

		click := sdl.GetLastMouseClick()
		if click.OK {
			sdl.Clear(0, 0, 0)

			if click.Button == sdl.LEFT {
				state.zoom_click(click.X, click.Y, 2)
			} else {
				state.zoom_click(click.X, click.Y, 0.5)
			}
		}

		state.progress(1000)

		select {

		case <- ticker.C:

			sdl.Present()

		default:

			// nothing

		}
	}
}
