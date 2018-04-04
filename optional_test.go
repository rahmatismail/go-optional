package optional_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/rahmatismail/goptional"
)

type TestA struct {
	A optional.Int64   `json:"a"`
	B []optional.Int64 `json:"b,omitempty"`
}

type TestB struct {
	A optional.String   `json:"a"`
	B []optional.String `json:"b,omitempty"`
}

func TestInt64(t *testing.T) {
	// Scenario: Optional defaults to invalid
	// Given: Some optional value
	for i := int64(0); i < int64(20); i++ {
		v := optional.Int64{}
		v.Set(i, false)
		// When: Checked for validity using Ok
		// Then: Returns false
		if v.Ok() {
			t.Errorf("[optional] Invalid default value, expected: false, got: true")
		}
		// When: Checked for validity using Get
		// Then: Returns stored value and true
		if val, ok := v.Get(); ok || val != i {
			t.Errorf("[optional] Invalid default value, expected: %v and false, got: %v and %v", i, val, ok)
		}
	}

	// Scenario: Optional is invalid if unmarshal found no such key
	// Given: JSON message with missing key
	testMsgA := `{"b": [null]}`
	ta := TestA{}
	if err := json.Unmarshal([]byte(testMsgA), &ta); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgA)
	}
	if _, ok := ta.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}
	for _, k := range ta.B {
		if _, ok := k.Get(); ok {
			t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
		}
	}

	// Scenario: Optional is invalid if unmarshal type is different than expected
	// Given: JSON message with wrong value type
	testMsgB := `{"a": "test"}`
	tb := TestA{}
	if err := json.Unmarshal([]byte(testMsgB), &tb); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgB)
	}
	if _, ok := tb.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}

	// Scenario: Marshal/unmarshal JSON <==> optional
	// Given: Some JSON and corresponding structure
	testCaseA := []struct {
		message       []byte
		expectedA     bool
		expectedBSize int
	}{
		{[]byte(`{}`), false, 0},
		{[]byte(`{"a": 4}`), true, 0},
		{[]byte(`{"a": 4, "b": []}`), true, 0},
		{[]byte(`{"a": 4, "b": [1, 2, 3]}`), true, 3},
		{[]byte(`{"b": [1, 2, 3]}`), false, 3},
	}

	for _, v := range testCaseA {
		// When: Unmarshaled using JSON.Unmarshal
		// Then: Returns structure with value as expected
		k := TestA{}
		if err := json.Unmarshal(v.message, &k); err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		if k.A.Ok() != v.expectedA {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v",
				v.expectedA, k.A)
		}
		if len(k.B) != v.expectedBSize {
			t.Errorf("[optional] Unexpected array value, expected: %v, got: %v",
				v.expectedBSize, k.B)
		}

		// When: Marshaled using JSON.Marshal
		// Then: Returns same message
		b, err := json.Marshal(k)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		b = SimplifyLiteral(b)
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		_, ok := m["a"]
		if v.expectedA != ok {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
		kb, _ := m["b"].([]interface{})
		if len(kb) != v.expectedBSize {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
	}
}

func TestString(t *testing.T) {
	// Scenario: Optional defaults to invalid
	// Given: Some optional value
	for i := int64(0); i < int64(20); i++ {
		testVal := fmt.Sprintln(i)
		v := optional.String{}
		v.Set(testVal, false)
		// When: Checked for validity using Ok
		// Then: Returns false
		if v.Ok() {
			t.Errorf("[optional] Invalid default value, expected: false, got: true")
		}
		// When: Checked for validity using Get
		// Then: Returns stored value and true
		if val, ok := v.Get(); ok || val != testVal {
			t.Errorf("[optional] Invalid default value, expected: %v and false, got: %v and %v", i, val, ok)
		}
	}

	// Scenario: Optional is invalid if unmarshal found no such key
	// Given: JSON message with missing key
	testMsgA := `{"b": [null]}`
	ta := TestB{}
	if err := json.Unmarshal([]byte(testMsgA), &ta); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgA)
	}
	if _, ok := ta.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}
	for _, k := range ta.B {
		if _, ok := k.Get(); ok {
			t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
		}
	}

	// Scenario: Optional is invalid if unmarshal type is different than expected
	// Given: JSON message with wrong value type
	testMsgB := `{"a": 0}`
	tb := TestB{}
	if err := json.Unmarshal([]byte(testMsgB), &tb); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgB)
	}
	if _, ok := tb.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}

	// Scenario: Marshal/unmarshal JSON <==> optional
	// Given: Some JSON and corresponding structure
	testCaseA := []struct {
		message       []byte
		expectedA     bool
		expectedBSize int
	}{
		{[]byte(`{}`), false, 0},
		{[]byte(`{"a": "4"}`), true, 0},
		{[]byte(`{"a": "4", "b": []}`), true, 0},
		{[]byte(`{"a": "4", "b": ["1", "2", "3"]}`), true, 3},
		{[]byte(`{"b": ["1", "2", "3"]}`), false, 3},
	}

	for _, v := range testCaseA {
		// When: Unmarshaled using JSON.Unmarshal
		// Then: Returns structure with value as expected
		k := TestB{}
		if err := json.Unmarshal(v.message, &k); err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		if k.A.Ok() != v.expectedA {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v",
				v.expectedA, k.A)
		}
		if len(k.B) != v.expectedBSize {
			t.Errorf("[optional] Unexpected array value, expected: %v, got: %v",
				v.expectedBSize, k.B)
		}

		// When: Marshaled using JSON.Marshal
		// Then: Returns same message
		b, err := json.Marshal(k)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		b = SimplifyLiteral(b)
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		_, ok := m["a"]
		if v.expectedA != ok {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
		kb, _ := m["b"].([]interface{})
		if len(kb) != v.expectedBSize {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
	}
}

