package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ack
func ack(w http.ResponseWriter) {
	m := map[string]interface{}{}
	m["state"] = "ok"
	fmt.Fprintf(w, toJSON(m))
}

// toJson convert object to json
func toJSON(m interface{}) string {
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal("fail encoding datapoint to JSON")
	}
	return string(data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, _ := template.ParseFiles("api.html")
	t.Execute(w, nil)
}

// setTemperature min and max
// usage:
//  /set?t_max=31
//  /set?t_minx=31
func setHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query["t_max"] != nil {
		max, err := strconv.ParseFloat(query["t_max"][0], 64)
		if err != nil {
			fail(w, "fail decoding t_max")
			return
		}
		if max < MinTemperature {
			fail(w, "t_max < t_min <=> "+query["t_max"][0]+" < "+strconv.FormatFloat(MinTemperature, 'E', -1, 64))
			return
		}
		MaxTemperature = max
	} else if query["t_min"] != nil {
		min, err := strconv.ParseFloat(query["t_min"][0], 64)
		if err != nil {
			fail(w, "fail decoding t_min")
			return
		}

		if min > MaxTemperature {
			fail(w, "t_min > t_max <=> "+query["t_min"][0]+" > "+strconv.FormatFloat(MaxTemperature, 'E', -1, 64))
			return
		}
		MinTemperature = min
	} else {
		fail(w, "expect t_max=30.0 or t_min=15.0 ")
		return
	}

	m := map[string]interface{}{}
	m["state"] = "ok"
	m["thermometer_max"] = MaxTemperature
	m["thermometer_min"] = MinTemperature
	fmt.Fprintf(w, toJSON(m))
}

func fail(w http.ResponseWriter, msg string) {
	m := map[string]interface{}{}
	m["state"] = "fail"
	m["error"] = msg
	fmt.Fprintf(w, toJSON(m))
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	data["state"] = "ok"
	data["circulation"] = circulation_pump_state
	data["heater"] = heater_state
	data["pump1"] = pump1_state
	data["pump2"] = pump2_state
	data["thermometer"] = ReadDatapoint(GetRoot())
	data["thermometer_max"] = MaxTemperature
	data["thermometer_min"] = MinTemperature
	fmt.Fprintf(w, toJSON(data))
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
