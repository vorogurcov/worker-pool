package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	keysFile := "data/keys.txt"
	baseURL := "http://localhost:30099/v1/get?key="

	file, err := os.Open(keysFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var keys []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := strings.TrimSpace(scanner.Text())
		if key != "" {
			keys = append(keys, key)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR scanning file: %v\n", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	out := bufio.NewWriter(os.Stdout)
	count := 0
	for _, key := range keys {
		fmt.Fprintf(out, "GET %s%s\n", baseURL, key)
		count++

		if count%100 == 0 {
			out.Flush()
		}
	}
	out.Flush()

	fmt.Fprintf(os.Stderr, "\n[DONE] Всего ключей отправлено в очередь: %d\n", count)
}
