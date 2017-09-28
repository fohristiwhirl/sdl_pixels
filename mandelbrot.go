package main

import (
	"math/cmplx"
	"sync"
	"time"
	"./pixels"
)

const (
	WIDTH = 1920
	HEIGHT = 1080
	ZOOM = 400
	X_OFFSET = 1000
	MAX_ITERATIONS = 1000
	THREADS = 4
)

var mutex sync.Mutex
var threads_running int
var done_chan chan bool = make(chan bool)

func main() {
	pixels.Init(WIDTH, HEIGHT)
	defer pixels.Shutdown()

	render_time := time.Now()

	pixel_x := 0
	pixel_y := 0

	for {
		pixels.HandleEvents()

		if pixels.MustQuit() {
			return
		}

		if pixels.GetKeyDown("Escape") {
			return
		}

		for {

			if threads_running >= THREADS {
				<- done_chan						// Wait till a thread ends before starting a new one
				go iterator(pixel_x, pixel_y)
			} else {
				go iterator(pixel_x, pixel_y)
				threads_running += 1
			}

			pixel_x += 1
			if pixel_x >= WIDTH {
				pixel_x = 0
				pixel_y += 1
			}

			if time.Now().Sub(render_time) > 20 * time.Millisecond {
				break
			}
		}

		mutex.Lock()
		pixels.Present()
		mutex.Unlock()
		render_time = time.Now()
	}
}

func iterator(pixel_x, pixel_y int) {

	x := float64(pixel_x - WIDTH / 2) / ZOOM
	y := float64(pixel_y - HEIGHT / 2) / ZOOM

	var c complex128 = complex(x, y)
	var z complex128 = c

	for n := 0 ; n < MAX_ITERATIONS ; n++ {

		z = z * z + c

		if cmplx.Abs(z) > 2 {               // The particle does escape
			mutex.Lock()
			pixels.Set(pixel_x, pixel_y, n, n, 0)
			mutex.Unlock()
			break
		}
	}

	done_chan <- true
}
