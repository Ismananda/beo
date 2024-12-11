# Beo

Beo is a lightweight AI framework written in Go, designed for simple question-and-answer systems. This is not a machine learning model but a rule-based engine utilizing a combination of TF-IDF (Term Frequency-Inverse Document Frequency), Levenshtein Distance, and cosine similarity. Beo is typo-tolerant and suitable for basic automated response systems. The framework emphasizes simplicity and provides flexible configuration options for AI behavior.

The underlying model, named **BEE** (Basic Expression Evaluator), powers Beo's functionality. BEE combines the aforementioned techniques to ensure efficient and reliable responses for various textual inputs.

## Features
- **Custom Question-Answer Pairs**: Train Beo with your own questions and answers.
- **Reusable Hooks**: Define hooks for reusable, predefined responses.
- **Dynamic Placeholders**: Use placeholders to create adaptable responses.
- **Persistence**: Save and load Beo's knowledge base in YAML format.
- **Customizable Logic**: Easily extend or modify the core behavior.
- **Multi-Question Handling**: Parse and answer multiple questions in a single query.

## Getting Started

### Installation
1. Clone this repository:
   ```bash
   git clone https://github.com/Ismananda/beo.git
   cd beo
   ```

2. Train a question and answer:
   ```bash
   go run cmd/main.go --train "Hello, who are you?" "I am Beo."
   ```

3. Run a query on trained data:
   ```bash
   go run cmd/main.go --ask "Helo, who are you?"
   ```

### Using Beo

- **CLI Application**: Interact with Beo through commands like `--train` and `--ask` after compiling the project.
- **Go Module**: Integrate Beo into your Go projects as a module. See examples in the [Usage](#usage) section below.

---

## Usage

### Training the AI
Use the `Train` function to add questions, answers, and optional hooks to Beo's knowledge base.

Example:
```go
ai.Train("What is your name?", []string{"I am Beo."}, "")
ai.Train("What is the capital of France?", []string{"Paris"}, "")
```

### Asking Questions
Use the `Ask` function to query Beo and receive answers.

Example:
```go
answer := ai.Ask("What is your name?")
fmt.Println(answer) // Output: I am Beo.
```

### Adding Hooks
Define reusable hooks with predefined responses using `AddHook`.

Example:
```go
ai.AddHook("greeting", []string{"Hello!", "Hi there!"})
```

### Adding Placeholders
Dynamically define placeholders for use in responses.

Example:
```go
ai.AddPlaceholder("name", "Beo")
```

### Saving and Loading
Save Beo's knowledge base to a file:
```go
if err := ai.Save(); err != nil {
    log.Fatalf("Error saving knowledge base: %v", err)
}
```

Load Beo's knowledge base from a file:
```go
ai, err := beo.NewAI(file)
if err != nil {
    log.Fatalf("Error loading AI: %v", err)
}
```

### Handling Multiple Questions
Beo can split inputs based on punctuation marks (e.g., `.`, `?`, `!`) to handle multiple questions in one query.

Example:
```go
input := "Who are you? What is your purpose?"
responses := ai.Ask(input)
for _, response := range responses {
    fmt.Println(response)
}
// Output:
// I am Beo. My purpose is to assist with basic queries.
```

---

## Example `main.go`
Below is an example of how to set up and run Beo. This file is located in the `cmd/` directory.

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Ismananda/beo"
)

func main() {
    // Initialize Beo with a temporary knowledge base file
    file, err := os.CreateTemp("", "knowledgebase_*.yml")
    if err != nil {
        log.Fatalf("Error creating temp file: %v", err)
    }
    defer os.Remove(file.Name())

    ai, err := beo.NewAI(file)
    if err != nil {
        log.Fatalf("Error initializing Beo: %v", err)
    }

    // Train Beo
    ai.Train("What is your name?", []string{"I am Beo."}, "")
    ai.Train("What is the capital of France?", []string{"Paris"}, "")

    // Add a hook
    ai.AddHook("greeting", []string{"Hello!", "Hi there!"})

    // Add a placeholder
    ai.AddPlaceholder("name", "Beo")

    // Save the knowledge base
    if err := ai.Save(); err != nil {
        log.Fatalf("Error saving knowledge base: %v", err)
    }

    // Ask a question
    answer := ai.Ask("What is your name?")
    fmt.Println(answer)
}
```

---

## Configuration Example

Below is an example of a model file for Beo, written in YAML format. It demonstrates how to define questions, answers, hooks, and placeholders:

```yaml
name: Beo AI
model: BEE
trainer: You
fallbacks:
    noanswer: Sorry, I don't understand your question.
formats:
    date: "02 Jan 2006"
    time: "15:04:05"
    timezone: UTC
placeholders:
    test: This is a test placeholder
questions:
    - question: What is your name?
      answers:
        - Hi %user%, I am %ainame%.
        - %ainame%.
    - question: Test
      answers:
        - This is a text, and %test%.
    - question: How's your day?
      hook: status
hooks:
    status:
        answers:
            - Just another regular day.
            - My day is very good.
```

---

## Running Tests
Unit tests are located in the `test/` folder. Run the tests using:
```bash
go test ./test
```

---

## Contribution
Contributions are welcome! Submit issues or pull requests to improve the framework.

---

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Notes
- This project was created out of curiosity about AI systems and is not intended for production use.
- Almost all code, including this `README.md`, was generated with assistance from ChatGPT.
- The project is no longer actively maintained. Feel free to fork it, provide suggestions, or extend its functionality through pull requests or issues.

---