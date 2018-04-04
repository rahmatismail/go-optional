package goptional

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

// Optional is a interface for type that has additional boolean to indicate whether
// its value valid or not.
type Optional interface {
	Ok() bool
}

// Int is optional form of int
type Int struct {
	val int
	set bool
}

// NewInt creates a new Int optional
func NewInt(val int, ok bool) Int {
	return Int{val, ok}
}

// Ok returns true if optional valus is valid.
func (i *Int) Ok() bool {
	return i.set
}

// Get returns value-flag pair of this optional.
func (i *Int) Get() (int, bool) {
	return i.val, i.set
}

// Set sets value-flag pair of this optional.
func (i *Int) Set(v int, b bool) {
	i.val = v
	i.set = b
}

// UnmarshalJSON is used to unmarshal JSON into optional value.
// If unmarshal failed, optional value is invalid (`Optional.Ok()` would return false).
func (i *Int) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		i.set = false
		return nil
	}
	if err = json.Unmarshal(b, &i.val); err != nil {
		i.set = false
		return nil
	}
	i.set = true
	return nil
}

// MarshalJSON marshals optional into JSON.
// Returns error if optional is invalid.
func (i Int) MarshalJSON() ([]byte, error) {
	if v, ok := i.Get(); ok {
		return json.Marshal(v)
	}
	var v *int64
	return json.Marshal(v)
}

// Int64 is optional form of int64.
type Int64 struct {
	val int64
	set bool
}

// NewInt64 creates a new optional
func NewInt64(val int64, ok bool) Int64 {
	return Int64{val, ok}
}

// Ok returns true if optional value is valid.
func (i *Int64) Ok() bool {
	return i.set
}

// Get returns value-flag pair of this optional.
func (i *Int64) Get() (int64, bool) {
	return i.val, i.set
}

// Set sets value-flag pair of this optional.
func (i *Int64) Set(v int64, b bool) {
	i.val = v
	i.set = b
}

// UnmarshalJSON is used to unmarshal JSON into optional value.
// If unmarshal failed, optional value is invalid (`Optional.Ok()` would return false).
func (i *Int64) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		i.set = false
		return nil
	}
	if err = json.Unmarshal(b, &i.val); err != nil {
		i.set = false
		return nil
	}
	i.set = true
	return nil
}

// MarshalJSON marshals optional into JSON.
// Returns error if optional is invalid.
func (i Int64) MarshalJSON() ([]byte, error) {
	if v, ok := i.Get(); ok {
		return json.Marshal(v)
	}
	var v *int64
	return json.Marshal(v)
}

// SetBSON implements bson.Setter
func (i *Int64) SetBSON(raw bson.Raw) error {
	if raw.Kind == 0x0A {
		i.set = false
		return nil
	}
	if err := raw.Unmarshal(&i.val); err != nil {
		i.set = false
		return fmt.Errorf("Unable to unmarshal data with kind %v to int64", raw.Kind)
	}
	i.set = true
	return nil
}

// GetBSON implements bson.Getter
func (i Int64) GetBSON() (interface{}, error) {
	if v, ok := i.Get(); ok {
		return v, nil
	}
	return nil, nil
}

// String is optional form of string.
type String struct {
	val string
	set bool
}

// NewString creates a new optional
func NewString(val string, ok bool) String {
	return String{val, ok}
}

// Ok returns true if optional value is valid.
func (i *String) Ok() bool {
	return i.set
}

// Get returns value-flag pair of this optional.
func (i *String) Get() (string, bool) {
	return i.val, i.set
}

// Set sets value-flag pair of this optional.
func (i *String) Set(v string, b bool) {
	i.val = v
	i.set = b
}

// UnmarshalJSON is used to unmarshal JSON into optional value.
// If unmarshal failed, optional value is invalid (`Optional.Ok()` would return false).
func (i *String) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		i.set = false
		return nil
	}
	if err = json.Unmarshal(b, &i.val); err != nil {
		i.set = false
		return nil
	}
	i.set = true
	return nil
}

