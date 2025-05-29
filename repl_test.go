package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " pikachu   bulbasaur  charmander",
			expected: []string{"pikachu", "bulbasaur", "charmander"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		expected := c.expected

		if len(actual) != len(expected) {
			t.Errorf("lengths do not match")
		}
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("words do not match")
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
