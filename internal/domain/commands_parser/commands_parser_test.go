package commandsparser

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	cp := NewCommandsParser()

	tests := []struct {
		input       string
		expectedKey string
		expectedCmd Command
		expectedErr bool
	}{
		// Positive test cases
		{"addItem('key1', 'value1')", "key1", Command{Type: AddItem, Key: "key1", Value: "value1"}, false},
		{"deleteItem('key1')", "key1", Command{Type: DeleteItem, Key: "key1"}, false},
		{"getItem('key1')", "key1", Command{Type: GetItem, Key: "key1"}, false},
		{"getAllItems()", "", Command{Type: GetAllItems}, false},

		// Negative test cases
		{"addItem('key1', )", "", Command{}, true},
		{"addItem('key1', 'value1'", "", Command{}, true},
		{"deleteItem(key1)", "", Command{}, true},
		{"getItem('key1'", "", Command{}, true},
		{"getAllItems(", "", Command{}, true},
		{"unknownCommand('key1')", "", Command{}, true},
		{"", "", Command{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, cmd, err := cp.parseCommand(tt.input)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if key != tt.expectedKey {
				t.Errorf("expected key: %s, got: %s", tt.expectedKey, key)
			}
			if cmd != tt.expectedCmd {
				t.Errorf("expected cmd: %+v, got: %+v", tt.expectedCmd, cmd)
			}
		})
	}
}

func TestGetCommandID(t *testing.T) {
	cp := NewCommandsParser()

	tests := []struct {
		input       string
		expectedKey string
		expectedErr bool
	}{
		// Positive test cases
		{"addItem('key1', 'value1')", "key1", false},
		{"deleteItem('key1')", "key1", false},
		{"getItem('key1')", "key1", false},
		{"getAllItems()", "", false},

		// Negative test cases
		{"addItem('key1', )", "", true},
		{"addItem('key1', 'value1'", "", true},
		{"deleteItem(key1)", "", true},
		{"getItem('key1'", "", true},
		{"getAllItems(", "", true},
		{"unknownCommand('key1')", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			key, err := cp.GetCommandID(tt.input)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if key != tt.expectedKey {
				t.Errorf("expected key: %s, got: %s", tt.expectedKey, key)
			}
		})
	}
}

func TestParseCommandOnly(t *testing.T) {
	cp := NewCommandsParser()

	tests := []struct {
		input       string
		expectedCmd Command
		expectedErr bool
	}{
		// Positive test cases
		{"addItem('key1', 'value1')", Command{Type: AddItem, Key: "key1", Value: "value1"}, false},
		{"deleteItem('key1')", Command{Type: DeleteItem, Key: "key1"}, false},
		{"getItem('key1')", Command{Type: GetItem, Key: "key1"}, false},
		{"getAllItems()", Command{Type: GetAllItems}, false},

		// Negative test cases
		{"addItem('key1', )", Command{}, true},
		{"addItem('key1', 'value1'", Command{}, true},
		{"deleteItem(key1)", Command{}, true},
		{"getItem('key1'", Command{}, true},
		{"getAllItems(", Command{}, true},
		{"unknownCommand('key1')", Command{}, true},
		{"", Command{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd, err := cp.ParseCommand(tt.input)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}
			if cmd != tt.expectedCmd {
				t.Errorf("expected cmd: %+v, got: %+v", tt.expectedCmd, cmd)
			}
		})
	}
}
