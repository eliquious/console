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
	RequiredFlags    []string
	ValidateArgs     ValidationFunc
	Run              func(env *Environment, cmd *Command, args []string) error
	Suggestions      func(env *Environment, args []string) []string
	EagerSuggestions bool
	IsBuiltIn        bool
	ShouldPropagate  bool

	flags         *pflag.FlagSet
	requiredFlags []string
}

// Flags returns the pflag.FlagSet. It will initialize the FlagSet if nil.
func (cmd *Command) Flags() *pflag.FlagSet {
	if cmd.flags == nil {
		cmd.flags = pflag.NewFlagSet(cmd.Use, pflag.ContinueOnError)
	}
	return cmd.flags
}

// Execute executes the command with the given args. Flags are reset before execution.
func (cmd *Command) Execute(env *Environment, args []string) error {
	cmd.flags.Visit(func(f *pflag.Flag) { f.Changed = false })

	// Parse flags
	if err := cmd.Flags().Parse(args); err != nil {
		return err
	}

	helpFlag := cmd.Flags().Lookup("help")
	if helpFlag != nil && helpFlag.Changed {
		fmt.Println(cmd.Usage())
		return nil
	}

	// Validate flags
	if len(cmd.RequiredFlags) > 0 {
		if err := cmd.validateRequiredFlags(); err != nil {
			return err
		}
	}

	// Validate args
	if cmd.ValidateArgs != nil {
		if err := cmd.ValidateArgs(args); err != nil {
			return err
		}
	}

	if cmd.Run != nil {
		return cmd.Run(env, cmd, cmd.Flags().Args())
	}
	return errors.New("'" + cmd.Use + "' command has no run function")
}

func (cmd *Command) validateRequiredFlags() error {
	for index := 0; index < len(cmd.requiredFlags); index++ {
		flag := cmd.Flags().Lookup(cmd.requiredFlags[index])
		if flag != nil && flag.Changed {
			return fmt.Errorf("%s flag is required", flag.Name)
		}
	}
	return nil
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
	cmd.Flags().SetOutput(&buf)
	cmd.Flags().PrintDefaults()

	if len(cmd.Aliases) > 0 {
		fmt.Fprintf(&buf, "\nAliases:\n  %s\n", strings.Join(cmd.Aliases, ", "))
	}
	return buf.String()
}
