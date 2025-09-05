package assembler

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenComment
	TokenLabel
	TokenMnemonic
	TokenDirective
	TokenImmediate
	TokenAbsolute
	TokenIdentifier
	TokenNumber
	TokenString
	TokenComma
	TokenHash
	TokenDollar
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type Lexer struct {
	source   string
	position int
	line     int
	column   int
	tokens   []Token
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
		column: 1,
		tokens: make([]Token, 0),
	}
}

func (l *Lexer) Tokenize() ([]Token, error) {
	for l.position < len(l.source) {
		if err := l.nextToken(); err != nil {
			return nil, err
		}
	}
	
	l.addToken(TokenEOF, "")
	return l.tokens, nil
}

func (l *Lexer) nextToken() error {
	char := l.currentChar()
	
	switch {
	case char == 0:
		return nil
		
	case char == '\n':
		l.addToken(TokenNewline, string(char))
		l.advance()
		l.line++
		l.column = 1
		
	case char == '\r':
		l.advance()
		
	case unicode.IsSpace(char):
		l.skipWhitespace()
		
	case char == ';':
		l.readComment()
		
	case char == '#':
		l.addToken(TokenHash, "#")
		l.advance()
		
	case char == '$':
		l.advance()
		if err := l.readHexNumber(); err != nil {
			return err
		}
		
	case char == ',':
		l.addToken(TokenComma, ",")
		l.advance()
		
	case char == '.':
		return l.readDirective()
		
	case unicode.IsLetter(char) || char == '_':
		return l.readIdentifier()
		
	case unicode.IsDigit(char):
		return l.readNumber()
		
	case char == '"':
		return l.readString()
		
	default:
		return fmt.Errorf("unexpected character '%c' at line %d, column %d", char, l.line, l.column)
	}
	
	return nil
}

func (l *Lexer) currentChar() rune {
	if l.position >= len(l.source) {
		return 0
	}
	return rune(l.source[l.position])
}

func (l *Lexer) advance() {
	if l.position < len(l.source) {
		l.position++
		l.column++
	}
}

func (l *Lexer) addToken(tokenType TokenType, value string) {
	l.tokens = append(l.tokens, Token{
		Type:   tokenType,
		Value:  value,
		Line:   l.line,
		Column: l.column - len(value),
	})
}

func (l *Lexer) skipWhitespace() {
	for l.position < len(l.source) && unicode.IsSpace(rune(l.source[l.position])) && rune(l.source[l.position]) != '\n' {
		l.advance()
	}
}

func (l *Lexer) readComment() {
	start := l.position
	for l.position < len(l.source) && rune(l.source[l.position]) != '\n' {
		l.advance()
	}
	
	comment := l.source[start:l.position]
	l.addToken(TokenComment, comment)
}

func (l *Lexer) readDirective() error {
	start := l.position
	l.advance()
	
	for l.position < len(l.source) && (unicode.IsLetter(rune(l.source[l.position])) || unicode.IsDigit(rune(l.source[l.position]))) {
		l.advance()
	}
	
	directive := l.source[start:l.position]
	l.addToken(TokenDirective, directive)
	return nil
}

func (l *Lexer) readIdentifier() error {
	start := l.position
	
	for l.position < len(l.source) {
		char := rune(l.source[l.position])
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			break
		}
		l.advance()
	}
	
	identifier := l.source[start:l.position]
	
	if l.position < len(l.source) && rune(l.source[l.position]) == ':' {
		l.advance()
		l.addToken(TokenLabel, identifier)
	} else if l.isMnemonic(identifier) {
		l.addToken(TokenMnemonic, strings.ToUpper(identifier))
	} else {
		l.addToken(TokenIdentifier, identifier)
	}
	
	return nil
}

func (l *Lexer) readNumber() error {
	start := l.position
	
	for l.position < len(l.source) {
		char := rune(l.source[l.position])
		if !unicode.IsDigit(char) {
			break
		}
		l.advance()
	}
	
	number := l.source[start:l.position]
	l.addToken(TokenNumber, number)
	return nil
}

func (l *Lexer) readHexNumber() error {
	start := l.position
	
	for l.position < len(l.source) {
		char := rune(l.source[l.position])
		if !unicode.IsDigit(char) && char != 'A' && char != 'B' && char != 'C' && char != 'D' && char != 'E' && char != 'F' && char != 'a' && char != 'b' && char != 'c' && char != 'd' && char != 'e' && char != 'f' {
			break
		}
		l.advance()
	}
	
	number := l.source[start:l.position]
	l.addToken(TokenAbsolute, "$"+number)
	return nil
}

func (l *Lexer) readString() error {
	start := l.position
	l.advance()
	
	for l.position < len(l.source) && rune(l.source[l.position]) != '"' {
		l.advance()
	}
	
	if l.position >= len(l.source) {
		return fmt.Errorf("unterminated string at line %d", l.line)
	}
	
	l.advance()
	str := l.source[start:l.position]
	l.addToken(TokenString, str)
	return nil
}

func (l *Lexer) isMnemonic(word string) bool {
	mnemonics := []string{
		"ADC", "AND", "ASL", "BCC", "BCS", "BEQ", "BIT", "BMI", "BNE", "BPL",
		"BRK", "BVC", "BVS", "CLC", "CLD", "CLI", "CLV", "CMP", "CPX", "CPY",
		"DEC", "DEX", "DEY", "EOR", "INC", "INX", "INY", "JMP", "JSR", "LDA",
		"LDX", "LDY", "LSR", "NOP", "ORA", "PHA", "PHP", "PLA", "PLP", "ROL",
		"ROR", "RTI", "RTS", "SBC", "SEC", "SED", "SEI", "STA", "STX", "STY",
		"TAX", "TAY", "TSX", "TXA", "TXS", "TYA",
	}
	
	upper := strings.ToUpper(word)
	for _, mnemonic := range mnemonics {
		if upper == mnemonic {
			return true
		}
	}
	return false
}