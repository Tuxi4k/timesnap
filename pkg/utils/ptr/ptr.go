// Package ptr provides type-safe pointer helpers for any type.
// It is primarily used to simplify the creation of pointers to literal values
package ptr

// To returns a pointer to the value passed as an argument.
// Example:
//
//	type User struct {
//	    Age *int
//	}
//
//	u := User{Age: ptr.To(25)}
func To[T any](v T) *T {
	return new(v)
}
