package main

import (
	"log"
	"sync"
)

// want -1, 0, 1, 2, 3, 3, 5, 9, 10, 10, 12, 15
func mergeChannels(ch1 <-chan int64, ch2 <-chan int64) <-chan int64 {
	var (
		sources = []<-chan int64{ch1, ch2}
		nums    = [][]int64{{}, {}}
		wg      = &sync.WaitGroup{}
		merged  = make(chan int64)
	)

	go func() {
		wg.Add(1)
		for idx, source := range sources {
			for value := range source {
				nums[idx] = append(nums[idx], value)
			}
		}

		result := merge(nums[0], nums[1])
		go func() {
			defer wg.Done()

			for _, value := range result {
				merged <- value
			}

		}()

		wg.Wait()
		close(merged)

	}()

	return merged
}

func initChannel(data []int64) chan int64 {
	ch := make(chan int64, len(data))
	for _, d := range data {
		ch <- d
	}
	close(ch)

	return ch
}

func merge(first []int64, second []int64) []int64 {
	var (
		i   = len(first) - 1
		j   = len(second) - 1
		end = len(first) + len(second) - 1
	)

	for k := 1; k <= len(second); k++ {
		first = append(first, 0)
	}

	for j >= 0 {
		if i >= 0 && first[i] > second[j] {
			first[end] = first[i]
			i--
		} else {
			first[end] = second[j]
			j--
		}
		end--
	}

	return first
}

func main() {
	ch1 := initChannel([]int64{1, 2, 3, 10, 12, 15})
	ch2 := initChannel([]int64{-1, 0, 3, 5, 9, 10})

	for value := range mergeChannels(ch1, ch2) {
		log.Printf("got value: %d", value)
	}

	log.Print("done")
}
