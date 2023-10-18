package usecase

import (
	"fmt"
	"sync"
)

func StartGor(withChannel bool, withWG bool, withMutex bool) {
	if withChannel {
		gorWithChannel()
	} else if withWG {
		gorWithWG()
	} else if withMutex {
		gorWithMutex()
	} else {
		simpleGor()
	}
}

func simpleGor() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(fmt.Sprintf("Горутина №%d", i))
		}(i)
	}
}

func gorWithChannel() {
	gorNums := make(chan int)
	gorValues := make(chan string)

	for i := 0; i < 10; i++ {
		go func(i int) {
			gorNums <- i
		}(i)
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			gorValues <- fmt.Sprintf("Значение горутины №%d", i)
		}(i)
	}

	// основной метод для чтения из канала
	/*val, ok := <-gorValues
	if ok {
		fmt.Println(val)
	}*/

	// Для чтения из одного канала
	/*for gorNum := range gorNums {
		fmt.Println(fmt.Sprintf("Горутина №%d", gorNum))
	}*/

	// Для чтения из нескольких каналов
	/*for {
		select {
		case num := <-gorNums:
			fmt.Println(fmt.Sprintf("Горутина №%d", num))
		case val := <-gorValues:
			fmt.Println(val)
		}
	}*/

	// гоуртина для чтения из канала
	/*go func(gorValues chan string) {
		for {
			select {
			case val := <-gorValues:
				fmt.Println(val)
			}
		}
	}(gorValues)*/
}

func gorWithWG() {
	var wg sync.WaitGroup
	gorNums := make(chan int)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			gorNums <- i
		}(i, &wg)
	}

	for gorNum := range gorNums {
		fmt.Println(fmt.Sprintf("Горутина №%d", gorNum))
	}

	wg.Wait()
}

type Container struct {
	mu       sync.Mutex
	counters map[string]int
}

func gorWithMutex() {
	c := Container{
		counters: map[string]int{"a": 0, "b": 0},
	}

	var wg sync.WaitGroup

	doIncrement := func(name string, n int) {
		for i := 0; i < n; i++ {
			c.inc(name)
		}
		wg.Done()
	}

	wg.Add(3)
	go doIncrement("a", 10000)
	go doIncrement("a", 10000)
	go doIncrement("b", 10000)

	wg.Wait()
	fmt.Println(c.counters)
}

func (c *Container) inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counters[name]++
}
