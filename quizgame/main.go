package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvPath, timeLimit := handleFlags()
	data, err := parseCsv(csvPath)
	if err != nil {
		fmt.Println("csv parsing error:", err)
		return
	}
	err = validateCsv(data)
	if err != nil {
		fmt.Println("csv validating error:", err)
		return
	}
	score, maxScore := startQuiz(data, timeLimit)
	fmt.Printf("You scored %d out of %d.\n", score, maxScore)
	if score == maxScore {
		fmt.Println("Good job!")
	}
}

func handleFlags() (csvPath string, timeLimit int) {
	flag.StringVar(&csvPath, "csv", "problems.csv", "a csv file in the format of 'question,answer' (default \"problems.csv\")")
	flag.IntVar(&timeLimit, "limit", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()
	return
}

func parseCsv(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	result, err := reader.ReadAll()
	return result, err
}

func validateCsv(list [][]string) error {
	if len(list) == 0 {
		return errors.New("no valid entries in the file")
	}
	if len(list[0]) != 2 {
		return errors.New("should be 2 columns in the file")
	}
	return nil
}

func startQuiz(data [][]string, timeLimit int) (score int, maxScore int) {
	score = 0
	maxScore = len(data)
	chInput := make(chan string)
	chTimerRestart := make(chan bool)
	chTimerTimeout := make(chan bool)
	go userInput(chInput, data, &score)
	go timer(timeLimit, chTimerRestart, chTimerTimeout)
	for {
		select {
		case _, ok := <-chInput:
			chTimerRestart <- true
			if !ok {
				return
			}
		case <-chTimerTimeout:
			fmt.Println("\nTime out")
			return
		}
	}
}

func userInput(out chan<- string, data [][]string, score *int) {
	*score = 0
	i := 1
	for _, v := range data {
		fmt.Printf("Problem #%d: %s\n", i, v[0])
		answer := ""
		for answer == "" {
			fmt.Scanln(&answer)
			answer = strings.TrimSpace(answer)
			out <- answer
		}
		if answer == v[1] {
			*score++
		}
		i++
	}
	close(out)
	return
}

func timer(timeLimit int, restart <-chan bool, timeout chan<- bool) {
	startTime := time.Now()
	limit := time.Duration(timeLimit * int(time.Second))
	for time.Since(startTime) < limit {
		select {
		case <-restart:
			startTime = time.Now()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	timeout <- true
}
