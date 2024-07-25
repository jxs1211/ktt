package print

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Simple string",
			input:    "Hello, World!",
			expected: "\"Hello, World!\"\n",
		},
		{
			name:     "Integer",
			input:    42,
			expected: "42\n",
		},
		{
			name:  "Simple struct",
			input: struct{ Name string }{"John Doe"},
			expected: `{
  "Name": "John Doe"
}
`,
		},
		{
			name:  "Nested struct",
			input: struct{ Person struct{ Name string } }{Person: struct{ Name string }{"Jane Doe"}},
			expected: `{
  "Person": {
    "Name": "Jane Doe"
  }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to capture output
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrettyPrint(tt.input)

			// Restore stdout
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			actual := buf.String()

			if actual != tt.expected {
				t.Errorf("PrettyPrint(%v) = %v, want %v", tt.input, actual, tt.expected)
			}
		})
	}
}
