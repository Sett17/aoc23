package day03

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Space struct {
	line  int
	index int
	rune  rune
}

func (s Space) lookup(grid [][]rune) rune {
	return grid[s.line][s.index]
}

type Number struct {
	line   int
	index  int
	length int

	number int
}

// gear are all * that is adjaccent to exactly 3 numbers (ratio is * of the two)
type Gear struct {
	line  int
	index int
	Nums  [2]Number
}

func (g Gear) ratio() int {
	return g.Nums[0].number * g.Nums[1].number
}

func Solve() (string, error) {
	file, err := os.Open("input")
	if err != nil {
		return "", err
	}
	defer file.Close()

	grid := make([][]rune, 0)

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, make([]rune, len(line)))
		for i, r := range line {
			grid[lineCount][i] = r
		}
		lineCount++

	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	totalNumber := 0

	//generate numbers, can be anywhere from 1 to 3 digits, save number itself in struct also, should always havve maximum lengths, e.g. number is 215, hen 15 and 5 are not seperate numbers
	numbers := make([]Number, 0)
	for i, line := range grid {
		for j, r := range line {
			if r < '0' || r > '9' {
				continue
			}
			// Check if the previous character is a digit; if it is, continue to the next character
			if j > 0 && line[j-1] >= '0' && line[j-1] <= '9' {
				continue
			}
			startIndex := j
			numberString := ""
			for j < len(line) && line[j] >= '0' && line[j] <= '9' {
				numberString += string(line[j])
				j++
			}
			number := Number{
				line:   i,
				index:  startIndex,
				length: j - startIndex,
				number: stringToNumber(numberString),
			}
			numbers = append(numbers, number)
		}
	}

	for _, number := range numbers {
		adjacents := findAdjacents(grid, number)

		fmt.Printf("Number: %v,\tadjacents: ", number)
		for _, r := range adjacents {
			fmt.Printf("%c ", r.rune)
		}

		hasSymbol := false
		for _, r := range adjacents {
			if isSymbol(r.rune) {
				hasSymbol = true
				break
			}
		}
		fmt.Printf(",\thasSymbol: %v", hasSymbol)
		println()
		if hasSymbol {
			totalNumber += number.number
		}
	}

	println()

	//find gears
	gears := make([]Gear, 0)
	gearSpaces := make([]Space, 0)
	for i, line := range grid {
		for j, r := range line {
			if r == '*' {
				gearSpaces = append(gearSpaces, Space{line: i, index: j, rune: r})
			}
		}
	}
	ratiosSum := 0
	for _, gearSpace := range gearSpaces {
		multiplicants := make([]Number, 0)
		for _, number := range numbers {
			for _, adj := range findAdjacents(grid, number) {
				if adj.line == gearSpace.line && adj.index == gearSpace.index {
					multiplicants = append(multiplicants, number)
					if len(multiplicants) == 2 {
						gears = append(gears, Gear{line: gearSpace.line, index: gearSpace.index, Nums: [2]Number{multiplicants[0], multiplicants[1]}})
						ratiosSum += multiplicants[0].number * multiplicants[1].number
					}
				}
			}
		}
		if len(multiplicants) == 2 {
			gear := Gear{line: gearSpace.line, index: gearSpace.index, Nums: [2]Number{multiplicants[0], multiplicants[1]}}
			duplicate := false
			for _, g := range gears {
				if g.Nums[0].number == gear.Nums[0].number && g.Nums[1].number == gear.Nums[1].number {
					duplicate = true
					break
				}
			}
			if duplicate {
				continue
			}
			gears = append(gears, gear)
			ratiosSum += multiplicants[0].number * multiplicants[1].number
		}
	}

	for _, gear := range gears {
		fmt.Printf("Gear at %d,%d [%d,%d] with ratio %d\n", gear.line, gear.index, gear.Nums[0].number, gear.Nums[1].number, gear.ratio())
	}

	return fmt.Sprintf("%d", ratiosSum), nil
}

func stringToNumber(s string) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return ret
}

