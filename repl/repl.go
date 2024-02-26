package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"ronnie/evaluator"
	"ronnie/lexer"
	"ronnie/object"
	"ronnie/parser"
	"strings"
)

const (
	PROMPT       = ">> "
	QUIT_COMMAND = "quit"
	ERROR_FACE   = `
	⣿⣿⣿⣿⣿⣿⣿⡿⠛⠥⠖⠈⠩⠭⠭⠉⠉⠠⢤⣬⡙⢿⣿⣿⣿⣿⣿⣿⣿⣿
	⣿⣿⣿⣿⣿⣿⣿⢁⣶⠶⢦⡄⢸⣾⣿⡆⠀⣤⣶⢶⣬⡀⢿⣿⣿⣿⣿⣿⣿⣿
	⣿⣿⣿⣿⣿⣿⡏⠈⠷⠤⠼⢃⣨⣤⣤⣤⡈⠿⢅⣀⣿⡇⢸⣿⣿⣿⣿⣿⣿⣿
	⣿⣿⣿⣿⣿⣿⠇⣶⡶⢠⢾⣿⣿⣿⣿⣿⢻⡄⢰⣶⣶⣶⠘⣿⣿⣿⣿⣿⣿⣿
	⣿⣿⣿⣿⣿⡏⢰⣿⠁⣿⣾⣿⣿⣿⣿⣿⣾⣿⡆⢹⣿⣿⡄⢿⣿⣿⣿⣿⣿⣿
	⣿⣿⣿⣿⣿⡇⢸⡿⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠈⣿⣿⡿⠦⢬⣙⠻⣿⣿⣿
	⣿⣿⣿⣿⣿⣷⡈⣧⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⣿⡿⠁⠀⠀⠈⢷⡌⢿⣿
	⣿⣿⣿⣿⣿⣿⣷⢸⣧⡙⠿⣿⣿⣿⣿⣿⣿⡿⠟⢠⣿⠇⠀⠀⠀⠀⢸⣿⠘⣿
	⣿⣿⣿⣿⡿⠟⣋⣄⠻⣄⠀⢨⢙⡛⣩⢭⡤⡀⠀⣸⠏⣰⣆⡀⠀⠀⢸⣿⡆⣿
	⣿⣿⠟⣉⣤⣾⣿⣿⣷⣌⠑⠶⣦⣤⣥⣴⣶⡾⠟⢁⣼⣿⣿⣿⣶⣄⡀⣿⣧⠹
	⡿⢡⣾⣿⣿⡟⢛⠻⣿⣿⣷⣤⣄⠉⠉⠉⢁⣠⣴⣿⣿⣿⣿⣿⣿⣿⣿⣌⠻⡆
	⣱⣿⣿⣿⡿⢰⣿⣷⣌⡛⠿⣿⣿⣿⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡄
	⣿⣿⣿⡟⣰⣿⣿⣿⣿⣿⣷⣦⣤⣭⣉⣭⣭⣭⣤⣤⣶⣆⠻⠿⠿⣿⣿⣿⣿⣿
	⣿⣿⠟⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⢸⣿⣿⣿⣿⠀
	`
)

// Start the REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == QUIT_COMMAND {
			return
		}

		// Check if the input is a file with .ro extension
		if strings.HasSuffix(line, ".ro") {
			// Open the file
			file, err := os.Open(line)
			if err != nil {
				fmt.Fprintf(out, "Error opening file: %s\n", err)
				continue
			}
			defer file.Close()

			// Create a scanner for the file
			fileScanner := bufio.NewScanner(file)

			// Read and execute Ronnie code line by line
			for fileScanner.Scan() {
				line := fileScanner.Text()

				// Lex, parse, and evaluate Ronnie code
				evaluateInput(line, env, out)
			}

			if err := fileScanner.Err(); err != nil {
				fmt.Fprintf(out, "Error reading file: %s\n", err)
			}
		} else {
			// Lex, parse, and evaluate Ronnie code
			evaluateInput(line, env, out)
		}
	}
}

func evaluateInput(input string, env *object.Environment, out io.Writer) {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParseErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, ERROR_FACE+"\n")
	io.WriteString(out, "Parser error(s): ")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
