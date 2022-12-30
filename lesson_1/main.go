package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkUserAnswer(question string, answer int, answerChannel chan bool) {
	fmt.Printf("What is %s ?\n", question)
	var ans int
	fmt.Scan(&ans)
	if ans == answer {
		answerChannel <- true
	} else {
		answerChannel <- false
	}
}

func countTrueAnswers(allAnswers []bool) (int, int) {
	var correctAnswers int
	var totalAnswers int

	for i, ans := range allAnswers {
		if ans {
			correctAnswers++
		}

		totalAnswers = i + 1
	}

	return correctAnswers, totalAnswers
}

func parseArgs() (definedTimeout int, fileWithQuestions string) {
	pathToQuestionsFile := flag.String("questions", "problems.csv", "Path to the .csv file with questions")
	timeout := flag.Int("timeout", 30, "Timeout for one question")
	flag.Parse()
	return *timeout, *pathToQuestionsFile
}

func readQuestionsFromFile(pathToQuestionsFile string) ([]byte, error) {
	return os.ReadFile(pathToQuestionsFile)
}

func main() {
	oneQuestionTimeout, pathToQuestionFile := parseArgs()
	fmt.Printf("Defined timeout: %d\n", oneQuestionTimeout)

	dat, err := readQuestionsFromFile(pathToQuestionFile)

	check(err)

	r := csv.NewReader(strings.NewReader(string(dat)))

	var answers []bool
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(oneQuestionTimeout)*time.Second)
		defer cancel()
		var answerConvertedToInt int
		answerConvertedToInt, err = strconv.Atoi(record[1])
		check(err)
		fmt.Printf("Question nr %d\n", len(answers)+1)
		answerChannel := make(chan bool)
		answer := false
		go checkUserAnswer(record[0], answerConvertedToInt, answerChannel)
		select {
		case <-ctx.Done():
			fmt.Print("Time out!\n")
			answer = false
		case answerFromChannel := <-answerChannel:
			answer = answerFromChannel
		}

		answers = append(answers, answer)

	}
	correct, total := countTrueAnswers(answers)
	fmt.Printf("Total amount of: correct %d, total %d\n", correct, total)
}