func printGrid(grid [][]rune) {
	for _, line := range grid {
		for _, r := range line {
			print(string(r))
		}
		println()
	}
}

// also diagonal
func findAdjacents(grid [][]rune, number Number) []Space {
	allAdjacents := make([]Space, 0)
	if number.length == 1 {
		//top
		if number.line-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index, rune: grid[number.line-1][number.index]})
		}
		//bottom
		if number.line+1 < len(grid) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index, rune: grid[number.line+1][number.index]})
		}
		//left
		if number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index - 1, rune: grid[number.line][number.index-1]})
		}
		//right
		if number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index + 1, rune: grid[number.line][number.index+1]})
		}
		//top left
		if number.line-1 >= 0 && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index - 1, rune: grid[number.line-1][number.index-1]})
		}
		//top right
		if number.line-1 >= 0 && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 1, rune: grid[number.line-1][number.index+1]})
		}
		//bottom left
		if number.line+1 < len(grid) && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index - 1, rune: grid[number.line+1][number.index-1]})
		}
		//bottom right
		if number.line+1 < len(grid) && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 1, rune: grid[number.line+1][number.index+1]})
		}
	} else if number.length == 2 {
		//top1
		if number.line-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index, rune: grid[number.line-1][number.index]})
		}
		//top2
		if number.line-1 >= 0 && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 1, rune: grid[number.line-1][number.index+1]})
		}
		//bottom1
		if number.line+1 < len(grid) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index, rune: grid[number.line+1][number.index]})
		}
		//bottom2
		if number.line+1 < len(grid) && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 1, rune: grid[number.line+1][number.index+1]})
		}
		//left
		if number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index - 1, rune: grid[number.line][number.index-1]})
		}
		//right
		if number.index+2 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index + 2, rune: grid[number.line][number.index+2]})
		}
		//top left
		if number.line-1 >= 0 && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index - 1, rune: grid[number.line-1][number.index-1]})
		}
		//top right
		if number.line-1 >= 0 && number.index+2 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 2, rune: grid[number.line-1][number.index+2]})
		}
		//bottom left
		if number.line+1 < len(grid) && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index - 1, rune: grid[number.line+1][number.index-1]})
		}
		//bottom right
		if number.line+1 < len(grid) && number.index+2 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 2, rune: grid[number.line+1][number.index+2]})
		}
	} else if number.length == 3 {
		//top1
		if number.line-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index, rune: grid[number.line-1][number.index]})
		}
		//top2
		if number.line-1 >= 0 && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 1, rune: grid[number.line-1][number.index+1]})
		}
		//top3
		if number.line-1 >= 0 && number.index+2 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 2, rune: grid[number.line-1][number.index+2]})
		}
		//bottom1
		if number.line+1 < len(grid) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index, rune: grid[number.line+1][number.index]})
		}
		//bottom2
		if number.line+1 < len(grid) && number.index+1 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 1, rune: grid[number.line+1][number.index+1]})
		}
		//bottom3
		if number.line+1 < len(grid) && number.index+2 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 2, rune: grid[number.line+1][number.index+2]})
		}
		//left
		if number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index - 1, rune: grid[number.line][number.index-1]})
		}
		//right
		if number.index+3 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line, index: number.index + 3, rune: grid[number.line][number.index+3]})
		}
		//top left
		if number.line-1 >= 0 && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index - 1, rune: grid[number.line-1][number.index-1]})
		}
		//top right
		if number.line-1 >= 0 && number.index+3 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line - 1, index: number.index + 3, rune: grid[number.line-1][number.index+3]})
		}
		//bottom left
		if number.line+1 < len(grid) && number.index-1 >= 0 {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index - 1, rune: grid[number.line+1][number.index-1]})
		}
		//bottom right
		if number.line+1 < len(grid) && number.index+3 < len(grid[number.line]) {
			allAdjacents = append(allAdjacents, Space{line: number.line + 1, index: number.index + 3, rune: grid[number.line+1][number.index+3]})
		}
	}
	return allAdjacents
}

func isSymbol(r rune) bool {
	return (r < '0' || r > '9') && r != '.'
}
