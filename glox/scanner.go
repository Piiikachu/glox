package glox

import "strconv"

type Scanner struct {
	source  string
	start   int
	current int
	line    int
}

// var scanner = new(Scanner)

func (s *Scanner) init(source string) {
	s.source = source
	s.start = 0
	s.current = 0
	s.line = 1
}

func (s *Scanner) scanToken() Token {
	s.skipWhiteSpace()

	s.start = s.current
	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF, "", nil)
	}
	char := s.advance()
	switch char {
	case '(':
		return s.makeTokenByType(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeTokenByType(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeTokenByType(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeTokenByType(TOKEN_RIGHT_BRACE)
	case ',':
		return s.makeTokenByType(TOKEN_COMMA)
	case '.':
		return s.makeTokenByType(TOKEN_DOT)
	case '-':
		return s.makeTokenByType(TOKEN_MINUS)
	case '+':
		return s.makeTokenByType(TOKEN_PLUS)
	case ';':
		return s.makeTokenByType(TOKEN_SEMICOLON)
	case '*':
		return s.makeTokenByType(TOKEN_STAR)
	case '!':
		if s.match('=') {
			return s.makeTokenByType(TOKEN_BANG_EQUAL)
		} else {
			return s.makeTokenByType(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			return s.makeTokenByType(TOKEN_EQUAL_EQUAL)
		} else {
			return s.makeTokenByType(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			return s.makeTokenByType(TOKEN_LESS_EQUAL)
		} else {
			return s.makeTokenByType(TOKEN_LESS)
		}
	case '>':
		if s.match('=') {
			return s.makeTokenByType(TOKEN_GREATER_EQUAL)
		} else {
			return s.makeTokenByType(TOKEN_GREATER)
		}
	case '/':
		return s.makeTokenByType(TOKEN_SLASH)
	case '"':
		return s.stringToken()
	default:
		if isDigit(char) {
			return s.numerToken()
		} else if isAlpha(char) {
			return s.identifier()
		} else {
			return s.errorToken("Unexpected character.")
		}
	}
}

func (s *Scanner) numerToken() Token {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return s.errorToken(err.Error())
	}
	return s.makeTokenByLiteral(TOKEN_NUMBER, num)
}

func (s *Scanner) identifier() Token {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := Keywords[text]
	if !ok {
		tokenType = TOKEN_IDENTIFIER
	}
	return s.makeTokenByType(tokenType)
}

func isAlphaNumeric(char byte) bool {
	return isAlpha(char) || isDigit(char)
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) stringToken() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken("Unterminated string.")
	}

	s.advance()
	value := s.source[s.start+1 : s.current-1]
	return s.makeTokenByLiteral(TOKEN_STRING, value)
}

func (s *Scanner) skipWhiteSpace() {
	for {
		c := s.peek()
		switch c {
		case ' ':
			s.advance()
		case '\r':
			s.advance()
		case '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) makeTokenByType(t TokenType) Token {
	return s.makeTokenByLiteral(t, nil)
}

func (s *Scanner) makeTokenByLiteral(t TokenType, literal interface{}) Token {
	text := s.source[s.start:s.current]
	return s.makeToken(t, text, literal)
}

func (s *Scanner) makeToken(t TokenType, lexeme string, literal interface{}) Token {
	token := &Token{
		tokenType: t,
		lexeme:    lexeme,
		literal:   literal,
		line:      s.line,
	}
	return *token
}

func (s *Scanner) errorToken(msg string) Token {
	return s.makeToken(TOKEN_ERROR, msg, nil)
}
