package main

import (
    "fmt"
    "os"
    "bufio"
    "io"
    "glox/util"
    "glox/scanner"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

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

func runFile(path string) {
    data, err := os.ReadFile(path)
    check(err)
    run(string(data))
    if util.HadError {
        os.Exit(65)
    }
}

func runPrompt() {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            check(err)
        } 
        run(line)
        util.HadError = false
    }
} 

func run(src string) {
    scan := scanner.NewScanner(src)
    tokens := scan.ScanTokens()
    for _, token := range tokens {
        fmt.Println(token)
    }
}

