package lex

import (
	"unicode"
	"unicode/utf8"
)

const (
	cr             = '\r'
	lf             = '\n'
	beginVCalender = "BEGIN:VCALENDAR"
	endVCalendar   = "END:VCALENDAR"
	beginVEvent    = "BEGIN:VEVENT"
	endVEvent      = "END:VEVENT"
)

// contentline   = name *(";" param ) ":" value CRLF
func lexContentLine(l *lexer) stateFunc {
	if l.hasPrefix(beginVCalender) {
		l.pos += len(beginVCalender)
		l.emit(CalendarBegin)
		return lexNewLine
	}

	if l.hasPrefix(endVCalendar) {
		l.pos += len(endVCalendar)
		l.emit(CalendarEnd)
		return lexNewLine
	}

	if l.hasPrefix(beginVEvent) {
		l.pos += len(beginVEvent)
		l.emit(EventBegin)
		return lexNewLine
	}

	if l.hasPrefix(endVEvent) {
		l.pos += len(endVEvent)
		l.emit(EventEnd)
		return lexNewLine
	}

	return lexName
}

func lexNewLine(l *lexer) stateFunc {
	r := l.next()
	if r == eof {
		l.emitEOF()
		return nil
	}

	if r == cr {
		r = l.next()
	} else if r == lf && l.strictLineBreaks {
		return l.errorf("missing carriage return (CR) at pos %d", l.pos)
	}

	if r != lf {
		return l.errorf("expected end of line at pos %d; got %s", l.pos, string(r))
	}

	if l.next() == eof {
		l.emitEOF()
		return nil
	}
	l.backup()
	l.ignore()

	return lexContentLine
}

// name          = iana-token / x-name
// iana-token    = 1*(ALPHA / DIGIT / "-")
// x-name        = "X-" [vendorid "-"] 1*(ALPHA / DIGIT / "-")
// vendorid      = 3*(ALPHA / DIGIT)
func lexName(l *lexer) stateFunc {
	for {
		r := l.next()
		if r == eof {
			return l.unexpectedEOF()
		}

		if isNameChar(r) {
			continue
		}

		l.backup()
		l.emitAdvanced(Name)

		r = l.next()

		switch r {
		case ':':
			l.ignore()
			return lexValue
		case ';':
			l.ignore()
			return lexParamName
		}

		return l.unexpected(r, ':', ';')
	}
}

// value         = *VALUE-CHAR
// VALUE-CHAR    = WSP / %x21-7E / NON-US-ASCII ; Any textual character
// NON-US-ASCII  = UTF8-2 / UTF8-3 / UTF8-4 ; UTF8-2, UTF8-3, and UTF8-4 are defined in [RFC3629]
// CONTROL       = %x00-08 / %x0A-1F / %x7F ; All the controls except HTAB
func lexValue(l *lexer) stateFunc {
	for {
		r := l.next()
		if r == eof {
			l.emitAdvanced(Value)
			return nil
		}

		if isValueChar(r) {
			continue
		}

		l.backup()
		l.emitAdvanced(Value)

		return lexNewLine
	}
}

// param         = param-name "=" param-value *("," param-value)
// param-name    = iana-token / x-name
// iana-token    = 1*(ALPHA / DIGIT / "-")
// x-name        = "X-" [vendorid "-"] 1*(ALPHA / DIGIT / "-")
// vendorid      = 3*(ALPHA / DIGIT)
func lexParamName(l *lexer) stateFunc {
	for {
		r := l.next()
		if r == eof {
			return l.unexpectedEOF()
		}

		if isNameChar(r) {
			continue
		}

		l.backup()
		l.emitAdvanced(ParamName)

		r = l.next()

		switch r {
		case '=':
			l.ignore()
			return lexParamValue
		case ':':
			l.ignore()
			return lexValue
		}

		return l.unexpected(r, '=', ':')
	}
}

// param-value   = paramtext / quoted-string
// paramtext     = *SAFE-CHAR
// quoted-string = DQUOTE *QSAFE-CHAR DQUOTE
// QSAFE-CHAR    = WSP / %x21 / %x23-7E / NON-US-ASCII ; Any character except CONTROL and DQUOTE
// SAFE-CHAR     = WSP / %x21 / %x23-2B / %x2D-39 / %x3C-7E / NON-US-ASCII ; Any character except CONTROL, DQUOTE, ";", ":", ","
// NON-US-ASCII  = UTF8-2 / UTF8-3 / UTF8-4 ; UTF8-2, UTF8-3, and UTF8-4 are defined in [RFC3629]
// CONTROL       = %x00-08 / %x0A-1F / %x7F ; All the controls except HTAB
func lexParamValue(l *lexer) stateFunc {
	for {
		r := l.next()
		if r == eof {
			l.emitAdvanced(ParamValue)
			return l.unexpectedEOF()
		}

		if isSafeChar(r) {
			continue
		}

		l.backup()
		l.emitAdvanced(ParamValue)

		r = l.next()

		switch r {
		case ':':
			l.ignore()
			return lexValue
		case ';':
			l.ignore()
			return lexParamName
		case ',':
			l.ignore()
			return lexParamValue
		}

		return l.unexpected(r, ':')
	}
}

// isNameChar checks if r is a unicode letter / digit or '-'
func isNameChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-'
}

// isQSafeChar check if r is a unicode letter / digit / '-' or '"'
func isQSafeChar(r rune) bool {
	return !unicode.IsControl(r) && r != '"'
}

// isSafeChar checks if r is a unicode letter / digit / '-' / '"' / ';' / ':' or ','
func isSafeChar(r rune) bool {
	return isQSafeChar(r) && r != ';' && r != ':' && r != ','
}

// isValueChar checks if r is a utf-8 control character or '\t'
func isValueChar(r rune) bool {
	return r == '\t' || (!unicode.IsControl(r) && utf8.ValidRune(r))
}
