package main

import (
    "fmt"
    "os"
    "glox/util"
    "strings"
)


// Tool to be run separately from glox module
func main() {
    if len(os.Args) < 2 {
        fmt.Print(os.Stderr, "Usage: go run generate_ast <output dir>")
        os.Exit(64)
    }
    outputDir := os.Args[1]

    defineAst(outputDir, "Expr", map[string][]string {
        "Assign": {"Name Token", "Value Expr"},
        "Binary": {"Left Expr", "Operator Token", "Right Expr"},
        "Call": {"Callee Expr", "Paren Token", "Arguments []Expr"},
        "Grouping": {"Expression Expr"},
        "Literal": {"Value Object"},
        "Logical": {"Left Expr", "Operator Token", "Right Expr"},
        "Unary": {"Operator Token", "Right Expr"},
        "Variable": {"Name Token"},
    })
    
    defineAst(outputDir, "Stmt", map[string][]string {
        "Block": {"Statements []Stmt"},
        "StmtExpression": {"Expression Expr"},
        "Function": {"Name Token", "Params []Token", "Body []Stmt"},
        "If": {"Condition Expr", "ThenBranch Stmt", "ElseBranch Stmt"},
        "Print": {"Expression Expr"},
        "Var": {"Name Token", "Initializer Expr"},
        "While": {"Condition Expr", "Body Stmt"},
    })
}

// Function to create the Abstract Syntax Tree
func defineAst(outputDir, baseName string, rules map[string][]string) {
    err := os.Mkdir(outputDir, 0755)
    if !os.IsExist(err) {
        util.Check(err)
    }

    fp, err := os.Create(outputDir + "/" + baseName + ".go")
    util.Check(err)

    defer fp.Close()
    
    // package name
    fp.WriteString("package " + outputDir + "\n\n")

    // Some unfortunate hardcoding to satisfy the go gods
    if (baseName == "Expr") {
        // imports
        fp.WriteString("import (\n")
        fp.WriteString("\t. \"glox/token\"\n")
        fp.WriteString("\t. \"glox/util\"\n")
        fp.WriteString(")\n\n")

    } else if (baseName == "Stmt") {
        // imports
        fp.WriteString("import (\n")
        fp.WriteString("\t. \"glox/token\"\n")
        fp.WriteString("\t. \"glox/util\"\n")
        fp.WriteString(")\n\n")
    }
    // create visitor interface
    defineVisitor(fp, baseName, rules)


    // create the base class
    fp.WriteString("type " + baseName + " interface{\n")
    fp.WriteString("\tAccept(v " + baseName + "Visitor) (Object, error)\n")
    fp.WriteString("}\n\n")

    // create all the types
    for className, fields := range rules {
        defineType(fp, baseName, className, fields)
    }
}

// Function to create the visitor interface and all the functions to implement
func defineVisitor(fp *os.File, baseName string, rules map[string][]string) {
    fp.WriteString("type " + baseName + "Visitor interface {\n")
    
    for className, _ := range rules {
        fp.WriteString("\tVisit" + className + "(obj " + className + ") (Object, error)\n")
    }

    fp.WriteString("}\n\n")
}

// Function to create a new type for a specified class
func defineType(fp *os.File, baseName, className string, fields []string) {
    // write type
    fp.WriteString("type " + className + " struct {\n")

    for _, field := range fields {
        fp.WriteString("\t" + field + "\n")
    }

    fp.WriteString("}\n\n")
    
    // write "constructors"
    fp.WriteString("func New" + className + "(")

    for i, field := range fields {
        fp.WriteString(field)
        if i != len(fields) - 1 {
            fp.WriteString(", ")
        }
    }

    fp.WriteString(") " + className + " {\n")

    fp.WriteString("\treturn " + className + "{")

    for i, field := range fields {
        slice := strings.Split(field, " ")
        fp.WriteString(slice[0] + ",")
        if i != len(fields) - 1 {
            fp.WriteString(" ")
        }
    }

    fp.WriteString("}\n}\n\n")

    // write accept
    fp.WriteString("func (obj " + className + ") Accept(v " + baseName + "Visitor) (Object, error) {\n")
    fp.WriteString("\treturn v.Visit" + className + "(obj)\n")
    fp.WriteString("}\n\n")
}
