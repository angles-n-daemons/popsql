package assert

import (
	"reflect"
	"testing"
)

// Comparable is a generic interface that requires the implementation of the
// Equal method. The Equal method should compare the current object with
// another object of the same type and return true if they are equal, otherwise
// false.
type Comparable[T any] interface {
	Equal(other T) bool
}

// fail logs a fatal error message indicating that two objects are not equal.
// It takes a testing object, and two objects of any type to compare.
func fail(t *testing.T, message string, expected, actual any) {
	t.Fatalf(message, expected, actual)
}

// equal is the helper function which executes the logic required for Equal
func equal[T any](expected, actual T) bool {
	if cmpA, ok := any(expected).(Comparable[T]); ok {
		if cmpA.Equal(actual) {
			return true
		}
	} else if reflect.DeepEqual(expected, actual) {
		return true
	}
	return false
}

// Equal compares two objects of any type for equality. If the objects
// implement the Comparable interface, it uses the Equal method for comparison.
// Otherwise, it uses reflect.DeepEqual to compare the objects. If the objects
// are not equal, it calls the fail function to log a fatal error.
func Equal[T any](t *testing.T, expected, actual T) {
	if !equal(expected, actual) {
		fail(t, "Values not equal, expected:\n\t%v\nactual:\n\t%v", expected, actual)
	}
}

// NotEqual is the exact inverse of Equal. It uses the Equal function if the
// passed values are Comparators, otherwise it uses reflect to assert that the
// two values are not equal.
func NotEqual[T any](t *testing.T, first, second T) {
	if equal(first, second) {
		fail(t, "expected values not to be equal, first:\n\t%v\nsecond:\n\t%v", first, second)
	}
}

// IsError checks if the provided error's message matches the expected message.
// It takes a testing object, an error, and a string message as parameters.
// If the message is empty, it assumes that there should be no error
// If the error is nil, it logs a fatal error indicating that an error was
// expected. Otherwise, the function uses the Equal function to compare the
// error's message with the expected message. If the messages do not match, it
// logs a fatal error.
func IsError(t *testing.T, err error, message string) {
	if message == "" {
		NoError(t, err)
		return
	} else if err == nil {
		t.Fatal("expected an error")
	} else if err.Error() != message {
		fail(t, "Expected error '%s', got '%s'", message, err.Error())
	}
}

// NoError checks if the provided error is nil.
// It takes a testing object and an error as parameters.
// If the error is not nil, it logs a fatal error with the error message.
func NoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// isNil checks if a specified object is nil or not.
// copied from testify/assert, link:
// https://github.com/stretchr/testify/blob/v1.10.0/assert/assertions.go#L685
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:

		return value.IsNil()
	}

	return false
}

// Nil checks if the provided value is nil.
// It takes a testing object and an interface value as parameters.
// If the value is not nil, it logs a fatal error with a message.
func Nil(t *testing.T, v interface{}) {
	if !isNil(v) {
		t.Fatalf("Expected nil, got %v", v)
	}
}
