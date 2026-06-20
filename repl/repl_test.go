package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello    world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "my      NaMe    ",
			expected: []string{"my", "name"},
		},
	}

	for _, c := range cases {
		output := CleanInput(c.input)
		if len(output) == 0 {
			t.Errorf("output and test case won't match:\noutput: %v\ntestcase: %v\n\n", output, c.expected)
			t.Fail()
		}
		for i := range output {
			if c.expected[i] != output[i] {
				t.Errorf("output and test case won't match:\noutput: %v\ntestcase: %v\n\n", output, c.expected)
				t.Fail()
			}
		}
	}
}