// Testing marshal/unmarshal BSON

type TestC struct {
	A optional.Int64   `bson:"a"`
	B []optional.Int64 `bson:"b"`
}

type TestD struct {
	A optional.String   `bson:"a"`
	B []optional.String `bson:"b"`
}

func TestBSONInt64(t *testing.T) {
	// Scenario: Unmarshaling optional into BSON-tagged struct
	// Given: JSON string with missing fields
	testTableA := []struct {
		msg      string
		testCase int
	}{
		// When: Missing fields are not sent
		{`{}`, 0},
		{`{"b": []}`, 0},
		{`{"b": [1, 2, 3, 4]}`, 0},
		// When: Missing fields value is null
		{`{"a": null, "b": []}`, 1},
		{`{"a": null, "b": [1, 2, 3, 4]}`, 1},
		// When: Missing fields value is undefined
		{`{"a": undefined, "b": []}`, 1},
		{`{"a": undefined, "b": [1, 2, 3, 4]}`, 1},
	}

	for _, v := range testTableA {
		var tc TestC
		msg := []byte(v.msg)
		if err := bson.UnmarshalJSON(msg, &tc); err != nil {
			t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
		}
		switch v.testCase {
		case 0:
			// Then: Optional value will indicate that respective field is missing
			fallthrough
		case 1:
			// Then: Optional value will indicate that respective field is missing
			fallthrough
		case 2:
			// Then: Optional value will indicate that respective field is missing
			if _, ok := tc.A.Get(); ok {
				t.Errorf("[optional] Fail to report missing field in message: %v", v.msg)
			}
		}
	}

	// Given: JSON string with complete fields
	testTableB := []struct {
		msg      string
		testCase int
	}{
		// When: JSON string match given structure
		{`{"a": 5, "b": [1, 2, 3, 4]}`, 0},
		{`{"a": 99, "b": [1, 2, 3, 4]}`, 0},
	}

	for _, v := range testTableB {
		var tc TestC
		msg := []byte(v.msg)
		if err := bson.UnmarshalJSON(msg, &tc); err != nil {
			t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
		}
		switch v.testCase {
		case 0:
			// Then: Optional will report that value exists
			if _, ok := tc.A.Get(); !ok {
				t.Errorf("[optional] Failed to indicate that value exists from message: %s", msg)
			}
		}
	}

	// Scenario: Marshal optional to BSON then unmarshal again should retain data
	// Given: Struct with BSON tags
	testTableC := []struct {
		singleVal optional.Int64
		arrayVal  []optional.Int64
		testCase  int
	}{
		// When: Optionals are filled
		{optional.NewInt64(int64(5), true), []optional.Int64{}, 0},
		{optional.NewInt64(int64(0), true), []optional.Int64{}, 0},
		{optional.NewInt64(int64(0), false), []optional.Int64{}, 0},
		// When: Unmarshal to different type
		{optional.NewInt64(int64(5), true), []optional.Int64{}, 1},
		{optional.NewInt64(int64(0), true), []optional.Int64{}, 1},
	}

	for _, v := range testTableC {
		tc := TestC{A: v.singleVal, B: v.arrayVal}
		msg, err := bson.Marshal(tc)
		if err != nil {
			t.Errorf("[optional] Fail to marshal data %#v, error: %v", tc, err)
		}
		switch v.testCase {
		case 0:
			// Then: Unmarshal result should match original data
			var tcu TestC
			err = bson.Unmarshal(msg, &tcu)
			if err := bson.Unmarshal(msg, &tcu); err != nil {
				t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
			}
			if !reflect.DeepEqual(tc, tcu) {
				t.Errorf("[optional] Marshal-unmarshal BSON mismatch")
				t.Errorf("[optional]	got: %v", tcu)
				t.Errorf("[optional]	expected: %v", tc)
			}
		case 1:
			// Then: Should report error
			var tcu TestD
			err = bson.Unmarshal(msg, &tcu)
			if err := bson.Unmarshal(msg, &tcu); err == nil {
				t.Errorf("[optional] Fail to detect unmarshal fail: %v, error: %v", tcu, err)
			}
		}
	}
}

