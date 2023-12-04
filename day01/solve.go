package day01

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Solve() (string, error) {
	file, err := os.Open("input")
	if err != nil {
		return "", err
	}
	defer file.Close()

	total := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		digits := make([]int, 0)

		validDigits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
		digitsLookup := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

		for i := 0; i < len(line); i++ {
			for validDigit := range validDigits {
				if strings.HasPrefix(line[i:], validDigits[validDigit]) {
					digits = append(digits, digitsLookup[validDigit])
					break
				}
			}
		}

		// for _, c := range line {
		// 	if c >= '0' && c <= '9' {
		// 		digits = append(digits, int(c))
		// 	}
		// }

		firstDigit := digits[0]
		number := firstDigit
		if len(digits) > 1 {
			lastDigit := digits[len(digits)-1]
			number *= 10
			number += lastDigit
		} else {
			number *= 10
			number += firstDigit
		}
		println(int(number))
		total += int(number)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", total), nil
}
