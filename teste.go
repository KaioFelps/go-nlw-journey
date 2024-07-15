package main

import (
	"fmt"
	"math/rand"
)

func sum(total chan int, exit chan bool) {
	soma := rand.Intn(100)
	soma += 1000

	for {
		select {
		case total <- soma:
			soma += 1000
		case <-exit:
			fmt.Println("finalizando o processo")
			return
		}
	}
}

func main() {
	somatoria := make(chan int)
	signal := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("soma: %d \n", <-somatoria)
		}

		signal <- true
	}()

	sum(somatoria, signal)
}
