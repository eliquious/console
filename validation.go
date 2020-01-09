package console

import "fmt"

// ValidationFunc validates input args
type ValidationFunc func([]string) error

// ExactArgs validates command args to an exact number
func ExactArgs(n int) ValidationFunc {
	return func(args []string) error {
		if len(args) != n {
			return fmt.Errorf("requires %d args", n)
		}
		return nil
	}
}

// MinimumArgs validates command args to greater than the given number
func MinimumArgs(n int) ValidationFunc {
	return func(args []string) error {
		if len(args) < n {
			return fmt.Errorf("requires at least %d args", n)
		}
		return nil
	}
}

// MaximumArgs validates command args to lesser than the given number
func MaximumArgs(n int) ValidationFunc {
	return func(args []string) error {
		if len(args) > n {
			return fmt.Errorf("requires no more than %d args", n)
		}
		return nil
	}
}

// CombineValidation validates command args against several validators
func CombineValidation(validators ...ValidationFunc) ValidationFunc {
	return func(args []string) error {
		for _, fn := range validators {
			if err := fn(args); err != nil {
				return err
			}
		}
		return nil
	}
}
