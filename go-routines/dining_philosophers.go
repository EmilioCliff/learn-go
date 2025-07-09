package main

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the
// table is laid for them ; each philosopher has their own place at the table.
// Their only difficulty - besides those of philosophy - is that the dish
// served is a very difficult kind of spaghetti which has to be eaten with
// two forks. There are two forks next to each plate, so that presents no difficulty
// As a consequence, however this means that no two neighbours may be eacting at the same time
// since there are five philosophers and only five forks.

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	leftFork  int
	rightFork int
}

var philosophers = []Philosopher{
	{"Plato", 0, 1},
	{"Socrates", 4, 0},
	{"Aristotle", 1, 2},
	{"Pascal", 2, 3},
	{"Locke", 3, 4},
}

var hunger = 3
var eatTime = 2
var thinkTime = 3
var forks = map[int]*sync.Mutex{}

var orderMutex = &sync.Mutex{}
var orderFinished = []string{}

func diningPhilosophers() {
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("------------------------------------")
	fmt.Println("The table is empty.")

	var wg sync.WaitGroup
	wg.Add(len(philosophers))

	var seated sync.WaitGroup
	seated.Add(len(philosophers))

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for _, philosopher := range philosophers {
		go func(p Philosopher) {
			defer wg.Done()

			fmt.Printf("%s seated at the table\n", p.name)
			seated.Done()

			seated.Wait()

			for i := hunger; i > 0; i-- {
				if philosopher.leftFork > philosopher.rightFork {
					forks[p.rightFork].Lock()
					fmt.Printf("\t %s takes the right fork.\n", philosopher.name)
					forks[p.leftFork].Lock()
					fmt.Printf("\t %s takes the left fork.\n", philosopher.name)
				} else {
					forks[p.leftFork].Lock()
					fmt.Printf("\t %s takes the left fork.\n", philosopher.name)
					forks[p.rightFork].Lock()
					fmt.Printf("\t %s takes the right fork.\n", philosopher.name)
				}

				fmt.Printf("\t %s is eating.\n", philosopher.name)
				time.Sleep(time.Duration(eatTime) * time.Second)

				fmt.Printf("\t %s is thinking.\n", philosopher.name)
				time.Sleep(time.Duration(thinkTime) * time.Second)

				forks[p.leftFork].Unlock()
				forks[p.rightFork].Unlock()

				fmt.Printf("\t %s puts down the forks.\n", philosopher.name)
			}

			orderMutex.Lock()
			orderFinished = append(orderFinished, philosopher.name)
			orderMutex.Unlock()

			fmt.Println(philosopher.name, "has finished eating.")
			fmt.Println(philosopher.name, "is leaving the table.")
		}(philosopher)
	}

	wg.Wait()

	fmt.Println("The table is empty again.")
	fmt.Println("All philosophers have finished eating.")
	fmt.Printf("Order of finishing: %s\n", strings.Join(orderFinished, ", "))

	fmt.Println("------------------------------------")
	fmt.Println("Philosophers can now think about their next meal.")
}
