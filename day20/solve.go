package day20

import (
    "bufio"
    "os"
)

func Solve() (string, error) {
    file, err := os.Open("input")
    if err != nil {
        return "", err
    }
    defer file.Close()



    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        _ = line
    }
    if err := scanner.Err(); err != nil {
        return "", err
    }



    return "result", nil
}
