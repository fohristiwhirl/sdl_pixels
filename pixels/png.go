package pixels

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func DumpPNG(filename string) error {

	outfile, err := os.Create(filename)
	if err != nil {
		if outfile != nil {
			outfile.Close()
		}
		return fmt.Errorf("Couldn't create output file '%s'", filename)
	}

	field := image.NewNRGBA(image.Rect(0, 0, logical_width, logical_height))

	// Conversion of BRGA format into RGBA...

	for i := 0; i < logical_width * logical_height * 4; i += 4 {
		field.Pix[i + 0] = pixels[i + 2]
		field.Pix[i + 1] = pixels[i + 1]
		field.Pix[i + 2] = pixels[i + 0]
		field.Pix[i + 3] = pixels[i + 3]
	}

	png.Encode(outfile, field)
	outfile.Close()
	return nil
}
