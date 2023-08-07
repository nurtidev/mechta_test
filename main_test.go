package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestReadItemsFromFile(t *testing.T) {
	// Создадим временный файл с тестовыми данными
	tempFile, err := ioutil.TempFile("", "test_data.json")
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Запишем тестовые данные во временный файл
	testData := []Item{{A: 1, B: 2}, {A: 3, B: 4}, {A: 5, B: 6}}
	data, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Ошибка при маршалинге тестовых данных: %v", err)
	}
	if _, err := tempFile.Write(data); err != nil {
		t.Fatalf("Ошибка при записи во временный файл: %v", err)
	}
	tempFile.Close()

	// Вызываем функцию для чтения данных из временного файла
	items, err := readItemsFromFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Ошибка при чтении данных из файла: %v", err)
	}

	// Проверяем, что данные были успешно прочитаны и соответствуют ожиданиям
	expectedData := testData
	if len(items) != len(expectedData) {
		t.Fatalf("Ожидается %d элементов, получено %d", len(expectedData), len(items))
	}
	for i := 0; i < len(expectedData); i++ {
		if items[i] != expectedData[i] {
			t.Fatalf("Ошибка в элементе %d. Ожидается %+v, получено %+v", i, expectedData[i], items[i])
		}
	}
}

func TestCalculateSum(t *testing.T) {
	// Генерируем тестовые данные
	testData := []Item{{A: 1, B: 2}, {A: 3, B: 4}, {A: 5, B: 6}}

	// Вычисляем ожидаемую сумму
	expectedSum := 0
	for _, item := range testData {
		expectedSum += item.A + item.B
	}

	// Вызываем функцию для вычисления суммы
	actualSum := calculateSum(testData, 1, 10)

	// Проверяем, что результат соответствует ожиданиям
	if actualSum != expectedSum {
		t.Fatalf("Ошибка в вычислении суммы. Ожидается %d, получено %d", expectedSum, actualSum)
	}
}

func TestCalculateSumWithMultipleGoroutines(t *testing.T) {
	// Генерируем тестовые данные
	testData := make([]Item, 1000)
	for i := 0; i < len(testData); i++ {
		testData[i] = Item{A: rand.Intn(19) - 9, B: rand.Intn(19) - 9}
	}

	// Вычисляем ожидаемую сумму
	expectedSum := 0
	for _, item := range testData {
		expectedSum += item.A + item.B
	}

	// Попробуем разное количество горутин
	for numGoroutines := 1; numGoroutines <= 10; numGoroutines++ {
		// Вызываем функцию для вычисления суммы с разным количеством горутин
		actualSum := calculateSum(testData, numGoroutines, 10)

		// Проверяем, что результат соответствует ожиданиям
		if actualSum != expectedSum {
			t.Fatalf("Ошибка в вычислении суммы с %d горутинами. Ожидается %d, получено %d",
				numGoroutines, expectedSum, actualSum)
		}
	}
}

func TestCalculateSumWithEmptyData(t *testing.T) {
	// Генерируем тестовые данные
	testData := []Item{}

	// Вызываем функцию для вычисления суммы
	actualSum := calculateSum(testData, 5, 10)

	// Проверяем, что результат равен 0, так как данных нет
	if actualSum != 0 {
		t.Fatalf("Ошибка в вычислении суммы пустого набора данных. Ожидается 0, получено %d", actualSum)
	}
}
