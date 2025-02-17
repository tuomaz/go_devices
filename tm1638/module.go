// Copyright 2015-2020, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tm1638

import "periph.io/x/periph/conn/gpio"

// TM1638 module should be conncted in following gpio pins
const (
	PinSTB  = "STB"
	PinCLK  = "CLK"
	PinDATA = "DATA"
)

// Color is type for LED colors
type Color byte

// Available colors for leds
const (
	Off Color = iota
	Green
	Red
)

// Module represents tm1638 based module
type Module struct {
	data      gpio.PinIO
	clk, stb  gpio.PinOut
	intensity uint8
}

// Open opens a tm1638 Module
// gpio driver should have following pins:
//   * data : data pin
//   * clk : clock pin
//   * stb : strobe pin
func Open(data gpio.PinIO, clk, stb gpio.PinOut, intensity uint8) (*Module, error) {
	m := &Module{
		data:      data,
		clk:       clk,
		stb:       stb,
		intensity: intensity,
	}

	if err := m.init(); err != nil {
		return nil, err
	}

	return m, nil
}

// SetLed sets led in pos to given color
func (m *Module) SetLed(pos int, led Color) {
	m.sendData(byte(pos<<1)+1, byte(led))
}

// SetFND sets FND in pos to data.
// The bits in the data are displayed by mapping bellow
//     -- 0 --
//    |       |
//    5       1
//     -- 6 --
//    4       2
//    |       |
//     -- 3 --  .7
func (m *Module) SetFND(pos int, data byte) {
	m.sendData(byte(pos)<<1, data)
}

// SetChar sets FND in given position to given character
func (m *Module) SetChar(pos int, c rune, dot bool) {
	data, ok := font[c]
	if !ok {
		data = 0x00
	}
	if dot {
		data |= 0x80
	}
	m.sendData(byte(pos)<<1, data)
}

// SetString sets FND to given str
func (m *Module) SetString(str string) {
	i := 0
	dot := false
	for _, r := range str {
		if r == '.' {
			dot = true
			continue
		}
		m.SetChar(i, r, dot)
		i++
		if r != '.' {
			dot = false
		}
	}
}

// GetButtons read buttons
func (m *Module) GetButtons() (keys byte) {
	m.stb.Out(gpio.Low)
	m.send(0x042)

	// TODO: why it repeats 4 time?
	for i := 0; i < 4; i++ {
		keys |= (m.receive() << uint(i))
	}
	m.stb.Out(gpio.High)

	return
}

func (m *Module) init() error {
	m.stb.Out(gpio.High)
	m.clk.Out(gpio.High)

	m.sendCmd(0x40)

	activeDisplay := byte(0x08)
	m.sendCmd(0x80 | m.intensity | activeDisplay)

	m.sendCmd(0xC0)
	for i := 0; i < 16; i++ {
		m.sendCmd(0x00)
	}

	return nil
}

func (m *Module) sendData(addr, data byte) {
	m.sendCmd(0x44)
	m.stb.Out(gpio.Low)
	m.send(0xC0 | addr)
	m.send(data)
	m.stb.Out(gpio.High)
}

func (m *Module) sendCmd(cmd byte) {
	m.stb.Out(gpio.Low)
	m.send(cmd)
	m.stb.Out(gpio.High)
}

func (m *Module) send(data byte) {
	for i := 0; i < 8; i++ {
		m.clk.Out(gpio.Low)
		if data&1 == 0 {
			m.data.Out(gpio.Low)
		} else {
			m.data.Out(gpio.High)
		}
		data >>= 1
		m.clk.Out(gpio.High)
	}
}

func (m *Module) receive() (data byte) {
	m.data.In(gpio.PullUp, gpio.RisingEdge)

	for i := 0; i < 8; i++ {
		data >>= 1
		m.clk.Out(gpio.Low)
		if b := m.data.Read(); b == gpio.High {
			data |= 0x80
		}
		m.clk.Out(gpio.High)
	}

	m.data.Out(gpio.Low)

	return
}
