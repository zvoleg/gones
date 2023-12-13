package ppu

type color struct {
	r, g, b, a byte
}

func newColor(r, g, b byte) color {
	return color{
		r: r,
		g: g,
		b: b,
		a: 0xFF,
	}
}

type image struct {
	buff   []byte
	width  int
	height int
}

func newImage(width, height int) image {
	buff := make([]byte, width*height*4)
	return image{
		buff:   buff,
		width:  width,
		height: height,
	}
}