func TestBSONString(t *testing.T) {
	// Scenario: Unmarshaling optional into BSON-tagged struct
	// Given: BSON string with missing fields
	testTable := []struct {
		msg      string
		testCase int
	}{
		// When: Missing fields are not sent
		{`{}`, 0},
		{`{"b": []}`, 0},
		{`{"b": ["1", "2", "3", "4"]}`, 0},
		// When: Missing fields value is null
		{`{"a": null, "b": []}`, 1},
		{`{"a": null, "b": ["1", "2", "3", "4"]}`, 1},
		// When: Missing fields value is undefined
		{`{"a": undefined, "b": []}`, 1},
		{`{"a": undefined, "b": ["1", "2", "3", "4"]}`, 1},
	}

	for _, v := range testTable {
		var td TestD
		msg := []byte(v.msg)
		if err := bson.UnmarshalJSON(msg, &td); err != nil {
			t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
		}
		switch v.testCase {
		case 0:
			// Then: Optional value will indicate that respective field is missing
			fallthrough
		case 1:
			// Then: Optional value will indicate that respective field is missing
			fallthrough
		case 2:
			// Then: Optional value will indicate that respective field is missing
			if _, ok := td.A.Get(); ok {
				t.Errorf("[optional] Fail to report missing field in message: %v", v.msg)
			}
		}
	}

	// Given: JSON string with complete fields
	testTableB := []struct {
		msg      string
		testCase int
	}{
		// When: JSON string match given structure
		{`{"a": "abc", "b": ["abc", "abc"]}`, 0},
		{`{"a": "esd", "b": ["tst", "src"]}`, 0},
	}

	for _, v := range testTableB {
		var tc TestD
		msg := []byte(v.msg)
		if err := bson.UnmarshalJSON(msg, &tc); err != nil {
			t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
		}
		switch v.testCase {
		case 0:
			// Then: Optional will report that value exists
			if _, ok := tc.A.Get(); !ok {
				t.Errorf("[optional] Failed to indicate that value exists from message: %s", msg)
			}
		}
	}

	// Scenario: Marshal optional to BSON then unmarshal again should retain data
	// Given: Struct with BSON tags
	testTableC := []struct {
		singleVal optional.String
		arrayVal  []optional.String
		testCase  int
	}{
		// When: Optionals are filled
		{optional.NewString("abc", true), []optional.String{}, 0},
		{optional.NewString("bca", true), []optional.String{}, 0},
		{optional.NewString("", false), []optional.String{}, 0},
		// When: Unmarshal to different type
		{optional.NewString("abc", true), []optional.String{}, 1},
		{optional.NewString("bca", true), []optional.String{}, 1},
	}

	for _, v := range testTableC {
		tc := TestD{A: v.singleVal, B: v.arrayVal}
		msg, err := bson.Marshal(tc)
		if err != nil {
			t.Errorf("[optional] Fail to marshal data %#v, error: %v", tc, err)
		}
		switch v.testCase {
		case 0:
			// Then: Unmarshal result should match original data
			var tcu TestD
			err = bson.Unmarshal(msg, &tcu)
			if err := bson.Unmarshal(msg, &tcu); err != nil {
				t.Errorf("[optional] Fail to unmarshal message: %s, error: %v", msg, err)
			}
			if !reflect.DeepEqual(tc, tcu) {
				t.Errorf("[optional] Marshal-unmarshal BSON mismatch")
				t.Errorf("[optional]	got: %v", tcu)
				t.Errorf("[optional]	expected: %v", tc)
			}
		case 1:
			// Then: Should report error
			var tcu TestC
			err = bson.Unmarshal(msg, &tcu)
			if err := bson.Unmarshal(msg, &tcu); err == nil {
				t.Errorf("[optional] Fail to detect unmarshal fail: %v, error: %v", tcu, err)
			}
		}
	}
}

