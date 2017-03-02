package interpreter

import "testing"

func TestParser(t *testing.T) {
	expectedRes := map[string]int64{
		"(1+2)":                       3,
		"(1+-2)":                      -1,
		"(9223372036854775807+1)":     -9223372036854775808,
		"((9223372036854775807+1)*2)": 0,
		"((1 && 2) + (0 || 3))":       2,
		"!!123":                       1,
	}

	c := 0
	var res int64

	for key, val := range expectedRes {
		if res = InterpretString(key); res != val {
			t.Errorf("Test #%d failed. Source=\"%s\", expected=%d, got=%d", c, key, val, res)
		}
		c++
	}
}
