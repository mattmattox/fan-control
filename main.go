package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
	"github.com/yryz/ds18b20"
)

//Setting Fan speed to zero at start
var fanSpeed int = 0
var currentTemperature float64

//Setting up fan pin
const (
	PWM_PIN = 12
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "temperature: %f\nfanspeed: %d\n", currentTemperature, fanSpeed)
}

func startWebServer() {
	http.HandleFunc("/metrics", handler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {

	//Parsing flags
	upperTemperatureLimitPar := flag.Float64("upperTemperatureLimit", 50.0, "Upper temperature limit")
	lowerTemperatureLimitPar := flag.Float64("lowerTemperatureLimit", 32.0, "Lower temperature limit")

	upperTemperatureLimit := *upperTemperatureLimitPar
	lowerTemperatureLimit := *lowerTemperatureLimitPar

	flag.Parse()

	go startWebServer()

	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	//Setting up fan controls
	motor_pin_pwm := rpio.Pin(PWM_PIN)
	motor_pin_pwm.Mode(rpio.Pwm)
	motor_pin_pwm.Freq(1000)
	motor_pin_pwm.DutyCycle(0, 100)
	fmt.Println("Starting sensors...")
	time.Sleep(5 * time.Second)

	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)

	for {
		for _, sensor := range sensors {
			t, err := ds18b20.Temperature(sensor)
			if err == nil {
				fmt.Printf("Temperature: %f\n", t)
				currentTemperature = t
				if currentTemperature > upperTemperatureLimit {
					fmt.Print("Need to ramp up fan\n")
					if fanSpeed != 100 {
						fanSpeed++
					}
				} else if currentTemperature < lowerTemperatureLimit {
					fmt.Print("Need to ramp down fan\n")
					if fanSpeed != 0 {
						fanSpeed--
					}
				} else {
					fmt.Print("Steady state\n")
				}
				motor_pin_pwm.DutyCycle(uint32(fanSpeed), 100)
				fmt.Printf("fanspeed: %d\n", fanSpeed)
			}
		}
		time.Sleep(60 * time.Second)
	}
}
