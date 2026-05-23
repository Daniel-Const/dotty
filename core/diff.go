package core

import (
	"bufio"
	"os"
)

type DiffLine struct {
	Kind string // "add", "remove", "same"
	Text string
}

func DiffFile(srcPath, destPath string) ([]DiffLine, error) {
	srcLines, err := readLines(srcPath)
	if err != nil {
		return nil, err
	}
	destLines, err := readLines(destPath)
	if err != nil {
		return nil, err
	}
	return diffLines(srcLines, destLines), nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func diffLines(src, dest []string) []DiffLine {
	m, n := len(src), len(dest)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if src[i-1] == dest[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	var result []DiffLine
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && src[i-1] == dest[j-1] {
			result = append([]DiffLine{{Kind: "same", Text: src[i-1]}}, result...)
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			result = append([]DiffLine{{Kind: "add", Text: dest[j-1]}}, result...)
			j--
		} else {
			result = append([]DiffLine{{Kind: "remove", Text: src[i-1]}}, result...)
			i--
		}
	}
	return result
}
