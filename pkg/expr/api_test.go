package expr

import (
	"testing"
)

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		expression string
		context    map[string]interface{}
		expected   bool
		err        error
	}{
		{
			expression: "",
			context:    map[string]interface{}{},
			expected:   false,
			err:        nil,
		},
		{
			expression: "true",
			context:    map[string]interface{}{},
			expected:   true,
			err:        nil,
		},
		{
			expression: "a > b",
			context:    map[string]interface{}{"a": 5, "b": 3},
			expected:   true,
			err:        nil,
		},
		{
			expression: "a > b",
			context:    map[string]interface{}{"a": 1, "b": 5},
			expected:   false,
			err:        nil,
		},
		{
			expression: `artifact_type == "report"`,
			context:    map[string]interface{}{"artifact_type": "report"},
			expected:   true,
			err:        nil,
		},
		// contains
		{
			expression: `contains(val, "a")`,
			context:    map[string]interface{}{"val": []string{"a", "b", "c"}},
			expected:   true,
			err:        nil,
		},
		{
			expression: `contains(val, "z")`,
			context:    map[string]interface{}{"val": []string{"a", "b", "c"}},
			expected:   false,
			err:        nil,
		},
		// containsKey
		{
			expression: `containsKey(val, "a")`,
			context:    map[string]interface{}{"val": map[string]string{"a": "b", "c": "d"}},
			expected:   true,
			err:        nil,
		},
		{
			expression: `containsKey(val, "z")`,
			context:    map[string]interface{}{"val": map[string]string{"a": "b", "c": "d"}},
			expected:   false,
			err:        nil,
		},
		// getMapValue
		{
			expression: `getMapValue(val, "a") == "b"`,
			context:    map[string]interface{}{"val": map[string]string{"a": "b", "c": "d"}},
			expected:   true,
			err:        nil,
		},
		{
			expression: `getMapValue(val, "a") == "z"`,
			context:    map[string]interface{}{"val": map[string]string{"a": "b", "c": "d"}},
			expected:   false,
			err:        nil,
		},
		// hasPrefix
		{
			expression: `hasPrefix("abc", "a")`,
			context:    nil,
			expected:   true,
			err:        nil,
		},
		{
			expression: `hasPrefix("abc", "z")`,
			context:    nil,
			expected:   false,
			err:        nil,
		},
		// inPath
		{
			expression: `inPath("ls")`,
			context:    nil,
			expected:   true,
			err:        nil,
		},
		{
			expression: `inPath("nonexistentcommand")`,
			context:    nil,
			expected:   false,
			err:        nil,
		},
		// regex
		{
			expression: `regex("abc", "[a-z]{3}")`,
			context:    nil,
			expected:   true,
			err:        nil,
		},
	}

	for _, test := range tests {
		result, err := EvalBooleanExpression(test.expression, test.context)

		if err != nil {
			t.Errorf("expected error: %v, but got: %v - expr: %s", test.err, err, test.expression)
		}

		if result != test.expected {
			t.Errorf("expected result: %v, but got: %v - expr: %s", test.expected, result, test.expression)
		}
	}
}
