package scan_test

// import (
// 	"bytes"
// 	"go/token"
// 	"reflect"
// 	"testing"
// )

// // This file is not runnable
// // Just an example of how to make tests more readable

// // create type for tests
// type scanTest struct {
// 	name  string
// 	input string
// 	want  []token.Token
// }

// // separate out test closure function
// func (st scanTest) run(t *testing.T) {
// 	b := bytes.NewBufferString(st.input)
// 	s := NewScanner(ScanConfig{}, st.name, b)

// 	var got []token.Token
// 	for tok := s.Next(); tok.Type != token.EOF; tok = s.Next() {
// 		got = append(got, tok)
// 	}

// 	if !reflect.DeepEqual(st.want, got) {
// 		t.Errorf("line %q, wanted %v, got %v", st.input, st.want, got)
// 	}
// }

// // create table of test inputs
// var scanTests = []scanTest{
// 	{
// 		name:  "simple-add-comma",
// 		input: "2 1 +, 3+",
// 		want: []token.Token{
// 			{Type: token.Number, Line: 1, Text: "2"},
// 			{Type: token.Number, Line: 1, Text: "1"},
// 			{Type: token.Operator, Line: 1, Text: "+"},
// 			{Type: token.Comma, Line: 1, Text: ","},
// 			{Type: token.Number, Line: 2, Text: "3"},
// 			{Type: token.Operator, Line: 2, Text: "+"},
// 		},
// 	},
// 	...
// }

// // run sub-tests for items in table
// func TestScanner(t *testing.T) {
// 	for _, st := range scanTests {
// 		t.Run(st.name, st.run)
// 	}
// }

// MORE REFACTORING
// type checker interface {
// 	check(*testing.T, string, string) bool
// }

// type subTest struct {
// 	name string
// 	shouldFail bool
// 	checker checker <-- parameterize how we check results
// 	...
// }

// we can now define different checker types
// type checkGolden struct {...}

// func (c checkGolden) check(t *testing.T, got, want string) bool {
// 	...
// }
