package beo

import (
	"math"
)

func findBestMatches(inputTokens []string, kb KnowledgeBase) []Question {
	matches := []Question{}
	usedTokens := make([]bool, len(inputTokens)) // Tandai token yang sudah digunakan

	// Cache TF-IDF untuk pertanyaan dalam KnowledgeBase
	questionTFIDFCache := make(map[string]map[string]float64)
	isSingleQuestion := len(kb.Questions) == 1
	idfAvailable := len(kb.IDF) > 0
	for _, question := range kb.Questions {
		questionTokens := tokenize(question.Question)
		questionTF := termFrequency(questionTokens)
		if isSingleQuestion || !idfAvailable {
			// Gunakan TF langsung jika hanya satu pertanyaan atau tidak ada data IDF
			questionTFIDFCache[question.Question] = questionTF
		} else {
			questionTFIDFCache[question.Question] = tfidfScore(questionTF, kb.IDF)
		}
	}

	start := 0
	for start < len(inputTokens) {
		var bestMatch Question
		highestSimilarity := 0.0
		bestMatchLength := 0

		// Tentukan panjang maksimum subTokens yang masuk akal
		maxLength := min(10, len(inputTokens)-start) // Misalnya, maksimal 10 kata

		for length := 1; length <= maxLength; length++ {
			end := start + length

			// Lewati jika rentang token sudah digunakan
			if isUsedRange(usedTokens, start, end) {
				continue
			}

			subTokens := inputTokens[start:end]
			subTF := termFrequency(subTokens)
			var subTFIDF map[string]float64
			if isSingleQuestion || !idfAvailable {
				subTFIDF = subTF
			} else {
				subTFIDF = tfidfScore(subTF, kb.IDF)
			}

			for _, question := range kb.Questions {
				questionTFIDF := questionTFIDFCache[question.Question]
				similarity := cosineSimilarity(subTFIDF, questionTFIDF)
				if similarity > highestSimilarity {
					highestSimilarity = similarity
					bestMatch = question
					bestMatchLength = length
				}
			}
		}

		if highestSimilarity > 0.1 {
			matches = append(matches, bestMatch)
			markUsedRange(usedTokens, start, start+bestMatchLength)
			start += bestMatchLength
		} else {
			start++
		}
	}

	return matches
}

// Cek apakah rentang token sudah digunakan
func isUsedRange(usedTokens []bool, start, end int) bool {
	for i := start; i < end; i++ {
		if usedTokens[i] {
			return true
		}
	}
	return false
}

// Tandai rentang token sebagai digunakan
func markUsedRange(usedTokens []bool, start, end int) {
	for i := start; i < end; i++ {
		usedTokens[i] = true
	}
}

// Hitung Term Frequency (TF)
func termFrequency(doc []string) map[string]float64 {
	tf := make(map[string]float64)
	for _, word := range doc {
		tf[word]++
	}
	for word := range tf {
		tf[word] = tf[word] / float64(len(doc))
	}
	return tf
}

// Hitung Inverse Document Frequency (IDF)
func inverseDocumentFrequency(corpus [][]string) map[string]float64 {
	idf := make(map[string]float64)
	totalDocs := float64(len(corpus))
	for _, doc := range corpus {
		seen := make(map[string]bool)
		for _, word := range doc {
			if !seen[word] {
				idf[word]++
				seen[word] = true
			}
		}
	}
	for word := range idf {
		idf[word] = math.Log(totalDocs / idf[word])
	}
	return idf
}

// Hitung skor TF-IDF
func tfidfScore(tf map[string]float64, idf map[string]float64) map[string]float64 {
	tfidf := make(map[string]float64)
	for word, tfValue := range tf {
		if idfValue, exists := idf[word]; exists {
			tfidf[word] = tfValue * idfValue
		}
	}
	return tfidf
}
