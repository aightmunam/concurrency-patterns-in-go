package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

func readFromFile(filePath string, ch chan string, wg *sync.WaitGroup) {
	fd, err := os.Open(filePath)
	defer fd.Close()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	r := csv.NewReader(fd)

	// ignore headers row
	_, err = r.Read()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var totalData string
	for {
		const title = 1
		const author = 2
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		totalData += fmt.Sprintf("%v (%v)\n", record[title], record[author])
	}

	defer wg.Done()
	ch <- totalData
}

func AggregateData(ch chan string, rsch chan string) {
	var aggregatedData string
	for data := range ch {
		aggregatedData += data
	}

	rsch <- aggregatedData
}

func main() {
	fileNames := []string{"static/input1.csv", "static/input2.csv"}
	dataChannel := make(chan string, len(fileNames))

	var wg sync.WaitGroup
	for _, fileName := range fileNames {
		wg.Add(1)
		go readFromFile(fileName, dataChannel, &wg)
	}

	wg.Wait()
	close(dataChannel) // Close the channel after all subtasks are done

	rsch := make(chan string, 1)
	go AggregateData(dataChannel, rsch)
	aggregatedData := <-rsch
	close(rsch)

	fmt.Printf("Final Data: \n%v", aggregatedData)

}
