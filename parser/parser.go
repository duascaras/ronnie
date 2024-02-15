package parser

import (
	"ronnie/ast"
	"ronnie/lexer"
	"ronnie/token"
)

type Parser struct {
	l *lexer.Lexer

	// Instead of pointing to the cur/next character, 
	// they will point to the cur/next TOKEN.
	// This happens in the nextToken() helper function bellow.
	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
