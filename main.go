package main

import (
	"aoc23/day01"
	"aoc23/day02"
	"aoc23/day03"
	"aoc23/day04"
	"aoc23/day05"
	"aoc23/day06"
	"aoc23/day07"
	"aoc23/day08"
	"aoc23/day09"
	"aoc23/day10"
	"aoc23/day11"
	"aoc23/day12"
	"aoc23/day13"
	"aoc23/day14"
	"aoc23/day15"
	"aoc23/day16"
	"aoc23/day17"
	"aoc23/day18"
	"aoc23/day19"
	"aoc23/day20"
	"aoc23/day21"
	"aoc23/day22"
	"aoc23/day23"
	"aoc23/day24"
	"aoc23/day25"
	"fmt"
	"os"
)

// Function type for day solutions
type SolveFunc func() (string, error)

// Map of day functions
var dayFunctions = map[string]SolveFunc{
	"1":  day01.Solve,
	"2":  day02.Solve,
	"3":  day03.Solve,
	"4":  day04.Solve,
	"5":  day05.Solve,
	"6":  day06.Solve,
	"7":  day07.Solve,
	"8":  day08.Solve,
	"9":  day09.Solve,
	"10": day10.Solve,
	"11": day11.Solve,
	"12": day12.Solve,
	"13": day13.Solve,
	"14": day14.Solve,
	"15": day15.Solve,
	"16": day16.Solve,
	"17": day17.Solve,
	"18": day18.Solve,
	"19": day19.Solve,
	"20": day20.Solve,
	"21": day21.Solve,
	"22": day22.Solve,
	"23": day23.Solve,
	"24": day24.Solve,
	"25": day25.Solve,
}

func main() {
	println("=== Aoc 23 Codeframe ===")
	args := os.Args[1:]
	if len(args) == 0 {
		println("No day specified :(")
		return
	}
	day := args[0]
	fmt.Println("Day:", day)

	if solve, ok := dayFunctions[day]; ok {
		result, err := solve()
		if err != nil {
			fmt.Printf("Error solving day %s: %s\n", day, err)
			return
		}
		fmt.Println("Result:", result)
	} else {
		fmt.Printf("No solution available for day %s\n", day)
	}
}
