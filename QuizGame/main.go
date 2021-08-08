package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "un archivo csv que tiene el formato de 'pregunta,respuesta'")
	timeLimit := flag.Int("limite", 5, "el tiempo limite para responder en segundos")
	flag.Parse()
	_ = csvFilename

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("falla al abreir el archivo CSV: %s", *csvFilename))
		os.Exit(1)
		_ = file
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("falla al parsear el archivo CSV"))
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problema #%d: %s = ", i+1, p.question)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTu puntaje es %d de un maximo de %d.\n", correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("\nTu puntaje es %d de un maximo de %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
