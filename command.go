package console

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// Command represents a command in the console.
type Command struct {
	Use              string
	Short            string
	Long             string
	Aliases          []string
	Flags            *pflag.FlagSet
	EagerSuggestions bool
	Suggestions      func() []string
	Run              func(env *Environment, args []string) error

	builtin bool
}

// Execute executes the command with the given args. Flags are reset before execution.
func (cmd *Command) Execute(env *Environment, args []string) error {
	cmd.Flags.Visit(func(f *pflag.Flag) { f.Changed = false })

	// Parse flags
	if err := cmd.Flags.Parse(args); err != nil {
		return err
	}

	helpFlag := cmd.Flags.Lookup("help")
	if helpFlag != nil && helpFlag.Changed {
		fmt.Println(cmd.Usage())
		return nil
	}

	if cmd.Run != nil {
		return cmd.Run(env, cmd.Flags.Args())
	}
	return errors.New("'" + cmd.Use + "' command has no run function")
}

// Usage returns the command usage.
func (cmd *Command) Usage() string {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "\n"+cmd.Short)
	if len(cmd.Long) > 0 {
		fmt.Fprintln(&buf, cmd.Long)
	}
	fmt.Fprintf(&buf, "\nUsage:\n  %s [flags] [args...]\n", cmd.Use)
	fmt.Fprintln(&buf, "\nFlags:")
	cmd.Flags.SetOutput(&buf)
	cmd.Flags.PrintDefaults()

	if len(cmd.Aliases) > 0 {
		fmt.Fprintf(&buf, "\nAliases:\n  %s\n", strings.Join(cmd.Aliases, ", "))
	}
	return buf.String()
}
