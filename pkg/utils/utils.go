// https://github.com/surrealdb/surrealdb/blob/v2.0.4/core/src/sql/escape.rs

package utils

import (
	"strings"
)

const (
	BracketL       = "⟨"
	BracketR       = "⟩"
	BracketESC     = "\\⟩"
	Backtick       = "`"
	BacktickESC    = "\\`"
	SingleQuote    = "'"
	DoubleQuote    = `"`
	DoubleQuoteESC = `\"`
	Underscore     = '_'
)

func isAsciiDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isAsciiAlpha(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

func isAsciiAlphaNumeric(r rune) bool {
	return isAsciiDigit(r) || isAsciiAlpha(r)
}

func escape(str, left, right, escaped string) string {
	return left + strings.ReplaceAll(str, right, escaped) + right
}

func escapeNormal(str, left, right, escaped string) string {
	return escapeStartsNumeric(str, left, right, escaped)
}

func escapeStartsNumeric(str, left, right, escaped string) string {
	for i, r := range str {
		if i == 0 && isAsciiDigit(r) {
			return escape(str, left, right, escaped)
		}
		if !(isAsciiAlphaNumeric(r) || r == Underscore) {
			return escape(str, left, right, escaped)
		}
	}
	return str
}

func escapeFullNumeric(str, left, right, escaped string) string {
	numeric := true

	for _, r := range str {
		if !(isAsciiAlphaNumeric(r) || r == Underscore) {
			return escape(str, left, right, escaped)
		}
		if numeric && !isAsciiDigit(r) {
			numeric = false
		}
	}

	if numeric {
		return escape(str, left, right, escaped)
	}
	return str
}

func QuoteStr(str string) string {
	if str == "" {
		return SingleQuote + SingleQuote
	}

	str = strings.ReplaceAll(str, "\\", "\\\\")

	if !strings.Contains(str, SingleQuote) {
		return SingleQuote + str + SingleQuote
	}

	return DoubleQuote + strings.ReplaceAll(str, DoubleQuote, DoubleQuoteESC) + DoubleQuote
}

func QuoteKey(key string) string {
	if key == "" {
		return DoubleQuote + DoubleQuote
	}
	return escapeNormal(key, DoubleQuote, DoubleQuote, DoubleQuoteESC)
}

func QuoteRID(rid string) string {
	if rid == "" {
		return BracketL + BracketR
	}
	return escapeFullNumeric(rid, BracketL, BracketR, BracketESC)
}

func QuoteIdent(ident string) string {
	if ident == "" {
		return Backtick + Backtick
	}
	return escapeStartsNumeric(ident, Backtick, Backtick, BacktickESC)
}
