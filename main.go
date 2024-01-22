package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numberOfBarbers    = 3
	numberOfSeats      = 5
	shopOperatingHours = 10 * time.Second
)

type Customer struct {
	id int
}

type Barber struct {
	id int
}

var (
	waitingRoom     = make(chan Customer, numberOfSeats)
	barberReady     = make(chan int, numberOfBarbers) // Buffer to prevent deadlock
	customerDone    = make(chan bool)
	wg              sync.WaitGroup
	customerCounter = 0
	shopClosed      = make(chan bool)
)

func (b Barber) barberWork() {
	for {
		select {
		case customer := <-waitingRoom:
			fmt.Printf("Barber %d is cutting hair of customer %d\n", b.id, customer.id)
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond) // Simulate haircut time
			fmt.Printf("Barber %d finished with customer %d\n", b.id, customer.id)
			customerDone <- true
		case <-shopClosed:
			fmt.Printf("Barber %d is going home\n", b.id)
			wg.Done()
			return
		}
	}
}

func (c Customer) customerVisit() {
	defer wg.Done()

	select {
	case waitingRoom <- c:
		fmt.Printf("Customer %d is waiting\n", c.id)
		<-customerDone
		fmt.Printf("Customer %d is leaving with a haircut\n", c.id)
	default:
		fmt.Printf("Customer %d left because the waiting room is full\n", c.id)
	}
}

func main() {
	startTime := time.Now()

	// Start barbers
	for i := 1; i <= numberOfBarbers; i++ {
		wg.Add(1)
		go Barber{id: i}.barberWork()
	}

	// Generate customers for a fixed duration
	go func() {
		for {
			if time.Now().Sub(startTime) > shopOperatingHours {
				fmt.Println("Barbershop is closing now.")
				close(shopClosed) // Signal barbers the shop is closed
				return
			}
			customerCounter++
			wg.Add(1)
			go Customer{id: customerCounter}.customerVisit()
			time.Sleep(time.Duration(rand.Intn(250)) * time.Millisecond) // Customers arrive at random intervals
		}
	}()

	wg.Wait() // Wait for all barbers to go home
	fmt.Println("Barbershop is closed, all barbers have gone home.")
}
