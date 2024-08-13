package main

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResp(t *testing.T) {
	msg := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	correctValue := Value{
		typ: TYPE_ARRAY,
		array: []Value{
			{typ: TYPE_BULK, bulk: "SET"},
			{typ: TYPE_BULK, bulk: "foo"},
			{typ: TYPE_BULK, bulk: "bar"},
		},
	}

	reader := strings.NewReader(msg)
	resp := NewResp(reader)

	Value, err := resp.Read()
	// Test the decryption
	assert.Equal(t, correctValue, Value)

	if err != nil {
		slog.Error("Failed due to error in reading.", "Err", err)
		t.FailNow()
	}

	byteArr := Value.Marshal()
	// Test the encryption
	assert.Equal(t, string(byteArr), msg)
}
