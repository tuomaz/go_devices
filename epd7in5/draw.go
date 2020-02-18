// Copyright 2019, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epd7in5

import (
	"fmt"
	"image"
	"image/color"
)

// DrawImage draws a image
func (d *Display) DrawImage(i image.Image) error {
	b, err := d.Image2Buffer(i)
	if err != nil {
		return err
	}
	return d.DrawBuffer(b)
}

func (d *Display) refresh() {
	d.sendCmd(0x12)
	// time.Sleep(100 * time.Millisecond)
	d.waitTillNotBusy()
}

// Clear fill display with given patten in byte (8 pixel)
func (d *Display) Clear() error {
	db := make([]byte, d.w*d.h/2)
	for i := 0; i < len(db); i++ {
		db[i] = 0x33
	}
	// for i := 0; i < d.w/4*d.h; i++ {
	// 	for j := 0; j < 4; j++ { // TODO:
	// 		d.sendData(0x33)
	// 	}
	// }
	d.sendCmd(0x10)
	for i := 0; i < len(db); i += 4096 {
		d.sendDatas(db[i : i+4096])
	}
	d.refresh()

	return nil
}

// DrawBuffer draws buffer to display
func (d *Display) DrawBuffer(b []byte) error {
	// log.Println("DrawBuffer start")
	db := make([]byte, 0)
	for i := 0; i < d.w/4*d.h; i++ {
		tmp1 := b[i]
		var tmp2 byte
		j := 0
		for j < 4 {
			if tmp1&0xC0 == 0xC0 {
				tmp2 = 0x03
			} else if tmp1&0xC0 == 0x00 {
				tmp2 = 0x00
			} else {
				tmp2 = 0x04
			}
			tmp2 = (tmp2 << 4) & 0xFF
			tmp1 = (tmp1 << 2) & 0xFF
			j++
			if tmp1&0xC0 == 0xC0 {
				tmp2 |= 0x03
			} else if tmp1&0xC0 == 0x00 {
				tmp2 |= 0x00
			} else {
				tmp2 |= 0x04
			}
			tmp1 = (tmp1 << 2) & 0xFF
			// d.sendData(tmp2)
			db = append(db, tmp2)
			j++
		}
	}
	// log.Println("making display buffer done db len =", len(db))

	d.sendCmd(0x10)
	for i := 0; i < len(db); i += 4096 {
		d.sendDatas(db[i : i+4096])
	}
	d.refresh()
	// log.Println("DrawBuffer end")
	return nil
}

// Image2Buffer returns monochrome image buffer from image.Image
func (d *Display) Image2Buffer(img image.Image) ([]byte, error) {
	// log.Println("Image2Buffer")
	b := make([]byte, d.w*d.h/4)

	imgW, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgW == d.w && imgH == d.h {
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				switch checkColor(img.At(x, y)) {
				case black:
					b[(x+y*d.w)/4] &= ^(0xC0 >> (x % 4 * 2))
				case gray:
					b[(x+y*d.w)/4] &= ^(0xC0 >> (x % 4 * 2))
					b[(x+y*d.w)/4] |= (0x40 >> (x % 4 * 2))
				case white:
					b[(x+y*d.w)/4] |= (0xC0 >> (x % 4 * 2))
				}
			}
		}
	} else if imgW == d.h && imgH == d.w {
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				nx := y           // 160
				ny := d.h - x - 1 // 383
				switch checkColor(img.At(x, y)) {
				case black:
					b[(nx+ny*d.w)/4] &= ^(0xC0 >> (y % 4 * 2))
				case gray:
					b[(nx+ny*d.w)/4] &= ^(0xC0 >> (y % 4 * 2))
					b[(nx+ny*d.w)/4] |= (0x40 >> (y % 4 * 2))
				case white:
					b[(nx+ny*d.w)/4] |= (0xC0 >> (y % 4 * 2))
				}
			}
		}
	} else {
		return nil, fmt.Errorf("image size should be %dx%d of %dx%d", d.w, d.h, d.h, d.w)
	}

	return b, nil
}

func checkColor(c color.Color) inkColor {
	g := color.GrayModel.Convert(c).(color.Gray)

	if g.Y < 64 {
		return black
	} else if g.Y < 192 {
		return gray
	}
	return white
}
