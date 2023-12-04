#!/bin/bash

for i in {1..25}
do
    dir_name="day$(printf "%02d" $i)" # Formats the day as day01, day02, etc.
    file_name="$dir_name/solve.go"

    # Check if directory exists, if not, create it
    if [ ! -d "$dir_name" ]; then
        mkdir "$dir_name"
        echo "Created directory: $dir_name"
    fi

    # Check if solve.go exists in the directory, if not, create it
    if [ ! -f "$file_name" ]; then
        cat > "$file_name" <<- EOM
package $dir_name

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
EOM
        echo "Created file: $file_name"
    fi
done

echo "Setup complete."
