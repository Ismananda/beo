package test

import (
	"os"
	"testing"

	"github.com/Ismananda/beo" // Ganti dengan path import yang sesuai
)

// Test fungsi NewAI untuk memastikan AI dapat dibuat dan knowledge base dimuat dengan benar
func TestNewAI(t *testing.T) {
	file, err := os.CreateTemp("", "knowledgebase_test_*.yml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	ai, err := beo.NewAI(file)
	if err != nil {
		t.Fatalf("Error initializing AI: %v", err)
	}

	if len(ai.KnowledgeBase.Questions) != 0 {
		t.Errorf("Expected 0 questions, got %d", len(ai.KnowledgeBase.Questions))
	}

	if ai.KnowledgeBase.Fallbacks.NoAnswer != "I'm sorry, I don't know the answer to that." {
		t.Errorf("Expected default response, got %v", ai.KnowledgeBase.Fallbacks.NoAnswer)
	}
}

// Test fungsi Ask untuk memastikan AI memberikan jawaban yang benar
func TestAsk(t *testing.T) {
	file, err := os.CreateTemp("", "knowledgebase_test_*.yml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	ai, err := beo.NewAI(file)
	if err != nil {
		t.Fatalf("Error initializing AI: %v", err)
	}

	ai.Train("What is your name?", []string{"My name is TestBot."}, "")
	if err := ai.Save(); err != nil {
		t.Fatalf("Error saving AI model: %v", err)
	}

	answer := ai.Ask("What is your name?")
	expectedAnswer := "My name is TestBot."
	if answer != expectedAnswer {
		t.Errorf("Expected answer %v, but got %v", expectedAnswer, answer)
	}
}

// Test fungsi Train untuk memastikan pertanyaan dan jawaban ditambahkan dengan benar
func TestTrain(t *testing.T) {
	file, err := os.CreateTemp("", "knowledgebase_test_*.yml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	ai, err := beo.NewAI(file)
	if err != nil {
		t.Fatalf("Error initializing AI: %v", err)
	}

	ai.Train("What is the capital of France?", []string{"Paris"}, "")
	if err := ai.Save(); err != nil {
		t.Fatalf("Error saving AI model: %v", err)
	}

	if len(ai.KnowledgeBase.Questions) != 1 {
		t.Errorf("Expected 1 question, but got %d", len(ai.KnowledgeBase.Questions))
	}

	answer := ai.Ask("What is the capital of France?")
	expectedAnswer := "Paris"
	if answer != expectedAnswer {
		t.Errorf("Expected answer %v, but got %v", expectedAnswer, answer)
	}
}

// Test fungsi AddHook untuk memastikan hook dapat ditambahkan dengan benar
func TestAddHook(t *testing.T) {
	file, err := os.CreateTemp("", "knowledgebase_test_*.yml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	ai, err := beo.NewAI(file)
	if err != nil {
		t.Fatalf("Error initializing AI: %v", err)
	}

	ai.AddHook("greeting", []string{"Hello!", "Hi there!"})
	if err := ai.Save(); err != nil {
		t.Fatalf("Error saving AI model: %v", err)
	}

	if len(ai.KnowledgeBase.Hooks) != 1 {
		t.Errorf("Expected 1 hook, but got %d", len(ai.KnowledgeBase.Hooks))
	}

	hookAnswers := ai.KnowledgeBase.Hooks["greeting"].Answers
	if len(hookAnswers) != 2 {
		t.Errorf("Expected 2 answers in hook, but got %d", len(hookAnswers))
	}
}

// Test fungsi AddPlaceholder untuk memastikan placeholder dapat ditambahkan dengan benar
func TestAddPlaceholder(t *testing.T) {
	file, err := os.CreateTemp("", "knowledgebase_test_*.yml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	ai, err := beo.NewAI(file)
	if err != nil {
		t.Fatalf("Error initializing AI: %v", err)
	}

	ai.AddPlaceholder("name", "TestBot")
	if err := ai.Save(); err != nil {
		t.Fatalf("Error saving AI model: %v", err)
	}

	if len(ai.KnowledgeBase.Placeholders) != 1 {
		t.Errorf("Expected 1 placeholder, but got %d", len(ai.KnowledgeBase.Placeholders))
	}

	value := ai.KnowledgeBase.Placeholders["name"]
	expectedValue := "TestBot"
	if value != expectedValue {
		t.Errorf("Expected placeholder value %v, but got %v", expectedValue, value)
	}
}
