package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    excludes := os.Args[1:]
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        excluded := false
        for _, pattern := range excludes {
            if strings.Contains(line, pattern) {
                excluded = true
                break
            }
        }
        if !excluded {
            fmt.Println(line)
        }
    }
}