// MarshalJSON marshals optional into JSON.
// Returns error if optional is invalid.
func (i String) MarshalJSON() ([]byte, error) {
	if v, ok := i.Get(); ok {
		return json.Marshal(v)
	}
	var v *string
	return json.Marshal(v)
}

// SetBSON implements bson.Setter
func (i *String) SetBSON(raw bson.Raw) error {
	if raw.Kind == 0x0A {
		i.set = false
		return nil
	}
	if err := raw.Unmarshal(&i.val); err != nil {
		i.set = false
		return fmt.Errorf("Unable to unmarshal data with kind %v to string", raw.Kind)
	}
	i.set = true
	return nil
}

// GetBSON implements bson.Getter
func (i String) GetBSON() (interface{}, error) {
	if v, ok := i.Get(); ok {
		return v, nil
	}
	return nil, nil
}

// Float64 is optional form of int64.
type Float64 struct {
	val float64
	set bool
}

// NewFloat64 creates a new optional
func NewFloat64(val float64, ok bool) Float64 {
	return Float64{val, ok}
}

// Ok returns true if optional value is valid.
func (f *Float64) Ok() bool {
	return f.set
}

// Get returns value-flag pair of this optional.
func (f *Float64) Get() (float64, bool) {
	return f.val, f.set
}

// Set sets value-flag pair of this optional.
func (f *Float64) Set(v float64, b bool) {
	f.val = v
	f.set = b
}

// UnmarshalJSON is used to unmarshal JSON into optional value.
// If unmarshal failed, optional value is invalid (`Optional.Ok()` would return false).
func (f *Float64) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte("null")) {
		f.set = false
		return nil
	}
	if err = json.Unmarshal(b, &f.val); err != nil {
		f.set = false
		return nil
	}
	f.set = true
	return nil
}

// MarshalJSON marshals optional into JSON.
// Returns error if optional is invalid.
func (f Float64) MarshalJSON() ([]byte, error) {
	if v, ok := f.Get(); ok {
		return json.Marshal(v)
	}
	var v *int64
	return json.Marshal(v)
}

// Bool is optional form of bool.
type Bool struct {
	val bool
	set bool
}

// NewBool creates a new optional
func NewBool(val bool, ok bool) Bool {
	return Bool{val, ok}
}

// Ok returns true if optional value is valid.
func (b *Bool) Ok() bool {
	return b.set
}

// Get returns value-flag pair of this optional.
func (b *Bool) Get() (bool, bool) {
	return b.val, b.set
}

// GetInt returns bool as int64
func (b *Bool) GetInt() (int64, bool) {
	if b.val {
		return 1, b.set
	}
	return 0, b.set
}

// Set sets value-flag pair of this optional.
func (b *Bool) Set(v interface{}, s bool) {
	r := reflect.ValueOf(v)
	z := reflect.Zero(r.Type())
	if v == false || v == z.Interface() || v == "" {
		b.val = false
	} else {
		b.val = true
	}
	b.set = s
}

// UnmarshalJSON is used to unmarshal JSON into optional value.
// If unmarshal failed, optional value is invalid (`Optional.Ok()` would return false).
func (b *Bool) UnmarshalJSON(dt []byte) (err error) {
	if bytes.Equal(dt, []byte("null")) {
		b.set = false
		return nil
	}
	s := string(dt)
	if s == "1" || s == "true" {
		b.val = true
	} else if s == "0" || s == "false" {
		b.val = false
	} else {
		b.set = false
		return nil
	}
	b.set = true
	return nil
}

// MarshalJSON marshals optional into JSON.
// Returns error if optional is invalid.
func (b Bool) MarshalJSON() ([]byte, error) {
	if v, ok := b.Get(); ok {
		return json.Marshal(v)
	}
	var v *int64
	return json.Marshal(v)
}
