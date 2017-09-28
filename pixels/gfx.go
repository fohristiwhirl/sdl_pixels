package pixels

func Clear(r, g, b uint8) {
	renderer.SetDrawColor(r, g, b, 255)
	renderer.Clear()
}

func Set(x, y int, r, g, b, a uint8) {
	if x >= 0 && x < logical_width && y >= 0 && y < logical_height {
		i := y * logical_width * 4 + x * 4
		pixels[i + 0] = byte(r)
		pixels[i + 1] = byte(g)
		pixels[i + 2] = byte(b)
		pixels[i + 3] = byte(a)
	}
}

func Present() {
	virtue.Update(nil, pixels, int(logical_width) * 4)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Copy(virtue, nil, nil)
	renderer.Present()
}
