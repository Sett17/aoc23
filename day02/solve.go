package day02

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bag struct {
	id       int
	minRed   int
	minGreen int
	minBlue  int
}

func (b Bag) isPossible(maxRed, maxGreen, maxBlue int) bool {
	return b.minRed <= maxRed && b.minGreen <= maxGreen && b.minBlue <= maxBlue
}

func (b Bag) Power() int {
	return b.minRed * b.minGreen * b.minBlue
}

func Solve() (string, error) {
	file, err := os.Open("input")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bags := make([]Bag, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		id, err := strconv.Atoi(line[5:strings.Index(line, ":")])
		if err != nil {
			return "", err
		}
		bag := Bag{id: id}
		remainder := line[strings.Index(line, ":")+2:]
		pulls := strings.Split(remainder, ";")
		for _, pull := range pulls {
			singleColors := strings.Split(pull, ",")
			for _, singleColor := range singleColors {
				singleColor = strings.TrimSpace(singleColor)
				number, err := strconv.Atoi(singleColor[0:strings.Index(singleColor, " ")])
				if err != nil {
					return "", err
				}
				color := singleColor[strings.Index(singleColor, " ")+1:]
				switch color {
				case "red":
					if bag.minRed == 0 || number > bag.minRed {
						bag.minRed = number
					}
				case "green":
					if bag.minGreen == 0 || number > bag.minGreen {
						bag.minGreen = number
					}
				case "blue":
					if bag.minBlue == 0 || number > bag.minBlue {
						bag.minBlue = number
					}
				}
			}
		}
		bags = append(bags, bag)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	bagIdSum := 0
	for _, bag := range bags {
		if bag.isPossible(12, 13, 14) {
			bagIdSum += bag.id
		}
	}

	powerSum := 0
	for _, bag := range bags {
		powerSum += bag.Power()
	}

	return fmt.Sprintf("%d", powerSum), nil
}
