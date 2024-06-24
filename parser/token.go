package parser

import (
	"github.com/madlitz/go-dsl"
)

func NewTokenSet() dsl.TokenSet {
	return dsl.NewTokenSet(
		"LITERAL",
		"PLUS",
		"MINUS",
		"MULTIPLY",
		"DIVIDE",
		"OPEN_PAREN",
		"CLOSE_PAREN",
		"COLON",
		"SEMICOLON",
		"ASSIGN",
		"VARIABLE",
		"COMMENT",
		"NL",
		"EOF",
		// IEC 61131-3 specific tokens
		"PROGRAM",
		"END_PROGRAM",
		"VAR",
		"END_VAR",
		"BOOL",
		"INT",
		"REAL",
		"STRING",
		"IF",
		"THEN",
		"ELSE",
		"ELSIF",
		"END_IF",
		"AND",
		"OR",
		"NOT",
		"GT", // >
		"LT", // <
		"GE", // >=
		"LE", // <=
		"EQ", // =
		"NE", // <>
	)
}
