package day04

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	id             int
	winningNumbers []int
	cardNumbers    []int
}

func (c *Card) points() int {
	points := 0
	//first match gets 1 points, each other match multiplies points by 2
	firstMatch := true
	for _, n := range c.cardNumbers {
		for _, wn := range c.winningNumbers {
			if n == wn {
				if firstMatch {
					points += 1
					firstMatch = false
				} else {
					points *= 2
				}
			}
		}
	}
	return points
}

func (c *Card) howManyWinningNumbers() int {
	count := 0
	for _, n := range c.cardNumbers {
		for _, wn := range c.winningNumbers {
			if n == wn {
				count += 1
			}
		}
	}
	return count
}

func Solve() (string, error) {
	file, err := os.Open("input")
	if err != nil {
		return "", err
	}
	defer file.Close()

	cards := make([]Card, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		idString := line[strings.Index(line, " ")+1 : strings.Index(line, ":")]
		idString = strings.Trim(idString, " ")
		id, err := strconv.Atoi(idString)
		if err != nil {
			return "", err
		}
		winnngString := line[strings.Index(line, ":")+2 : strings.Index(line, "|")]
		winnngString = strings.Trim(winnngString, " ")
		cardString := line[strings.Index(line, "|")+2:]
		cardString = strings.Trim(cardString, " ")
		winningNumbers := make([]int, 0)
		cardNumbers := make([]int, 0)
		for _, n := range strings.Split(winnngString, " ") {
			if n == "" || n == " " {
				continue
			}
			num, err := strconv.Atoi(n)
			if err != nil {
				return "", err
			}
			winningNumbers = append(winningNumbers, num)
		}
		for _, n := range strings.Split(cardString, " ") {
			if n == "" || n == " " {
				continue
			}
			num, err := strconv.Atoi(n)
			if err != nil {
				return "", err
			}
			cardNumbers = append(cardNumbers, num)
		}
		card := Card{id, winningNumbers, cardNumbers}
		fmt.Println(card)
		cards = append(cards, card)

	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	totalPoints := 0
	for _, c := range cards {
		totalPoints += c.points()
	}

	instancesPerCard := make(map[int]int)

	//implementing the second part of the challenge

	totalCards := 0
	for _, v := range instancesPerCard {
		totalCards += v
	}

	return fmt.Sprintf("Total points: %d, Total cards: %d", totalPoints, totalCards), nil
}
