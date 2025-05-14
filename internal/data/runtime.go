package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Declare a custom Runtime type, which has the underlying type int32 (the same as our
// Movie struct field).
type Runtime int32

// custom error type when decoding json from the movies Runtime key
var ErrInvalidRuntimTypeFormat = errors.New("invalid runtime property type format")

// Implement a MarshalJSON() method on the Runtime type so that it satisfies the
// json.Marshaler interface. This should return the JSON-encoded value for the movie
// runtime (in our case, it will return a string in the format "<runtime> mins").
func (r *Runtime) MarshalJSON() ([]byte, error) {
	jsValue := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*
	quotedJSONValue := strconv.Quote(jsValue)

	return []byte(quotedJSONValue), nil
}

// Implement a UnmarshalJSON() method on the Runtime type so that it satisfies the
// json.Unmarshaler interface. This should take the JSON-encoded data and convert it
// back into a Runtime type.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// Use strconv.Unquote() to remove the surrounding quotes from the JSON string
	jsStr, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return err
	}

	// Split the string to isolate the part containing the number.
	jsStrParts := strings.Split(jsStr, " ")
	if len(jsStrParts) != 2 || jsStrParts[1] != "mins" {
		return ErrInvalidRuntimTypeFormat
	}

	// Parse the string containing the number into an int32.
	i, err := strconv.ParseInt(jsStrParts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimTypeFormat
	}

	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer
	*r = Runtime(i)

	return nil
}
