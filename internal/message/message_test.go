package message

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		input   string
		results []interface{}
	}{
		{
			"[SENT]hello\n",
			[]interface{}{
				ClientMessage{"hello"},
			},
		},
		{
			"[RECEIVED]pwet:hello\n[RECEIVED]plop:world\n",
			[]interface{}{
				ServerMessage{"pwet", "hello"},
				ServerMessage{"plop", "world"},
			},
		},
	}

	for _, test := range tests {
		reader := NewReader(strings.NewReader(test.input))
		results, err := reader.ReadAll()

		t.Log(results)

		if err != nil {
			t.Errorf("Unable to read command, error %v", err)
		} else if !reflect.DeepEqual(results, test.results) {
			t.Errorf("Expected: %v. Got: %v", test.results, results)
		}
	}
}

func TestWriter(t *testing.T) {
	tests := []struct {
		messages []interface{}
		result string
	}{
		{
			[]interface{}{
				ServerMessage{"pwet", "hello"},
			},
			"[RECEIVED]pwet:hello\n",
		},
	}

	buf := new(bytes.Buffer)

	for _, test := range tests {
		buf.Reset()
		cmdWriter := NewWriter(buf)

		for _, cmd := range test.messages {
			if cmdWriter.Write(cmd) != nil {
				t.Errorf("Unable to write message %v", cmd)
			}
		}

		t.Log(buf.String())

		if buf.String() != test.result {
			t.Errorf("Expected %v. Got %v", test.result, buf.String())
		}
	}

}
