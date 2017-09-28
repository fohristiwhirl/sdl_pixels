package pixels

// These things are not thread-safe.
// Also, only the goroutine that called Init() should call these, I think.

// We seem to be assuming BGRA format...

func Clear(r, g, b int) {

	// The whole point of this package is for situations where we don't clear, but meh...

	r_val := byte(clamp(r, 0, 255))
	g_val := byte(clamp(g, 0, 255))
	b_val := byte(clamp(b, 0, 255))

	i := 0

	for i < len(pixels) {
		pixels[i] = b_val
		i++
		pixels[i] = g_val
		i++
		pixels[i] = r_val
		i++
		pixels[i] = 255
		i++
	}
}

func Set(x, y int, r, g, b int) {
	if inbounds(x, y) {
		i := y * logical_width * 4 + x * 4
		pixels[i + 3] = 255
		pixels[i + 2] = byte(clamp(r, 0, 255))
		pixels[i + 1] = byte(clamp(g, 0, 255))
		pixels[i + 0] = byte(clamp(b, 0, 255))
	}
}

func Add(x, y int, r, g, b uint8) {

	if inbounds(x, y) {

		i := y * logical_width * 4 + x * 4

		new_r := min(255, int(pixels[i + 2]) + int(r))
		new_g := min(255, int(pixels[i + 1]) + int(g))
		new_b := min(255, int(pixels[i + 0]) + int(b))

		pixels[i + 3] = 255
		pixels[i + 2] = byte(new_r)
		pixels[i + 1] = byte(new_g)
		pixels[i + 0] = byte(new_b)
	}
}

func Present() {
	texture.Update(nil, pixels, int(logical_width) * 4)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func inbounds(x, y int) bool {
	if x >= 0 && x < logical_width && y >= 0 && y < logical_height {
		return true
	} else {
		return false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func clamp(val, min, max int) int {
	if val < min { val = min }
	if val > max { val = max }
	return val
}
