package repl

import (
	"bufio"
	"fmt"
	"io"
	"ronnie/lexer"
	"ronnie/token"
)

const (
	PROMPT = ">> "
	QUIT_COMMAND = "quit"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == QUIT_COMMAND {
			return
		}

		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
