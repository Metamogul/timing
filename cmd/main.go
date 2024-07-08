package main

import (
	"container/heap"
	"fmt"
	"time"
)

// Item represents an item with a timestamp from a channel
type Item struct {
	Timestamp time.Time
	Value     interface{}
	ChannelID int // Index of the channel this item came from
}

// PriorityQueue implements heap.Interface and holds Items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Timestamp.Before(pq[j].Timestamp)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// MergeChannels merges multiple timestamped channels into one sorted by timestamp
func MergeChannels(channels ...<-chan Item) <-chan Item {
	output := make(chan Item)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Initialize the heap with the first item from each channel
	for i, ch := range channels {
		if item, ok := <-ch; ok {
			heap.Push(&pq, &Item{
				Timestamp: item.Timestamp,
				Value:     item.Value,
				ChannelID: i,
			})
		}
	}

	go func() {
		defer close(output)
		for pq.Len() > 0 {
			// Pop the smallest item from the heap
			item := heap.Pop(&pq).(*Item)
			output <- *item

			// Push the next item from the same channel to the heap
			if nextItem, ok := <-channels[item.ChannelID]; ok {
				heap.Push(&pq, &Item{
					Timestamp: nextItem.Timestamp,
					Value:     nextItem.Value,
					ChannelID: item.ChannelID,
				})
			}
		}
	}()

	return output
}

func main() {
	// Example usage with dummy data

	ch1 := make(chan Item)
	ch2 := make(chan Item)
	ch3 := make(chan Item)

	go func() {
		defer func() {
			close(ch1)
			fmt.Println("ch1 closed")
		}()
		ch1 <- Item{Timestamp: time.Now().Add(1 * time.Second), Value: "A1"}
		ch1 <- Item{Timestamp: time.Now().Add(3 * time.Second), Value: "A2"}
		ch1 <- Item{Timestamp: time.Now().Add(4 * time.Second), Value: "A3"}
		ch1 <- Item{Timestamp: time.Now().Add(5 * time.Second), Value: "A4"}
		ch1 <- Item{Timestamp: time.Now().Add(6 * time.Second), Value: "A5"}
		ch1 <- Item{Timestamp: time.Now().Add(7 * time.Second), Value: "A6"}
	}()

	go func() {
		defer func() {
			close(ch2)
			fmt.Println("ch2 closed")
		}()
		ch2 <- Item{Timestamp: time.Now().Add(2 * time.Second), Value: "B1"}
		ch2 <- Item{Timestamp: time.Now().Add(4 * time.Second), Value: "B2"}
	}()

	go func() {
		defer func() {
			close(ch3)
			fmt.Println("ch3 closed")
		}()
		ch3 <- Item{Timestamp: time.Now().Add(2 * time.Second), Value: "C1"}
	}()

	output := MergeChannels(ch1, ch2, ch3)

	for item := range output {
		fmt.Printf("Received item: %v\n", item)
	}
}
