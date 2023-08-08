package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type Item struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	//if err := changeData(); err != nil {
	//	log.Fatal(err)
	//}
	numGoroutines := flag.Int("goroutines", 4, "Number of goroutines")
	numBlocks := flag.Int("blocks", 10, "Number of blocks")
	flag.Parse()

	if *numGoroutines <= 0 || *numBlocks <= 0 {
		log.Fatal("Number of goroutines and blocks must be greater than 0")
		return
	}

	items, err := readItemsFromFile("data.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	totalSum := calculateSum(items, *numGoroutines, *numBlocks)

	fmt.Println("Общая сумма всех чисел:", totalSum)
}

func calculateSum(items []Item, numGoroutines int, numBlocks int) int {
	blockSize := len(items) / numBlocks
	results := make(chan int, numBlocks)
	blocks := make(chan []Item, numBlocks)
	var wg sync.WaitGroup

	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := start + blockSize
		if i == numBlocks-1 && len(items)%numBlocks != 0 {
			end = len(items)
		}
		blocks <- items[start:end]
	}
	close(blocks)

	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go worker(blocks, &wg, results)
	}

	wg.Wait()
	close(results)

	totalSum := 0
	for sum := range results {
		totalSum += sum
	}

	// Замерим память
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Использовано памяти (в байтах): %d\n", mem.Alloc)

	// Время выполнения
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения: %v\n", duration)

	return totalSum
}

func worker(blocks <-chan []Item, wg *sync.WaitGroup, results chan<- int) {
	defer wg.Done()
	for block := range blocks {
		sum := 0
		for _, item := range block {
			sum += item.A + item.B
		}
		results <- sum
	}
}

func readItemsFromFile(fileName string) ([]Item, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var items []Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func changeData() error {
	rand.Seed(time.Now().UnixNano())

	items := make([]Item, 1000000)
	for i := 0; i < len(items); i++ {
		items[i] = Item{
			A: rand.Intn(19) - 9,
			B: rand.Intn(19) - 9,
		}
	}

	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(items)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return err
	}

	fmt.Println("Данные успешно записаны в файл data.json")
	return nil
}
