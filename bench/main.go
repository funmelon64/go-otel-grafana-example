package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const rpm = 100

func main() {
	go func() {
		for {
			go sendMockAddBooking()
			<-time.After(time.Minute / time.Duration(rpm))
		}
	}()

	go func() {
		for {
			go sendMockGetBooking()
			<-time.After(time.Minute / time.Duration(rpm))
		}
	}()

	forever := make(chan struct{})
	<-forever
}

func sendMockAddBooking() {
	// Эмуляция POST запроса на /bookings
	postReqBody := bytes.NewBuffer([]byte(`{"id":"1", "time":"2022-12-31T23:59:59Z"}`)) // Замените на реальные данные
	resp, err := http.Post("http://web-entry:8080/bookings", "application/json", postReqBody)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response to POST /bookings:", string(body))
}

func sendMockGetBooking() {
	// Эмуляция GET запроса на /bookings/{id}
	resp, err := http.Get("http://web-entry:8080/bookings/1") // Замените "1" на реальный ID
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response to GET /bookings/{id}:", string(body))
}