func TestFloat64(t *testing.T) {
	// Scenario: Optional defaults to invalid
	// Given: Some optional value
	for i := float64(0); i < float64(20); i++ {
		v := optional.Float64{}
		v.Set(i, false)
		// When: Checked for validity using Ok
		// Then: Returns false
		if v.Ok() {
			t.Errorf("[optional] Invalid default value, expected: false, got: true")
		}
		// When: Checked for validity using Get
		// Then: Returns stored value and true
		if val, ok := v.Get(); ok || val != i {
			t.Errorf("[optional] Invalid default value, expected: %v and false, got: %v and %v", i, val, ok)
		}
	}

	// Scenario: Optional is invalid if unmarshal found no such key
	// Given: JSON message with missing key
	testMsgA := `{"b": [null]}`
	ta := TestA{}
	if err := json.Unmarshal([]byte(testMsgA), &ta); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgA)
	}
	if _, ok := ta.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}
	for _, k := range ta.B {
		if _, ok := k.Get(); ok {
			t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
		}
	}

	// Scenario: Optional is invalid if unmarshal type is different than expected
	// Given: JSON message with wrong value type
	testMsgB := `{"a": "test"}`
	tb := TestA{}
	if err := json.Unmarshal([]byte(testMsgB), &tb); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgB)
	}
	if _, ok := tb.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}

	// Scenario: Marshal/unmarshal JSON <==> optional
	// Given: Some JSON and corresponding structure
	testCaseA := []struct {
		message       []byte
		expectedA     bool
		expectedBSize int
	}{
		{[]byte(`{}`), false, 0},
		{[]byte(`{"a": 4}`), true, 0},
		{[]byte(`{"a": 4, "b": []}`), true, 0},
		{[]byte(`{"a": 4, "b": [1, 2, 3]}`), true, 3},
		{[]byte(`{"b": [1, 2, 3]}`), false, 3},
	}

	for _, v := range testCaseA {
		// When: Unmarshaled using JSON.Unmarshal
		// Then: Returns structure with value as expected
		k := TestA{}
		if err := json.Unmarshal(v.message, &k); err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		if k.A.Ok() != v.expectedA {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v",
				v.expectedA, k.A)
		}
		if len(k.B) != v.expectedBSize {
			t.Errorf("[optional] Unexpected array value, expected: %v, got: %v",
				v.expectedBSize, k.B)
		}

		// When: Marshaled using JSON.Marshal
		// Then: Returns same message
		b, err := json.Marshal(k)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		b = SimplifyLiteral(b)
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		_, ok := m["a"]
		if v.expectedA != ok {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
		kb, _ := m["b"].([]interface{})
		if len(kb) != v.expectedBSize {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
	}
}

func TestBool(t *testing.T) {
	type bs struct {
		A optional.Bool `json:"a"`
	}
	testCaseA := []struct {
		message []byte
		vBool   bool
		vInt    int64
		ok      bool
	}{
		{[]byte(`{}`), false, 0, false},
		{[]byte(`{"a":"asd"}`), false, 0, false},
		{[]byte(`{"a":null}`), false, 0, false},
		{[]byte(`{"a":0}`), false, 0, true},
		{[]byte(`{"a":1}`), true, 1, true},
		{[]byte(`{"a":false}`), false, 0, true},
		{[]byte(`{"a":true}`), true, 1, true},
	}

	for _, v := range testCaseA {
		// test Unmarshall
		k := bs{}
		if err := json.Unmarshal(v.message, &k); err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}

		// test Ok
		if k.A.Ok() != v.ok {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", v.ok, k.A.Ok())
		}

		// test Get
		vBool, ok := k.A.Get()
		if ok {
			if vBool != v.vBool {
				t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", v.vBool, vBool)
			}
		} else {
			continue
		}

		// test GetInt
		vInt, _ := k.A.GetInt()
		if vInt != v.vInt {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", v.vInt, vInt)
		}
	}

	testCaseB := []struct {
		message []byte
		set     bool
		ok      bool
	}{
		{[]byte(`{"a":null}`), false, false},
		{[]byte(`{"a":false}`), false, true},
		{[]byte(`{"a":true}`), true, true},
	}

	for _, v := range testCaseB {
		// test NewBool
		k := bs{
			A: optional.NewBool(v.set, v.ok),
		}

		vBool, ok := k.A.Get()
		if ok {
			if vBool != v.set {
				t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", v.set, vBool)
			}
		}

		// test Marshall
		b, err := json.Marshal(k)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		if bytes.Compare(b, v.message) != 0 {
			t.Errorf("[optional] Unexpected JSON value, expected: %v, got: %v", string(v.message), string(b))
		}
	}

	// test Set false
	testCaseC := []interface{}{
		false, 0, int8(0), int16(0), int32(0), int64(0), float32(0), float64(0), "",
	}
	for _, v := range testCaseC {
		ob := optional.NewBool(false, false)

		ob.Set(v, true)
		val, ok := ob.Get()
		if ok && val != false {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", false, val)
		}
	}

	// test Set true
	testCaseD := []interface{}{
		true, -1, 1, 2, "asd",
	}
	for _, v := range testCaseD {
		ob := optional.NewBool(false, false)

		ob.Set(v, true)
		val, ok := ob.Get()
		if ok && val != true {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v", true, val)
		}
	}
}

