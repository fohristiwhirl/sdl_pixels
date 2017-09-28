package main

import (
	"math"
	"math/cmplx"
	"math/rand"
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

	for {
		pixels.HandleEvents()

		if pixels.MustQuit() {
			return
		}

		if pixels.GetKeyDown("Escape") {
			return
		}

		for {
			var c complex128 = complex(rand.Float64() * 3.32 - 2.11, rand.Float64() * 2.48 - 1.24)

			if threads_running >= THREADS {
				<- done_chan					// Wait till a thread ends before starting a new one
				go iterator(c)
			} else {
				go iterator(c)
				threads_running += 1
			}

			if time.Now().Sub(render_time) > 10 * time.Millisecond {

				for threads_running > 0 {		// Wait till all threads have finished
					<- done_chan
					threads_running -= 1
				}

				break
			}
		}

		pixels.Present()
		render_time = time.Now()
	}
}

func iterator(c complex128) {

	var z complex128 = c
	var list [MAX_ITERATIONS]complex128

	for n := 0 ; n < MAX_ITERATIONS ; n++ {

		z = z * z + c
		list[n] = z

		if cmplx.Abs(z) > 2 {               // The particle does escape, so draw the list of points
			for i := 0 ; i <= n ; i++ {
				pixel_x := int(math.Floor(real(list[i]) * ZOOM)) + X_OFFSET         // It's OK if x,y is out of bounds
				pixel_y := int(math.Floor(imag(list[i]) * ZOOM)) + HEIGHT / 2
				mutex.Lock()
				pixels.Add(pixel_x, pixel_y, 1, 1, 1)
				pixels.Add(pixel_x, HEIGHT - pixel_y, 1, 1, 1)
				mutex.Unlock()
			}
			break
		}
	}

	done_chan <- true
}
