package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mileusna/crontab"
)

/***
 * Application
 */

const (
	RUNALL_EVENT          string = "/run/all"
	CIRCULATION_ON_EVENT  string = "/pump/circulation/on"
	CIRCULATION_OFF_EVENT string = "/pump/circulation/off"
	STOP_EVENT            string = "/stop/all"
	WARM_EVENT            string = "/warm/31"
	COOL_EVENT            string = "/cool"
	PUMP1_ON              string = "/pump/1/on"
	PUMP1_OFF             string = "/pump/1/off"
	PUMP2_ON              string = "/pump/2/on"
	PUMP2_OFF             string = "/pump/2/off"
)

func setupLog() {
	f, err := os.OpenFile("spa.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}

// run main loop
func mainLoop(event <-chan string) {
	for {
		switch <-event {
		case RUNALL_EVENT:
			go RunAllPump()
		case STOP_EVENT:
			go StopAllPump()
		case WARM_EVENT:
			go WarmUp()
		case COOL_EVENT:
			go CoolDown()
		case PUMP1_ON:
			go Pump(1, true)
		case PUMP1_OFF:
			go Pump(1, false)
		case PUMP2_ON:
			go Pump(2, true)
		case PUMP2_OFF:
			go Pump(2, false)
		case CIRCULATION_ON_EVENT:
			go Pump(0, true)
		case CIRCULATION_OFF_EVENT:
			go Pump(0, false)
		default:
			log.Println("Error: event is not implemented")
		}
		log.Println("wait 5s before new event")
		time.Sleep(time.Second * 5)
		log.Println("listening for new event")
	}
}

// SafeShutdown when application is stopped by the user
func SafeShutdown(stop <-chan os.Signal, s *http.Server) {
	// wait for stop
	<-stop

	log.Println("turn everything off")
	CoolDown()
	StopAllPump()

	log.Fatal("stop web server")
	err := s.Close()
	if err != nil {
		log.Fatal("fail to stop web server")
	}
}

// main application
func main() {
	// Setup log
	//setupLog()
	log.Println("start application")

	// Load application
	BootStrap()

	// run crontab
	cron := crontab.New()

	// run circulation pump
	event := make(chan string)

	// Run all pump every day for 10 minutes around 4pm
	cron.AddJob("0 0 4 * *", RunAllPump)
	cron.AddJob("0 10 4 * *", StopAllPump)
	log.Println("cron: run pump 10min every day")

	// Check if hot tub does not overheat every 5 minutes
	cron.AddJob("*/2 * * * *", CheckHeat)
	log.Println("cron: check for overheat every 5minutes")

	// run main loop
	go mainLoop(event)

	fmt.Println("start web server.\n http://spa/")
	h := http.NewServeMux()

	// temperature terading
	h.HandleFunc("/thermometer/read", readTemperatureHandler)
	//h.HandleFunc("/thermometer/history", recordHandler)

	// register action
	h.HandleFunc(RUNALL_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- RUNALL_EVENT
		ack(w)
	})

	h.HandleFunc(CIRCULATION_OFF_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- CIRCULATION_OFF_EVENT
		ack(w)
	})

	h.HandleFunc(CIRCULATION_ON_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- CIRCULATION_ON_EVENT
		ack(w)
	})

	h.HandleFunc(STOP_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- STOP_EVENT
		ack(w)
	})

	h.HandleFunc(WARM_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- WARM_EVENT
		ack(w)
	})

	h.HandleFunc(COOL_EVENT, func(w http.ResponseWriter, r *http.Request) {
		event <- COOL_EVENT
		ack(w)
	})

	h.HandleFunc(PUMP1_ON, func(w http.ResponseWriter, r *http.Request) {
		event <- PUMP1_ON
		ack(w)
	})

	h.HandleFunc(PUMP1_OFF, func(w http.ResponseWriter, r *http.Request) {
		event <- PUMP1_OFF
		ack(w)
	})

	h.HandleFunc(PUMP2_OFF, func(w http.ResponseWriter, r *http.Request) {
		event <- PUMP2_OFF
		ack(w)
	})

	h.HandleFunc(PUMP2_ON, func(w http.ResponseWriter, r *http.Request) {
		event <- PUMP2_ON
		ack(w)
	})

	// Show log
	h.HandleFunc("/log", logReaderHandler)

	// serve static file
	fs := http.FileServer(http.Dir("static"))
	h.Handle("/static/", http.StripPrefix("/static/", fs))

	// setup main handler
	h.HandleFunc("/", indexHandler)
	h.HandleFunc("/state", stateHandler)

	server := &http.Server{Addr: ":80", Handler: h}

	// register interrupt handler
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go SafeShutdown(stop, server)

	server.ListenAndServe()

}
