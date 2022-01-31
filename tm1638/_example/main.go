package main

import (
	"time"

	"github.com/tuomaz/go_devices/tm1638"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	dev, err := tm1638.Open(
		gpioreg.ByName("17"), // data
		gpioreg.ByName("27"), // clk
		gpioreg.ByName("22"), // stb
		5,
	)
	if err != nil {
		panic(err)
	}

	dev.SetString(time.Now().Format("2006.01.02"))
	for i := 0; i < 8; i++ {
		if i%2 == 0 {
			dev.SetLed(i, tm1638.Green)
		} else {
			dev.SetLed(i, tm1638.Red)
		}
	}
}
