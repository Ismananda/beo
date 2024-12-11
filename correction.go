package beo

// Koreksi kata berdasarkan Levenshtein Distance
func correctWord(word string, vocabulary []string) string {
	minDistance := 5
	corrected := word
	for _, vocabWord := range vocabulary {
		distance := levenshtein(word, vocabWord)
		if distance < minDistance && distance <= 2 {
			minDistance = distance
			corrected = vocabWord
		}
	}
	return corrected
}

// Koreksi seluruh input
func correctInput(input []string, vocabulary []string) []string {
	corrected := []string{}
	for _, word := range input {
		corrected = append(corrected, correctWord(word, vocabulary))
	}
	return corrected
}
