package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status `json:"status"`
}

type Condition struct {
	Water string
	Wind  string
}

type Response struct {
	Status
	Condition
}

func updateData() {
	for {

		var data = Data{Status: Status{}}
		maxValue := 20

		data.Status.Water = rand.Intn(maxValue)
		data.Status.Wind = rand.Intn(maxValue)

		b, err := json.MarshalIndent(&data, "", " ")

		if err != nil {
			log.Fatalln("error while marshalling json data  =>", err.Error())
		}

		err = ioutil.WriteFile("data.json", b, 0644)

		if err != nil {
			log.Fatalln("error while writing value to data.json file  =>", err.Error())
		}
		fmt.Println("menggungu 5 detik")
		time.Sleep(time.Second * 5)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl, _ := template.ParseFiles("index.html")

		var data = Data{Status: Status{}}

		b, err := ioutil.ReadFile("data.json")

		if err != nil {
			fmt.Fprint(w, "error braderku")
			return
		}

		if err = json.Unmarshal(b, &data); err != nil {
			fmt.Fprint(w, "error braderku")
			return
		}

		response := Response{Status: data.Status}
		response.Condition.Water = generateWaterStatus(data.Status.Water)
		response.Condition.Wind = generateWindStatus(data.Status.Wind)

		tpl.ExecuteTemplate(w, "index.html", response)

	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}

func generateWaterStatus(status int) string {
	waterStatus := ""
	switch {
	case status <= 6:
		waterStatus = "Aman"
	case status > 6 && status <= 9:
		waterStatus = "Siaga"
	case status > 9:
		waterStatus = "Bahaya"
	}
	return waterStatus
}

func generateWindStatus(status int) string {
	windStatus := ""
	switch {
	case status <= 6:
		windStatus = "Aman"
	case status > 6 && status <= 15:
		windStatus = "Siaga"
	case status > 15:
		windStatus = "Bahaya"
	}
	return windStatus
}
