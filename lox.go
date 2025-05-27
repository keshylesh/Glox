package main

import (
    "fmt"
    "os"
    "bufio"
    "io"
    "glox/util"
    "glox/scanner"
    // "glox/ast"
    // "glox/token"
)

func main() {
    if len(os.Args) > 2 {
        fmt.Printf("Usage: %v <script>\n", os.Args[0])
        os.Exit(64)
    } else if len(os.Args) == 2 {
        runFile(os.Args[1])
    } else {
        runPrompt()
    }
}

// scan a file and interpret it
func runFile(path string) {
    data, err := os.ReadFile(path)
    util.Check(err)
    run(string(data))
    if util.HadError {
        os.Exit(65)
    }
}

// scan as a REPL and interpret line by line
func runPrompt() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("> ")
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        util.Check(err)
        run(line)
        util.HadError = false
    }
} 

// scan a line received from runPrompt() or runFile()
func run(src string) {
    scan := scanner.NewScanner(src)
    tokens := scan.ScanTokens()
    for _, token := range tokens {
        fmt.Println(token)
    }
}

