package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

func readFile(filePath string) (chan string, error) {
	fd, err := os.Open(filePath)
	r := csv.NewReader(fd)

	if err != nil {
		fmt.Println("Error")
		return nil, err
	}

	out := make(chan string)
	go parseCsvData(r, out)
	return out, nil
}

func parseCsvData(r *csv.Reader, ch chan string) {
	defer close(ch)

	// ignore headers row
	_, err := r.Read()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	var totalData string
	const (
		title  = 1
		author = 2
	)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading record:", err)
			continue
		}
		totalData += fmt.Sprintf("%v (%v)\n", record[title], record[author])
	}
	fmt.Printf("Sending all data from file: \n%v\n", totalData)
	ch <- totalData
}

func aggregateData(channels ...chan string) chan string {
	var wg sync.WaitGroup
	rsch := make(chan string)

	go combine(rsch, &wg, channels...)
	return rsch
}

func combine(ch chan string, wg *sync.WaitGroup, channels ...chan string) {
	var aggregatedData string
	for _, ch := range channels {
		aggregatedData += <-ch
	}
	wg.Wait()

	ch <- aggregatedData
	close(ch)
}

func main() {
	fileNames := []string{"static/input1.csv", "static/input2.csv"}
	var channels []chan string

	for _, fileName := range fileNames {
		channel, err := readFile(fileName)
		if err != nil {
			panic(fmt.Errorf("Could not read file %v. Error: %v", fileName, err))
		}
		channels = append(channels, channel)
	}

	result := aggregateData(channels...)
	for data := range result {
		fmt.Println(data)
	}
}
