package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	ShopCapacity    int
	CutDuration     time.Duration
	NumberOfBarbers int
	ClientsChan     chan string
	barberDoneChan  chan bool
	isOpen          bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++
	go func() {
		isSleeping := false
		color.Yellow("%s goes to the waiting room to check for clients...", barber)
		for {
			if len(shop.ClientsChan) == 0 && shop.isOpen {
				color.Yellow("There is no client here, so %s is going to take a nap.", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s goes to wake %s up.", client, barber)
					isSleeping = false
				}
				// cut client's hair.
				shop.cutHair(client, barber)
			} else {
				// shop is closed so send barber home.
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(client string, barber string) {
	color.Green("%s cutting %s's hair.", barber, client)
	time.Sleep(shop.CutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)

}

func (shop *BarberShop) sendBarberHome(barber string) {
	color.Cyan("%s is going to home.", barber)
	shop.barberDoneChan <- true
}

func (shop *BarberShop) addClient(client string) {
	color.Yellow("**** %s is arrives.", client)
	if shop.isOpen {
		select {
		case shop.ClientsChan <- client:
			color.Yellow("%s takes a seat in waiting room.", client)
		default:
			color.Red("the waiting room is full, so %s leaves.", client)
		}
	} else {
		color.Yellow("%s is leaved because the shop is already closed.", client)
	}
}

func (shop *BarberShop) closeShop() {
	color.Cyan("shop is closed.")
	close(shop.ClientsChan)
	shop.isOpen = false
	for i := 1; i <= shop.NumberOfBarbers; i++ {
		<-shop.barberDoneChan
	}

	close(shop.barberDoneChan)
	color.Green("---------------------------------------------------")
	color.Green("The shop is closed now, and everyone has gone home.")
}
