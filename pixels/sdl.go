package pixels

import (
	"runtime"
	"github.com/veandco/go-sdl2/sdl"
)

type key_map_query struct {
	response_chan	chan bool
	key				string
}

var logical_width int
var logical_height int

var pixels []byte
var texture *sdl.Texture
var renderer *sdl.Renderer
var window *sdl.Window

var keyboard = make(map[string]bool)
var key_map_query_chan = make(chan key_map_query)

var fn_chan = make(chan func())
var shutdown_chan = make(chan bool)

var must_quit = false

func MustQuit() bool {
	return must_quit
}

func Shutdown() {
	window.Destroy()
	sdl.Quit()
}

func GetKeyDown(key string) bool {
	return keyboard[key]
}

func Init(width, height int) {

	// The goroutine that interacts with SDL should be locked to a thread.
	// Otherwise, crashes are possible, apparently.

	runtime.LockOSThread()

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic("Init(): " + err.Error())
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")

	// We use fullscreen at whatever the current resolution is, then use SetLogicalSize()
	// so that we can pretend it's whatever we want it to be.

	var dm sdl.DisplayMode
	sdl.GetDesktopDisplayMode(0, &dm)

	window, err = sdl.CreateWindow("SDL Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int(dm.W), int(dm.H), sdl.WINDOW_FULLSCREEN)
	if err != nil {
		panic("Init(): " + err.Error())
	}

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic("Init(): " + err.Error())
	}

	renderer.SetLogicalSize(width, height)
	logical_width = width
	logical_height = height

	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic("Init(): " + err.Error())
	}

	pixels = make([]byte, width * height * 4)

	// sdl.ShowCursor(sdl.DISABLE)
}

func HandleEvents() {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {

		case *sdl.QuitEvent:
			must_quit = true

		case *sdl.KeyDownEvent:
			keyboard[sdl.GetKeyName(t.Keysym.Sym)] = true

		case *sdl.KeyUpEvent:
			keyboard[sdl.GetKeyName(t.Keysym.Sym)] = false
		}
	}
}
