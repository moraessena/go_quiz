package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the limit for the quiz in seconds")
	flag.Parse()
	file := getFile(*fileName)
	lines := readLines(file)
	problems := parseLines(lines)
	runQuestions(problems, time.Duration(*timeLimit))
}

func runQuestions(problems []problem, duration time.Duration) {
	timer := time.NewTimer(duration * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func getFile(f string) *os.File {
	file, err := os.Open(f)
	if err != nil {
		exit(f)
	}
	return file
}

func readLines(file io.Reader) [][]string {
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file")
	}
	return lines
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: strings.TrimSpace(line[0]),
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(m string) {
	fmt.Println("Failed to open the CSV file:", m)
	os.Exit(1)
}
