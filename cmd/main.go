package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Ismananda/beo"
)

func main() {
	const filename = "model.yml"
	const help = `
Use --ask, --train, --hook, or --placeholder
Examples:
--ask "What is AI?"
--train "What is AI?" "Artificial Intelligence"
--hook "greet" "Hello" "Hi"
--placeholder "date" "02 Jan 2006"
`

	if len(os.Args) < 2 {
		fmt.Print(help)
		return
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	ai, err := beo.NewAI(file)
	if err != nil {
		fmt.Printf("Failed to load model: %v\n", err)
		return
	}

	switch os.Args[1] {
	case "--ask":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a question.")
			return
		}
		question := strings.Join(os.Args[2:], " ")
		answer := ai.Ask(question)
		fmt.Println("Answer:", answer)

	case "--train":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a question and answers or a hook.")
			return
		}
		question := os.Args[2]
		lastArg := os.Args[len(os.Args)-1]

		var answers []string
		hook := ""
		if strings.HasPrefix(lastArg, "hook:") {
			hook = strings.TrimPrefix(lastArg, "hook:")
			answers = os.Args[3 : len(os.Args)-1]
		} else {
			answers = os.Args[3:]
		}

		ai.Train(question, answers, hook)
		if err := ai.Save(); err != nil {
			fmt.Printf("Failed to save model: %v\n", err)
			return
		}
		fmt.Println("Model successfully trained.")

	case "--hook":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a hook name and answers.")
			return
		}
		hookName := os.Args[2]
		answers := os.Args[3:]

		ai.AddHook(hookName, answers)
		if err := ai.Save(); err != nil {
			fmt.Printf("Failed to save model: %v\n", err)
			return
		}
		fmt.Println("Hook successfully added.")

	case "--placeholder":
		if len(os.Args) < 4 {
			fmt.Println("Please provide a placeholder name and its value.")
			return
		}
		key := os.Args[2]
		value := os.Args[3]

		ai.AddPlaceholder(key, value)
		if err := ai.Save(); err != nil {
			fmt.Printf("Failed to save placeholder: %v\n", err)
			return
		}
		fmt.Println("Placeholder successfully added.")

	default:
		fmt.Print("Unknown command.", help)
	}
}
