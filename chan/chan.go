package main

import "fmt"

func main() {
	ch1 := make(chan int)

	go func() { ch1 <- 1 }()

	s, ok := <-ch1
	_, ok2 := <-ch1

	fmt.Println(s, ok)
	fmt.Println(ok2)

}
