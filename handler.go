package main

import (
	"log/slog"
	"sync"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

// Handlers stores the mapping of the redis command strings to the parsing functions
var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

// ping is a function that returns the response Value for a ping message
func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "STRING", str: "PONG"}
	}
	return Value{typ: "STRING", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "ERROR", str: "ERROR Wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "STRING", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "ERROR", str: "ERROR: Wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	valueStr, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "NULL"}
	}

	return Value{typ: "STRING", str: valueStr}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "ERROR", str: "ERROR: Wrong number of arguments for 'hset' command"}
	}

	key := args[0].bulk
	hsetKey := args[1].bulk
	hsetVal := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[key]; !ok {
		HSETs[key] = map[string]string{}
	}
	HSETs[key][hsetKey] = hsetVal
	HSETsMu.Unlock()

	return Value{typ: "STRING", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "ERROR", str: "ERROR: Wrong number of arguments for 'hget' command"}
	}

	key := args[0].bulk
	hsetKey := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[key][hsetKey]
	HSETsMu.RUnlock()

	if !ok {
		slog.Error("value: %v, ok: %v", value, ok)
		return Value{typ: "NULL"}
	}

	return Value{typ: "BULK", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "ERROR", str: "ERROR: Wrong number of arguments for 'hgetall' command"}
	}

	value := Value{typ: "ARRAY", array: []Value{}}

	key := args[0].bulk
	HSETsMu.RLock()
	hashSet := HSETs[key]
	HSETsMu.RUnlock()

	for hsetKey, hsetValue := range hashSet {
		value.array = append(value.array, Value{typ: "BULK", bulk: hsetKey})
		value.array = append(value.array, Value{typ: "BULK", bulk: hsetValue})
	}

	return value
}
