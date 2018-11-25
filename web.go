package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// ack
func ack(w http.ResponseWriter) {
	m := map[string]interface{}{}
	m["state"] = "ok"
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal("fail encoding datapoint to JSON")
	}
	fmt.Fprintf(w, string(data))
}

func readTemperatureHandler(w http.ResponseWriter, r *http.Request) {
	thermometer := ReadDatapoint(GetRoot())
	data, err := json.Marshal(thermometer)
	if err != nil {
		log.Fatal("fail encoding datapoint to JSON")
	}
	fmt.Fprintf(w, string(data))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Current value to display
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	data["circulation"] = circulation_pump_state
	data["heater"] = heater_state
	data["pump1"] = pump1_state
	data["pump2"] = pump2_state
	data["thermometer"] = ReadDatapoint(GetRoot())
	data["max_temperature"] = MaxTemperature

	d, err := json.Marshal(data)
	if err != nil {
		log.Fatal("fail encoding datapoint to JSON")
	}
	fmt.Fprintf(w, string(d))
}

// Read log file and convert to json
func logReaderHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("spa.log")
	if err != nil {
		log.Fatal("fail reading log file")
	}
	lines := strings.Split(string(content), "\n")
	data, err := json.Marshal(lines)
	if err != nil {
		log.Fatal("fail encoding datapoint to JSON")
	}
	fmt.Fprintf(w, string(data))
}
