package console

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

// NewScope creates a new scope.
func NewScope(name string, description string) *Scope {
	scope := &Scope{
		Name:           name,
		Description:    description,
		InitializeFunc: func(env *Environment) {},
		commands:       map[string]*Command{},
		subScopes:      map[string]*Scope{},
	}

	scope.AddCommand(&Command{
		Use:   "use",
		Short: "Use pushes a new scope onto the environment",
		Suggestions: func(env *Environment, args []string) []string {
			return scope.AvailableScopes()
		},
		EagerSuggestions: true,
		Run: func(env *Environment, cmd *Command, args []string) error {
			if len(args) == 0 {
				return errors.New("use requires an argument")
			} else if len(args) > 1 {
				return errors.New("use requires only 1 argument")
			}

			sub, ok := scope.subScopes[args[0]]
			if !ok {
				return errors.New("unknown scope")
			}

			env.Push(sub)
			return nil
		},
		IsBuiltIn: true,
	})

	scope.AddCommand(&Command{
		Use:   "help",
		Short: "Prints help info",
		Suggestions: func(env *Environment, args []string) []string {
			return scope.AvailableCommands()
		},
		EagerSuggestions: true,
		Run: func(env *Environment, cmd *Command, args []string) error {
			if len(args) > 1 {
				return errors.New("help accepts only 1 argument")
			} else if len(args) == 1 {
				cmd, ok := scope.commands[args[0]]
				if ok {
					fmt.Println(cmd.Usage())
					return nil
				}

				sub, ok := scope.subScopes[args[0]]
				if ok {
					fmt.Println(sub.Usage())
					return nil
				}
				return errors.New("unknown argument")
			}

			fmt.Println(scope.Usage())
			return nil
		},
		IsBuiltIn: true,
	})

	return scope
}

// Scope represents related commands
type Scope struct {
	Name           string
	Description    string
	InitializeFunc func(*Environment)

	commands  map[string]*Command
	subScopes map[string]*Scope
}

// Commands returns the commands in the scope.
func (s *Scope) Commands() map[string]*Command {
	return s.commands
}

// SubScopes returns a list of subscopes.
func (s *Scope) SubScopes() map[string]*Scope {
	return s.subScopes
}

// AddCommand adds a command to the scope.
func (s *Scope) AddCommand(cmd *Command) {
	if cmd.Suggestions == nil {
		cmd.Suggestions = func(*Environment, []string) []string { return nil }
	}

	// Add help flag
	helpFlag := cmd.Flags().Lookup("help")
	if helpFlag == nil {
		cmd.Flags().BoolP("help", "h", false, "Prints this help")
		cmd.Flags().Lookup("help").Hidden = true
	}
	s.commands[cmd.Use] = cmd

	for _, alias := range cmd.Aliases {
		s.commands[alias] = cmd
	}

	// if the proagating command is added after the sub-scopes it should still propagate.
	if cmd.ShouldPropagate {
		for _, sub := range s.subScopes {
			sub.AddCommand(cmd)
		}
	}
}

// AddSubScope adds a sub-scope.
func (s *Scope) AddSubScope(sub *Scope) {
	for _, cmd := range s.commands {
		if cmd.ShouldPropagate {
			sub.AddCommand(cmd)
		}
	}
	s.subScopes[sub.Name] = sub
}

// AvailableCommands returns a list of commands.
func (s *Scope) AvailableCommands() []string {
	var commands []string
	for cmd := range s.commands {
		commands = append(commands, cmd)
	}
	sort.Strings(commands)
	return commands
}

// AvailableScopes returns the available sub-scopes.
func (s *Scope) AvailableScopes() []string {
	var scopes []string
	for sub := range s.subScopes {
		scopes = append(scopes, sub)
	}
	sort.Strings(scopes)
	return scopes
}

// Execute args in a scope.
func (s *Scope) Execute(env *Environment, args []string) error {
	if len(args) == 0 {
		return errors.New("no command given")
	}

	// Execute command
	if cmd, ok := s.commands[args[0]]; ok {
		if len(args) > 0 {
			return cmd.Execute(env, args[1:])
		}
		return cmd.Execute(env, nil)
	}
	return errors.New("unknown command")
}

// Usage returns the scope usage.
func (s *Scope) Usage() string {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, s.Description)

	commands := s.AvailableCommands()
	maxLen := getMaxLength(commands)
	sort.Strings(commands)

	fmt.Fprintln(&buf, "\nUser Commands:")
	for index := 0; index < len(commands); index++ {
		cmd := s.commands[commands[index]]
		if !cmd.IsBuiltIn {
			fmt.Fprintf(&buf, "  %s    %s\n", padRight(commands[index], " ", maxLen), cmd.Short)
		}
	}

	fmt.Fprintln(&buf, "\nBuilt-in Commands:")
	for index := 0; index < len(commands); index++ {
		cmd := s.commands[commands[index]]
		if cmd.IsBuiltIn {
			if cmd.Use == commands[index] {
				fmt.Fprintf(&buf, "  %s    %s\n", padRight(commands[index], " ", maxLen), cmd.Short)
			} else {
				fmt.Fprintf(&buf, "  %s    Alias for '%s' command\n", padRight(commands[index], " ", maxLen), cmd.Use)
			}
		}
	}

	if len(s.subScopes) > 0 {
		fmt.Fprintln(&buf, "\nSub-scopes:")
		scopes := s.AvailableScopes()
		maxLen = getMaxLength(scopes)
		sort.Strings(scopes)
		for index := 0; index < len(scopes); index++ {
			fmt.Fprintf(&buf, "  %s    %s\n", padRight(scopes[index], " ", maxLen), s.subScopes[scopes[index]].Description)
		}
	}
	return buf.String()
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func getMaxLength(input []string) int {
	var maxLen int
	for index := 0; index < len(input); index++ {
		if len(input[index]) > maxLen {
			maxLen = len(input[index])
		}
	}
	return maxLen
}
