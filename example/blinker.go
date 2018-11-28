package main

import (
	"fmt"
	"os"
	"rpio"
	"time"
)

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	relay0 = rpio.Pin(23)
	relay1 = rpio.Pin(22)
	relay2 = rpio.Pin(27)
	relay3 = rpio.Pin(17)
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Set pin to output mode
	relay0.Output()
	relay1.Output()
	relay2.Output()
	relay3.Output()

	relay0.Low()
	relay1.Low()
	relay2.Low()
	relay3.Low()
	// os.Exit(0)

	// Toggle pin 20 times
	for x := 0; x < 100; x++ {
		relay0.Toggle()
		relay1.Toggle()
		relay2.Toggle()
		relay3.Toggle()
		time.Sleep(time.Second * 2)
	}
}
