package parser

import (
	"github.com/madlitz/go-dsl"
)

var recovering bool

func Parse(p *dsl.Parser) (dsl.AST, []dsl.Error) {
	skipNewLines(p)

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"PROGRAM", parseProgram},
		}})

	skipNewLines(p)

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"EOF", nil},
		}})

	return p.Exit()
}

func parseProgram(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("PROGRAM")

	// Expect program name directly after PROGRAM keyword
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", nil},
		}})
	p.AddTokens()

	// Ignore newlines
	skipNewLines(p)

	// Expect VAR keyword blocks or assignment/call
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"NL", skipNewLine},
			{"VAR", parseVarBlock},
			{"VARIABLE", assignmentOrCall},
			{"IF", parseIfStatement},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"END_PROGRAM", nil},
		}})
	p.SkipToken()
}

func parseVarBlock(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("VAR_BLOCK")

	skipNewLines(p)

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", parseVar},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})

	skipNewLines(p)

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"END_VAR", nil},
		}})
	p.SkipToken()
	p.WalkUp()

}

func parseVar(p *dsl.Parser) {
	p.AddNode("VAR")

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"COLON", nil}}})
	p.SkipToken()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"BOOL", nil},
			{"INT", nil},
			{"REAL", nil},
			{"STRING", nil}}})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"SEMICOLON", nil}}})
	p.SkipToken()

	skipNewLines(p)

	p.AddTokens()
	p.WalkUp()

}

// Function to handle assignments or calls
func assignmentOrCall(p *dsl.Parser) {
	p.AddNode("ASSIGNMENT")
	p.AddTokens() // The VARIABLE token has already been consumed

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"ASSIGN", nil},
		}})
	p.AddTokens()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"LITERAL", nil},
			{"VARIABLE", nil},
		}})
	p.AddTokens()

	// Add expression node for arithmetic operations
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"PLUS", expression},
			{"MINUS", expression},
			{"MULTIPLY", expression},
			{"DIVIDE", expression},
		},
		Options: dsl.ParseOptions{Optional: true}})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"SEMICOLON", nil},
		}})
	p.SkipToken()
	p.WalkUp()

	skipNewLines(p)
}

// Expression handler for arithmetic operations
func expression(p *dsl.Parser) {
	p.AddNode("EXPRESSION")
	p.AddTokens()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"LITERAL", term},
			{"VARIABLE", term},
		}})
	p.WalkUp()
}

func term(p *dsl.Parser) {
	p.AddNode("TERMINAL")
	p.AddTokens()
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment
// parse -> assignmentOrCall -> assignment
func assignment(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("ASSIGNMENT")
	p.AddTokens()

	skipNewLines(p)

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", operator},
			{"LITERAL", operator},
			{"OPEN_PAREN", parenExpression},
		},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"SEMICOLON", nil}, // Handle semicolon after assignment
			{"NL", skipNewLines},
			{"EOF", nil},
		},
	})
	p.SkipToken()
}

// parse -> assignmentOrCall -> call
func call(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("CALL")
	p.AddTokens()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", operator},
			{"LITERAL", operator},
			{"OPEN_PAREN", parenExpression},
			{"CLOSE_PAREN", operator}}})
}

// parse -> assignmentOrCall -> assignment -> [expression, operator]
// parse -> assignmentOrCall -> call -> [expression, operator]
func operator(p *dsl.Parser) {
	p.AddNode("TERMINAL")
	p.AddTokens()
	p.WalkUp()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"SEMICOLON", nil}, // Allow for semicolon as statement end
			{"PLUS", expression},
			{"MINUS", expression},
			{"DIVIDE", expression},
			{"MULTIPLY", expression},
			{"OPEN_PAREN", parenExpression}},
		Options: dsl.ParseOptions{Optional: true}})
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment -> [expression, operator] -> parenExpression
// parse -> assignmentOrCall -> assignment -> parenExpression
// parse -> assignmentOrCall -> call -> [expression, operator] -> parenExpression
// parse -> assignmentOrCall -> call -> parenExpression
func parenExpression(p *dsl.Parser) {
	p.Peek([]dsl.PeekToken{
		{[]string{}, expression}})
	skipNewLines(p)
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"CLOSE_PAREN", operator}}})
	p.Recover(skipUntilLineBreak)
}

// parse -> assignmentOrCall -> [expression] -> addcomment
func addcomment(p *dsl.Parser) {
	p.AddNode("COMMENT")
	p.AddTokens()
	p.WalkUp()
}

func skipUntilLineBreak(p *dsl.Parser) {
	recovering = true
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"UNKNOWN", nil}}})
	recovering = false
}

// Function to handle IF statements
func parseIfStatement(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("IF_STATEMENT")
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", nil},
		}})
	p.AddTokens()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"THEN", nil},
		}})
	p.SkipToken()

	// Parse statements inside THEN block
	skipNewLines(p)
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", assignmentOrCall},
			{"IF", parseIfStatement},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})

	// Optional ELSE block
	skipNewLines(p)
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"ELSE", nil},
		},
		Options: dsl.ParseOptions{Optional: true}})
	p.SkipToken()

	// Parse statements inside ELSE block
	skipNewLines(p)
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", assignmentOrCall},
			{"IF", parseIfStatement},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})

	// Skip newlines before END_IF
	skipNewLines(p)

	// End IF statement
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"END_IF", nil},
		}})
	p.SkipToken()
	p.WalkUp() // Exit IF_STATEMENT context
}

// Function to handle THEN part of IF statement
func parseThenStatement(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("THEN_STATEMENT")

	// Ignore newlines
	skipNewLines(p)

	// Expect assignments or if/else statements
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", assignmentOrCall},
			{"IF", parseIfStatement},
			{"ELSE", parseElseStatement},
			{"END_IF", nil}, // Directly handle END_IF here
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})
	p.SkipToken() // Skip END_IF token
	p.WalkUp()    // Exit IF_STATEMENT context
}

// Function to handle ELSE part of IF statement
func parseElseStatement(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("ELSE_STATEMENT")

	// Ignore newlines
	skipNewLines(p)

	// Expect assignments or if statements
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", assignmentOrCall},
			{"IF", parseIfStatement},
			{"END_IF", nil},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})
}

func skipNewLines(p *dsl.Parser) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"NL", skipNewLine},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})
}

func skipNewLine(p *dsl.Parser) {
	p.SkipToken()
}
