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
)

func main() {
	readFromCommandLine()
}

func readFromCommandLine() {
	fileToUse := flag.String("file", "./problems.csv", "specify which csv file to use for the quiz")
	timeLimit := flag.Int("time", 10, "set a time limit for a question to be answered in")
	flag.Parse()
	file, err := os.Open(*fileToUse)
	defer file.Close()
	if err != nil {
		log.Fatalln("couldn't open csv file", err)
	}
	questions := csv.NewReader(file)
	scanner := bufio.NewScanner(os.Stdin)
	correctAnswers := 0
	totalAmountOfQuestions := 0

	checkTime := make(chan string)
	response := make(chan bool)
	for {
			question, err := questions.Read()
			if err == io.EOF {
				fmt.Println("Your final score is:", correctAnswers, "(out of", totalAmountOfQuestions, "questions)")
				break
			}
			fmt.Println(question[0])
			go timer(*timeLimit, checkTime, response)
			scanner.Scan()
			answer := scanner.Text()
			fmt.Println(answer, "answer")
			checkTime <- answer
			fmt.Println("kfowief")
			value := <- response
			fmt.Println("if yu see This it doesnt block")
			if !value {
				fmt.Println("You were too slow!")
				break
			}

			if question[1] == answer {
				correctAnswers += 1
			}
			totalAmountOfQuestions += 1

		}
	}

func timer(timeLimit int, answer <-chan string, response chan bool) {
	timer := time.NewTicker(time.Duration(timeLimit))
	fmt.Println("here")

	select {
	case value := <-answer:
		fmt.Println("here bitch")
		if value == "" {
			response <- false
		} else {
			response <- true
		}
		timer.Stop()
		break
	case <-timer.C:
		response <- false
		timer.Stop()
		break



	}

}