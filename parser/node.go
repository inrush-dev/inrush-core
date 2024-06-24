package parser

import (
	"github.com/madlitz/go-dsl"
)

func NewNodeSet() dsl.NodeSet {
	return dsl.NewNodeSet(
		"COMMENT",
		"EXPRESSION",
		"ASSIGNMENT",
		"TERMINAL",
		"CALL",
		// IEC 61131-3 specific nodes
		"PROGRAM",
		"END_PROGRAM",
		"VAR_BLOCK",
		"END_VAR",
		"VAR",
		"VAR_TYPE",
		"VAR_VALUE",
		"IF_STATEMENT",
		"THEN_STATEMENT",
		"ELSE_STATEMENT",
		"END_IF",
	)
}
