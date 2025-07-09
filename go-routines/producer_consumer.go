package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	errChan := make(chan error)
	p.quit <- errChan
	return <-errChan
}

func pizzaJob(pizzaJob *Producer) {
	var i = 0

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			select {
			case pizzaJob.data <- *currentPizza:
			case errChan := <-pizzaJob.quit:
				close(pizzaJob.data)
				close(errChan)
				return
			}
		}
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		rnd := rand.Intn(12)
		msg := ""
		success := false

		delay := rand.Intn(5) + 1
		fmt.Printf("Making pizza %d, pizzaNumber: %d delay\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd < 5 {
			msg = "Pizza failed to make"
			pizzasFailed++
		} else {
			msg = "Pizza made successfully"
			pizzasMade++
			success = true
		}
		total++

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func produceConsumer() {
	fmt.Printf("Welcome to the pizza factory!\n")

	pizzaFactory := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzaJob(pizzaFactory)

	for pizza := range pizzaFactory.data {
		if pizza.success {
			fmt.Printf("Pizza %d: %s\n", pizza.pizzaNumber, pizza.message)
		} else {
			fmt.Printf("Pizza %d: %s\n", pizza.pizzaNumber, pizza.message)
		}

		if pizza.pizzaNumber >= NumberOfPizzas {
			err := pizzaFactory.Close()
			if err != nil {
				fmt.Printf("Error closing factory: %v\n", err)
			}
			break
		}
	}

	fmt.Printf("Total pizzas made: %d, Failed: %d, Total: %d\n", pizzasMade, pizzasFailed, total)
	if pizzasFailed > 8 {
		fmt.Printf("Too many failed pizzas, shutting down the factory!\n")
	} else if pizzasFailed > 5 {
		fmt.Printf("Some pizzas failed, but we can continue.\n")
	} else if pizzasFailed > 2 {
		fmt.Printf("Most pizzas were made successfully, great job!\n")
	} else {
		fmt.Printf("All pizzas were made successfully, excellent work!\n")
	}
}
