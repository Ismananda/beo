package beo

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Struktur utama AI
type AI struct {
	KnowledgeBase KnowledgeBase
	file          *os.File
}

// KnowledgeBase merepresentasikan database pertanyaan dan jawaban
type KnowledgeBase struct {
	AIName       string            `yaml:"name"`
	Model        string            `yaml:"model"`
	Trainer      string            `yaml:"trainer"`
	Fallbacks    Fallbacks         `yaml:"fallbacks"`
	Formats      Formats           `yaml:"formats"`
	Placeholders map[string]string `yaml:"placeholders"`
	Questions    []Question        `yaml:"questions"`
	Hooks        map[string]Hook   `yaml:"hooks"`

	IDF        map[string]float64 `yaml:"-"`
	Corpus     [][]string         `yaml:"-"`
	Vocabulary []string           `yaml:"-"`
}

// Formats merepresentasikan struktur format placeholder
type Formats struct {
	Date     string `yaml:"date"`
	Time     string `yaml:"time"`
	TimeZone string `yaml:"timezone"`
}

// Fallbacks merepresentasikan struktur fallback untuk berbagai kondisi
type Fallbacks struct {
	NoAnswer string `yaml:"noanswer"`
}

// Question merepresentasikan sebuah pertanyaan dan jawaban
type Question struct {
	Question string   `yaml:"question"`
	Answers  []string `yaml:"answers,omitempty"`
	Hook     string   `yaml:"hook,omitempty"`
}

// Hook merepresentasikan hook yang memiliki jawaban
type Hook struct {
	Answers []string `yaml:"answers"`
}

// Membuat AI baru dan memuat knowledge base
func NewAI(file *os.File) (*AI, error) {
	ai := &AI{
		KnowledgeBase: KnowledgeBase{
			Questions:    []Question{},
			Hooks:        make(map[string]Hook),
			Placeholders: make(map[string]string),
			AIName:       "Beo Talk",
			Model:        "BEE",
			Trainer:      "You",
			Formats: Formats{
				Date:     "02 Jan 2006",
				Time:     "15:04:05",
				TimeZone: "UTC",
			},
			Fallbacks: Fallbacks{
				NoAnswer: "I'm sorry, I don't know the answer to that.",
			},
		},
		file: file,
	}

	// Memuat knowledge base dari file
	err := ai.load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("gagal memuat knowledge base: %w", err)
	}

	return ai, nil
}

// Memuat knowledge base dari file
func (ai *AI) load() error {
	var kb KnowledgeBase

	stat, err := ai.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to check file status: %v", err)
	}
	if stat.Size() == 0 {
		return nil
	}

	decoder := yaml.NewDecoder(ai.file)
	err = decoder.Decode(&kb)
	if err != nil {
		return fmt.Errorf("failed to decode YAML file: %v", err)
	}

	// Pastikan semua field memiliki nilai awal
	if kb.Hooks == nil {
		kb.Hooks = make(map[string]Hook)
	}
	if kb.Placeholders == nil {
		kb.Placeholders = make(map[string]string)
	}
	if (kb.Formats == Formats{}) {
		kb.Formats = Formats{
			Date:     "02 Jan 2006",
			Time:     "15:04:05",
			TimeZone: "UTC",
		}
	}
	if (kb.Fallbacks == Fallbacks{}) {
		kb.Fallbacks = Fallbacks{
			NoAnswer: "I'm sorry, I don't know the answer to that.",
		}
	}

	kb.updateIDF()
	kb.updateVocabularies()
	ai.KnowledgeBase = kb
	return nil
}

// Menyimpan knowledge base ke file
func (ai *AI) Save() error {
	// Reset isi file sebelum menulis ulang
	err := ai.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("gagal menghapus isi file: %w", err)
	}
	_, err = ai.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("gagal mengatur posisi file: %w", err)
	}

	// Menulis ulang knowledge base ke file
	encoder := yaml.NewEncoder(ai.file)
	defer encoder.Close()

	return encoder.Encode(&ai.KnowledgeBase)
}

// Mencari jawaban terbaik berdasarkan pertanyaan
func (ai *AI) Ask(question string) string {
	var bestMatches []Question
	var answers []string

	segments := splitByPunctuation(question)
	for _, segment := range segments {
		// Tokenisasi dan koreksi typo
		inputTokens := tokenize(segment)
		correctedTokens := correctInput(inputTokens, ai.KnowledgeBase.Vocabulary)

		// Cari pola yang cocok
		bestMatches = append(bestMatches, findBestMatches(correctedTokens, ai.KnowledgeBase)...)
	}

	for _, bestMatch := range bestMatches {
		if bestMatch.Hook != "" {
			hook, ok := ai.KnowledgeBase.Hooks[bestMatch.Hook]
			if ok {
				answers = append(answers, randomChoice(hook.Answers))
			}
		} else {
			answers = append(answers, randomChoice(bestMatch.Answers))
		}
	}

	// Gunakan fallback untuk jawaban default
	if len(bestMatches) < 1 {
		return ai.KnowledgeBase.Fallbacks.NoAnswer
	}

	// Mengganti placeholders
	return processPlaceholders(strings.Join(answers, " "), ai.KnowledgeBase)
}

// Melatih AI dengan pertanyaan, jawaban, atau hook
func (ai *AI) Train(question string, answers []string, hook string) {
	for i, q := range ai.KnowledgeBase.Questions {
		if q.Question == question {
			// Tambahkan jawaban baru yang belum ada
			for _, answer := range answers {
				if !contains(q.Answers, answer) {
					ai.KnowledgeBase.Questions[i].Answers = append(q.Answers, answer)
				}
			}
			ai.KnowledgeBase.updateIDF()
			ai.KnowledgeBase.updateVocabularies()
			return
		}
	}

	ai.KnowledgeBase.Questions = append(ai.KnowledgeBase.Questions, Question{
		Question: question,
		Answers:  answers,
		Hook:     hook,
	})

	ai.KnowledgeBase.updateIDF()
	ai.KnowledgeBase.updateVocabularies()
}

// Menambahkan hook baru
func (ai *AI) AddHook(hookName string, answers []string) {
	if ai.KnowledgeBase.Hooks == nil {
		ai.KnowledgeBase.Hooks = make(map[string]Hook)
	}
	ai.KnowledgeBase.Hooks[hookName] = Hook{Answers: answers}
}

// Menambahkan placeholder baru
func (ai *AI) AddPlaceholder(key, value string) {
	if ai.KnowledgeBase.Placeholders == nil {
		ai.KnowledgeBase.Placeholders = make(map[string]string)
	}
	ai.KnowledgeBase.Placeholders[key] = value
}

// updateIDF menghitung dan memperbarui nilai Inverse Document Frequency (IDF) di dalam KnowledgeBase.
func (kb *KnowledgeBase) updateIDF() {
	corpus := [][]string{}
	for _, question := range kb.Questions {
		corpus = append(corpus, tokenize(question.Question))
	}

	kb.Corpus = corpus
	kb.IDF = inverseDocumentFrequency(corpus)
}

// updateVocabularies memperbarui daftar kosakata (Vocabulary) di dalam KnowledgeBase.
func (kb *KnowledgeBase) updateVocabularies() {
	uniqueVocabularies := map[string]bool{}

	for _, question := range kb.Questions {
		for _, word := range tokenize(question.Question) {
			uniqueVocabularies[word] = true
		}
	}

	var vocabularyList []string
	for word := range uniqueVocabularies {
		vocabularyList = append(vocabularyList, word)
	}

	kb.Vocabulary = vocabularyList
}
