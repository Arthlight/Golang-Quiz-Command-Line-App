package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"strconv"
)

func main() {
	readFromCommandLine()
}

func readFromCommandLine() {
	fileToUse := flag.String("file", "./problems.csv", "specify which csv file to use for the quiz")
	timeLimit := flag.Int("time", 10, "set a time limit for a question to be answered in")
	flag.Parse()
	file, err := os.Open(*fileToUse)
	if err != nil {
		log.Fatalln("couldn't open csv file", err)
	}
	defer file.Close()
	questions := csv.NewReader(file)
	scanner := bufio.NewScanner(os.Stdin)
	correctAnswers := 0
	totalAmountOfQuestions := 0

	answer := make(chan string)
	response := make(chan interface{})

loop:
	for {
			question, err := questions.Read()
			if err == io.EOF {
				fmt.Println("Your final score is:", correctAnswers, "(out of", totalAmountOfQuestions, "questions)")
				break
			}
			fmt.Println(question[0])
			go timer(*timeLimit, answer, response)
			go scan(scanner, answer)

			value := <- response

			switch userAnswer := value.(type) {
			case bool:
				fmt.Println("You were too slow! But nice try ( ͡~ ͜ʖ ͡°)\nYour final score is:", correctAnswers)
				break loop

			case string:
				_, err := strconv.Atoi(userAnswer)
				if err != nil {
					fmt.Printf("This is not a number! You need to type a number! Preferably the correct one ( ͡° ͜ʖ ͡°)")
					break loop
				}
			}

			if question[1] == value {
				correctAnswers += 1
			}
			totalAmountOfQuestions += 1

		}
	}

func timer(timeLimit int, answer chan string, response chan<- interface{}) {
	timer := time.NewTicker(time.Second * time.Duration(timeLimit))

	select {
	case userAnswer := <-answer:
		response <- userAnswer
		timer.Stop()
	case <-timer.C:
		response <- false
		timer.Stop()
	}

}

func scan(scanner *bufio.Scanner, answer chan<- string) {
	scanner.Scan()
	userAnswer := scanner.Text()
	answer <- userAnswer
}