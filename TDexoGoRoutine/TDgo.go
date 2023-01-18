package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Ex 1
func Producer(msg []string) chan string {
	ch := make(chan string)
	go func() {
		for _, m := range msg {
			ch <- m
		}
		close(ch)
	}()
	return ch
}

func Consumer(ch chan string) {
	for i := range ch {
		fmt.Println(i)
	}
}

// Ex 3
func produceNumbers(n int) chan int {
	ch := make(chan int)
	go func() {
		for i := 1; i <= n; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func consumeNumbers(nChan chan int, wg *sync.WaitGroup, totNumberWorked *int) {
	defer wg.Done()
	for i := range nChan {
		isPrime(i)
		*totNumberWorked++
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	// Ex 1
	fmt.Println("Ex 1")
	msg := []string{"The world itself's",
		"just one big hoax.",
		"Spamming each other with our",
		"running commentary of bullshit,",
		"masquerading as insight, our social media",
		"faking as intimacy.",
		"Or is it that we voted for this?",
		"Not with our rigged elections,",
		"but with our things, our property, our money.",
		"I'm not saying anything new.",
		"We all know why we do this,",
		"not because Hunger Games",
		"books make us happy,",
		"but because we wanna be sedated.",
		"Because it's painful not to pretend,",
		"because we're cowards.", "- Elliot Alderson",
		"Mr. Robot"}
	chProducer := Producer(msg)
	Consumer(chProducer)

	// Ex 3
	fmt.Println()
	fmt.Println("Ex 3")
	t := time.Now()
	chNumbers := produceNumbers(10000) // create a channel that will contain the numbers to test
	totNumberWorked := 0               // number of prime numbers
	var wg3 sync.WaitGroup
	for k := 0; k < 50; k++ { // create 50 workers
		wg3.Add(1)
		go consumeNumbers(chNumbers, &wg3, &totNumberWorked)
	}
	wg3.Wait()
	_, ok := <-chNumbers
	fmt.Println(ok)                                 // false
	fmt.Println("totNumberPassed", totNumberWorked) // not always 10000, sometimes 9992 ...
	fmt.Println("Ex3 time : ", time.Since(t))

	// Ex 3bis
	// version 1
	fmt.Println()
	fmt.Println("Ex 3bis")
	fmt.Println("version 1")
	t = time.Now()
	world1 := []int{1, 2, 3, 4, 5}
	newWorld1 := make([]int, 0, len(world1))
	var wg sync.WaitGroup
	for _, state := range world1 {
		wg.Add(1)
		go func(state int) {
			defer wg.Done()
			newWorld1 = append(newWorld1, 2*state)
		}(state)
	}
	wg.Wait()
	world1 = newWorld1
	fmt.Println(time.Since(t))
	fmt.Println(world1)

	// version 2
	fmt.Println()
	fmt.Println("version 2")
	t = time.Now()
	world2 := []int{1, 2, 3, 4, 5}
	newWorld2 := make([]int, 0, len(world2))
	jobs2 := make(chan int, len(world2))
	finishedChan := make(chan bool, len(world2))
	for k := 0; k < len(world2); k++ { // create workers
		go func(ch chan int) { // create a worker
			for i := range ch { // get jobs from the channel
				newWorld2 = append(newWorld2, 2*i) // do the job
				finishedChan <- true               // notify the main thread that the job is done
			}
		}(jobs2)
	}
	for _, state := range world2 { // give jobs to the workers trough the channel
		jobs2 <- state // send the job to the channel
	}
	close(jobs2)
	for k := 0; k < len(world2); k++ { // wait for all the workers to finish
		<-finishedChan
	}
	world2 = newWorld2
	fmt.Println(time.Since(t))
	fmt.Println(world2)

	// version 3
	fmt.Println()
	fmt.Println("version 3")
	t = time.Now()
	world3 := []int{1, 2, 3, 4, 5}
	newWorld3 := make([]int, 0, len(world3))
	jobs3 := make(chan int, len(world3))
	done := make(chan bool)
	go func() { // create a worker
		for { // get jobs from the channel
			i, more := <-jobs3 // get the job from the channel
			if more {          // if there is a job
				newWorld3 = append(newWorld3, 2*i) // do the job
			} else { // if there is no more job
				done <- true // notify the main thread that the job is done
				return
			}
		}
	}()
	for j := 1; j <= 3; j++ { // give jobs to the workers trough the channel
		jobs3 <- j // send the job to the channel
	}
	close(jobs3)
	<-done // wait for the workers to finish
	world3 = newWorld3
	fmt.Println(time.Since(t))
	fmt.Println(world3) // [2 4 6] pas complet

	// Ex 4
	fmt.Println()
	fmt.Println("Ex 4")

	producerChan := make(chan string)
	// Producer

	go func() {
		var randNumber int
		rand.Seed(10)
		for k := 0; k < 10; k++ { // Give 1000 jobs to the Plus and Minus workers
			randNumber = rand.Intn(2)
			if randNumber == 0 {
				fmt.Println(k, "Add + to producerChan")
				producerChan <- "+"
			} else {
				fmt.Println(k, "Add - to producerChan")
				producerChan <- "-"
			}
		}
		fmt.Println("close(producerChan)")
		close(producerChan)
	}()

	countPlus := 0
	countMinus := 0
	plusChan := make(chan string)
	minusChan := make(chan string)
	var waitGroup sync.WaitGroup
	// Plus worker
	waitGroup.Add(1)
	go func() {
		fmt.Println("Plus worker starts")
		s, ok1 := <-producerChan
		fmt.Println("+wk ok1", s, ok1)
		_, ok2 := <-plusChan
		fmt.Println("+wk ok2", ok2)
		fmt.Println("ok1", ok1, "ok2", ok2)
		for ok1 || ok2 {
			fmt.Println("ok1", ok1, "ok2", ok2)
			if ok1 {
				if s == "+" {
					fmt.Println("direct countPlus !", countPlus)
					countPlus++
				} else {
					fmt.Println("+worker Gives a - to -worker")
					minusChan <- "-"
				}
			}
			if ok2 {
				fmt.Println("inter countPlus !", countPlus)
				countPlus++
			}
			s, ok1 = <-producerChan
			_, ok2 = <-plusChan
		}
		fmt.Println("Plus worker is done !")
		waitGroup.Done()
	}()
	// Minus worker
	waitGroup.Add(1)
	go func() {
		fmt.Println("Minus worker starts")
		s, ok1 := <-producerChan
		fmt.Println("-wk", s, ok1)
		_, ok2 := <-plusChan
		fmt.Println("-wk ok2", ok2)
		fmt.Println("ok1", ok1, "ok2", ok2)
		for ok1 || ok2 {
			fmt.Println("ok1", ok1, "ok2", ok2)
			if ok1 {
				if s == "-" {
					fmt.Println("direct countMinus !", countMinus)
					countMinus++
				} else {
					fmt.Println("-worker Gives a + to +worker")
					plusChan <- "+"
				}
			}
			if ok2 {
				fmt.Println("inter countMinus !", countMinus)
				countMinus++
			}
			s, ok1 = <-producerChan
			_, ok2 = <-plusChan
		}
		fmt.Println("Minus worker is done !")
		waitGroup.Done()
	}()

	fmt.Println()
	waitGroup.Wait()
	fmt.Println("countPlus", countPlus)
	fmt.Println("countMinus", countMinus)
	fmt.Println("countPlus + countMinus", countPlus+countMinus)

}
