package validation

import "testing"

type addTest struct {
	input  string
	result bool
}

func TestIsAlphaNumeric(t *testing.T) {
	var addTests = []addTest{
		{"Fjwro120j", true},
		{"nof1240bo", true},
		{"WWWW022", true},
		{"3r_14104bu", true},
		{"<sib>2uk", false},
		{"$!110ee", false},
		{"@testname", false},
		{"", false},
	}

	for _, test := range addTests {
		if got := IsAlphaNumeric(test.input); got != test.result {
			t.Errorf("input: %s, got %t, wanted %t", test.input, got, test.result)
		}
	}
}

func TestIsAlphabet(t *testing.T) {
	var addTests = []addTest{
		{"Fjwro120j", false},
		{"nof1240bo", false},
		{"flying", true},
		{"Elephant", true},
		{"letter", true},
		{"$!110ee", false},
		{"@testname", false},
		{"<script>fish</script>", false},
		{"my name", true},
		{"chua yi you", true},
		{"", false},
	}

	for _, test := range addTests {
		if got := IsAlphabet(test.input); got != test.result {
			t.Errorf("input: %s, got %t, wanted %t", test.input, got, test.result)
		}
	}
}

func TestIsTimeString(t *testing.T) {
	var addTests = []addTest{
		{"0000", true},
		{"2359", true},
		{"5000", false},
		{"Elephant", false},
		{"1200", true},
		{"1660", false},
		{"2400", false},
		{"250", false},
		{"", false},
	}

	for _, test := range addTests {
		if got := IsTimeString(test.input); got != test.result {
			t.Errorf("input: %s, got %t, wanted %t", test.input, got, test.result)
		}
	}
}

func TestIsDateString(t *testing.T) {
	var addTests = []addTest{
		{"19920801", true},
		{"20001231", true},
		{"17770011", false},
		{"19000132", false},
		{"20151200", false},
		{"00000801", true},
		{"150072", false},
		{"250", false},
		{"", false},
	}

	for _, test := range addTests {
		if got := IsDateString(test.input); got != test.result {
			t.Errorf("input: %s, got %t, wanted %t", test.input, got, test.result)
		}
	}
}
