package main

import (
	"log"
	"github.com/stianeikeland/go-rpio"
)

var thermometer []Thermometer

var (
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	circulation_pump_gpio          = rpio.Pin(23)
	heater_gpio                    = rpio.Pin(22)
	pump1_gpio                     = rpio.Pin(27)
	pump2_gpio                     = rpio.Pin(17)
	circulation_pump_state bool    = false
	heater_state           bool    = false
	pump1_state            bool    = false
	pump2_state            bool    = false
	MaxTemperature         float64 = 31.0
	MinTemperature         float64 = 15.0
)

// RunAllPump for 20 minutes
func RunAllPump() {
	log.Println("Run all pump")
	Pump(0, true)
	Pump(1, true)
	Pump(2, true)
}

// StopAllPump to prevent
func StopAllPump() {
	log.Println("Stop all pump")
	CoolDown()
	Pump(0, false)
	Pump(1, false)
	Pump(2, false)
}

// WarmUp enable heater
func WarmUp() {
	log.Println("Warm up the heater and run circulation pump")
	if circulation_pump_state == false {
		Pump(0, true)
	}

	heater_gpio.High()
	heater_state = true
}

// CoolDown stop heater
func CoolDown() {
	log.Println("Cool down")
	heater_gpio.Low()
	heater_state = false
}

// CheckHeat on the background to avoid overheat
func CheckHeat() {
	thermometer = ReadDatapoint(GetRoot())
	log.Println("check temperature: t0 = ", thermometer[0].Value, ", t1 = ", thermometer[1].Value)

	// too warm
	if thermometer[0].Value > MaxTemperature || thermometer[1].Value > MaxTemperature {
		log.Println("t > t_max = ", MaxTemperature)
		CoolDown()
		return
	}

	// too cool
	if thermometer[0].Value < MinTemperature && thermometer[1].Value < MinTemperature {
		log.Println("t < t_min = ", MinTemperature)
		WarmUp()
		return
	}
}

// Pump on/off
func Pump(num int, on bool) {
	var pump rpio.Pin

	switch num {
	case 0:
		pump = circulation_pump_gpio
		circulation_pump_state = on
	case 1:
		pump = pump1_gpio
		pump1_state = on
	case 2:
		pump = pump2_gpio
		pump2_state = on
	default:
		log.Fatalf("no pump: %d", num)
	}

	if on {
		pump.High()
	} else {
		pump.Low()
	}
	log.Println("pump", num, pump)
}
