package parser_test

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlitz/go-dsl"
	. "github.com/madlitz/inrush-core/parser"
)

func TestDSL(t *testing.T) {

	reader := bytes.NewBufferString(`
PROGRAM TestProgram
	VAR
		a : BOOL;
		b : INT;
	END_VAR
	a := TRUE;
	b := 5;
	IF a THEN
		b := b + 1;
	END_IF
END_PROGRAM
`)
	bufreader := bufio.NewReader(reader)
	ts := NewTokenSet()
	ns := NewNodeSet()
	logfilename := "TestDSL.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	ast, errs := dsl.ParseAndLog(Parse, Scan, ts, ns, bufreader, logfile)
	logfile.Close()
	if len(errs) != 0 {
		t.Fail()
		t.Error("Should report exactly 0 errors")
	}

	expectedNodes := []dsl.Node{
		{
			Type:   "PROGRAM",
			Tokens: []dsl.Token{{"PROGRAM", "TestProgram", 1, 1}},
			Children: []dsl.Node{
				{
					Type:   "VAR_BLOCK",
					Tokens: []dsl.Token{{"VAR", "", 2, 1}},
					Children: []dsl.Node{
						{
							Type:   "VAR",
							Tokens: []dsl.Token{{"VARIABLE", "a", 3, 3}},
							Children: []dsl.Node{
								{Type: "VAR_TYPE", Tokens: []dsl.Token{{"BOOL", "", 3, 6}}},
							},
						},
						{
							Type:   "VAR",
							Tokens: []dsl.Token{{"VARIABLE", "b", 4, 3}},
							Children: []dsl.Node{
								{Type: "VAR_TYPE", Tokens: []dsl.Token{{"INT", "", 4, 6}}},
							},
						},
					},
				},
				{
					Type:   "ASSIGNMENT",
					Tokens: []dsl.Token{{"VARIABLE", "a", 6, 1}, {"ASSIGN", ":=", 6, 3}, {"LITERAL", "TRUE", 6, 6}},
				},
				{
					Type:   "ASSIGNMENT",
					Tokens: []dsl.Token{{"VARIABLE", "b", 7, 1}, {"ASSIGN", ":=", 7, 3}, {"LITERAL", "5", 7, 6}},
				},
				{
					Type:   "IF_STATEMENT",
					Tokens: []dsl.Token{{"IF", "a", 8, 1}},
					Children: []dsl.Node{
						{
							Type:   "ASSIGNMENT",
							Tokens: []dsl.Token{{"VARIABLE", "b", 9, 5}, {"ASSIGN", ":=", 9, 7}, {"VARIABLE", "b", 9, 10}},
							Children: []dsl.Node{
								{
									Type:   "EXPRESSION",
									Tokens: []dsl.Token{{"PLUS", "+", 9, 12}},
									Children: []dsl.Node{
										{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "1", 9, 14}}},
									},
								},
							},
						},
					},
				},
				{
					Type:   "END_IF",
					Tokens: []dsl.Token{{"END_IF", "", 10, 1}},
				},
				{
					Type:   "END_PROGRAM",
					Tokens: []dsl.Token{{"END_PROGRAM", "", 11, 1}},
				},
			},
		},
	}

	if cmp.Diff(ast.RootNode.Children, expectedNodes) != "" {
		t.Fail()
		t.Error(cmp.Diff(ast.RootNode.Children, expectedNodes))
	}

	if errs != nil {
		t.Fail()
		for _, err := range errs {
			t.Error(err.String())
		}
	}
	//ast.Print()

}
