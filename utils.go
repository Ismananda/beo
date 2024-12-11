package beo

import (
	"math"
	"math/rand"
	"regexp"
	"strings"
)

// Mengecek apakah sebuah item terkandung dalam slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Hitung cosine similarity
func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	dotProduct, magnitude1, magnitude2 := 0.0, 0.0, 0.0
	for word, v1 := range vec1 {
		if v2, exists := vec2[word]; exists {
			dotProduct += v1 * v2
		}
		magnitude1 += v1 * v1
	}
	for _, v2 := range vec2 {
		magnitude2 += v2 * v2
	}
	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

// levenshtein menghitung jarak Levenshtein antara dua string
// Jarak Levenshtein adalah jumlah operasi penyuntingan minimum (penyisipan, penghapusan, atau penggantian) yang diperlukan untuk mengubah satu string menjadi string lainnya.
func levenshtein(a, b string) int {
	// n dan m adalah panjang dari string a dan b
	n, m := len(a), len(b)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	// Matriks d digunakan untuk menyimpan nilai jarak sementara antara sub-string a dan b
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
	}

	// Inisialisasi baris pertama dan kolom pertama dengan angka urutan (0, 1, 2, ..., n atau m)
	// Ini mewakili jarak antara string kosong dan string lainnya.
	for i := 0; i <= n; i++ {
		d[i][0] = i
	}
	for j := 0; j <= m; j++ {
		d[0][j] = j
	}

	// Menghitung jarak Levenshtein untuk setiap pasangan substring a[0..i-1] dan b[0..j-1]
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			// Jika karakter a[i-1] dan b[j-1] berbeda, maka ada biaya 1 untuk mengganti karakter
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			// Pilih nilai minimum antara penghapusan, penyisipan, atau penggantian karakter
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}

	// Mengembalikan jarak Levenshtein antara string a dan b
	return d[n][m]
}

// min mengembalikan nilai terkecil
func min(input ...int) int {
	if len(input) < 1 {
		return 0
	}

	min := input[0]
	for _, v := range input {
		if min > v {
			min = v
		}
	}

	return min
}

// randomChoice memilih salah satu jawaban secara acak dari daftar pilihan
func randomChoice(choices []string) string {
	choiceLength := len(choices)
	if choiceLength == 0 {
		return ""
	}
	return choices[rand.Intn(choiceLength)]
}

// Pisah input berdasarkan tanda baca
func splitByPunctuation(input string) []string {
	re := regexp.MustCompile(`[.?!\n]+`) // Pencocokan tanda baca
	segments := re.Split(input, -1)
	cleanedSegments := []string{}
	for _, segment := range segments {
		segment = strings.TrimSpace(segment) // Hilangkan spasi kosong
		if segment != "" {
			cleanedSegments = append(cleanedSegments, segment)
		}
	}
	return cleanedSegments
}

// tokenize mengubah
func tokenize(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}