// Copied from TestInt64
func TestInt(t *testing.T) {
	// Scenario: Optional defaults to invalid
	// Given: Some optional value
	for i := 0; i < 20; i++ {
		v := optional.Int{}
		v.Set(i, false)
		// When: Checked for validity using Ok
		// Then: Returns false
		if v.Ok() {
			t.Errorf("[optional] Invalid default value, expected: false, got: true")
		}
		// When: Checked for validity using Get
		// Then: Returns stored value and true
		if val, ok := v.Get(); ok || val != i {
			t.Errorf("[optional] Invalid default value, expected: %v and false, got: %v and %v", i, val, ok)
		}
	}

	// Scenario: Optional is invalid if unmarshal found no such key
	// Given: JSON message with missing key
	testMsgA := `{"b": [null]}`
	ta := TestA{}
	if err := json.Unmarshal([]byte(testMsgA), &ta); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgA)
	}
	if _, ok := ta.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}
	for _, k := range ta.B {
		if _, ok := k.Get(); ok {
			t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
		}
	}

	// Scenario: Optional is invalid if unmarshal type is different than expected
	// Given: JSON message with wrong value type
	testMsgB := `{"a": "test"}`
	tb := TestA{}
	if err := json.Unmarshal([]byte(testMsgB), &tb); err != nil {
		t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, testMsgB)
	}
	if _, ok := tb.A.Get(); ok {
		t.Errorf("[optional] Unexpected unmarshal result, expected: false, got: true")
	}

	// Scenario: Marshal/unmarshal JSON <==> optional
	// Given: Some JSON and corresponding structure
	testCaseA := []struct {
		message       []byte
		expectedA     bool
		expectedBSize int
	}{
		{[]byte(`{}`), false, 0},
		{[]byte(`{"a": 4}`), true, 0},
		{[]byte(`{"a": 4, "b": []}`), true, 0},
		{[]byte(`{"a": 4, "b": [1, 2, 3]}`), true, 3},
		{[]byte(`{"b": [1, 2, 3]}`), false, 3},
	}

	for _, v := range testCaseA {
		// When: Unmarshaled using JSON.Unmarshal
		// Then: Returns structure with value as expected
		k := TestA{}
		if err := json.Unmarshal(v.message, &k); err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		if k.A.Ok() != v.expectedA {
			t.Errorf("[optional] Unexpected single value, expected: %v, got: %v",
				v.expectedA, k.A)
		}
		if len(k.B) != v.expectedBSize {
			t.Errorf("[optional] Unexpected array value, expected: %v, got: %v",
				v.expectedBSize, k.B)
		}

		// When: Marshaled using JSON.Marshal
		// Then: Returns same message
		b, err := json.Marshal(k)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		b = SimplifyLiteral(b)
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			t.Errorf("[optional] Error unmarshal JSON: %v with message: %s", err, v.message)
			continue
		}
		_, ok := m["a"]
		if v.expectedA != ok {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
		kb, _ := m["b"].([]interface{})
		if len(kb) != v.expectedBSize {
			t.Errorf("[optional] Unexpected JSON message, expected: %s, got: %s",
				v.message, b)
		}
	}
}

// SimplifyLiteral removes whitespaces, newlines as well as null value from given msg.
func SimplifyLiteral(msg []byte) (res []byte) {
	// Removes null value
	res = nullRe.ReplaceAllLiteral(msg, []byte(""))
	res = nullEndRe.ReplaceAllLiteral(res, []byte("}"))
	// Removes newlines
	res = crRe.ReplaceAllLiteral(res, []byte(""))
	// Removes whitespaces
	res = wsRe.ReplaceAllLiteral(res, []byte(""))
	// Removes comma for last key-value in document
	res = commaRe.ReplaceAllLiteral(res, []byte("}"))
	return
}
