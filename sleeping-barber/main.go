package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// adding some variables

var (
	seatingCapacity = 10
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	openTime        = 10 * time.Second
)

func main() {
	color.Yellow("Sleeping Barber Problem")
	color.Yellow("------------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// creating  Barber shop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		NumberOfBarbers: 0,
		CutDuration:     cutDuration,
		ClientsChan:     clientChan,
		barberDoneChan:  doneChan,
		isOpen:          true,
	}
	color.Green("Now, The barber shop is open")

	// adding barber to the shop
	shop.addBarber("Ahmed")
	shop.addBarber("Omar")

	shopClosingChan := make(chan bool)
	closed := make(chan bool)
	go func() {
		<-time.After(openTime)
		shopClosingChan <- true
		shop.closeShop()
		closed <- true
	}()

	// add clients to the shop
	clientNumber := 1
	go func() {
		for {
			// random time for client to arrive
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosingChan:
				return
			case <-time.After(15 * time.Millisecond * time.Duration(randomMilliseconds)):
				shop.addClient(fmt.Sprintf("Client #%d", clientNumber))
				clientNumber++
			}
		}
	}()

	<-closed
}
