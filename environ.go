package console

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"

	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
	"github.com/kballard/go-shellquote"
	"github.com/spf13/pflag"
)

const Suggestions = "suggestions"

// NewEnvironment creates a new environment with a root scope.
func NewEnvironment(prefix string) *Environment {
	env := &Environment{ScopeStack: make([]*Scope, 0), Prefix: prefix, Configuration: viper.New()}
	return env
}

// Environment manages the various cmd scopes
type Environment struct {
	Prefix        string
	ScopeStack    []*Scope
	Configuration *viper.Viper
}

// LivePrefix allows for a dynamic prompt prefix
func (env *Environment) LivePrefix() (string, bool) {
	scopes := []string{}
	for index := 0; index < len(env.ScopeStack); index++ {
		scopes = append(scopes, env.ScopeStack[index].Name)
	}
	return strings.Join(scopes, ":") + env.Prefix, true
}

// // Set sets an environment variable
// func (env *Environment) Set(key, value string) {
// 	env.EnvVariables[strings.ToUpper(key)] = value
// }

// // Get returns an env variable.
// func (env *Environment) Get(key string) string {
// 	if val, ok := env.EnvVariables[strings.ToUpper(key)]; ok {
// 		return val
// 	}
// 	return ""
// }

// Push adds a scope to the environment
func (env *Environment) Push(scope *Scope) {
	if scope.InitializeFunc != nil {
		scope.InitializeFunc(env)
	}
	env.ScopeStack = append(env.ScopeStack, scope)
}

// Len returns the number of scopes. Should always be at least 1.
func (env *Environment) Len() int {
	return len(env.ScopeStack)
}

// Pop removes a scope from the environment. Should never remove the root scope.
func (env *Environment) Pop() *Scope {
	if len(env.ScopeStack) <= 1 {
		os.Exit(0)
		return nil
	}
	scope := env.CurrentScope()
	env.ScopeStack = env.ScopeStack[:env.Len()-1]
	return scope
}

// CurrentScope gets the current scope from the environment
func (env *Environment) CurrentScope() *Scope {
	return env.ScopeStack[env.Len()-1]
}

// ExecutorFunc executes the input.
func (env *Environment) ExecutorFunc(input string) {
	if input == "" {
		return
	}

	// Parse the input
	args, err := shellquote.Split(input)
	if err != nil {
		color.Warn.Println(err.Error())
		return
	}

	// Get the current scope
	scope := env.CurrentScope()
	if scope == nil {
		color.Warn.Println("current scope is nil")
		return
	}

	// Execute the command
	if err := scope.Execute(env, args); err != nil {
		color.Error.Println(err.Error())
		return
	}
}

// CompletorFunc gets the Completer from the current scope.
func (env *Environment) CompletorFunc(doc prompt.Document) []prompt.Suggest {
	line := doc.CurrentLine()
	if strings.TrimSpace(line) == "" {
		return []prompt.Suggest{}
	}

	// Parse the input
	args, err := shellquote.Split(doc.TextBeforeCursor())
	if err != nil {
		// color.Warn.Println(err.Error())
		return []prompt.Suggest{}
	}

	// Get suggestions from current scope
	scope := env.CurrentScope()
	suggestions := GetSuggestions(env, line, scope.Commands(), doc.GetWordBeforeCursor(), args)
	return prompt.FilterFuzzy(suggestions, doc.GetWordBeforeCursor(), true)
}

// GetSuggestions returns the suggestions for the given input and commands.
func GetSuggestions(env *Environment, line string, commands map[string]*Command, prevWord string, args []string) []prompt.Suggest {
	rootCompletions := []prompt.Suggest{}

	var commandNames []string
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	for index := 0; index < len(commandNames); index++ {
		name := commandNames[index]
		cmd := commands[name]

		if strings.HasPrefix(line, cmd.Use) {
			if len(args) > 0 && args[0] == cmd.Use {
				return getCommandSuggestions(env, line, cmd, prevWord, args[0:])
			}
			return getCommandSuggestions(env, line, cmd, prevWord, args)
		}

		for _, alias := range cmd.Aliases {
			if strings.HasPrefix(line, alias) {
				if len(args) > 0 && args[0] == alias {
					return getCommandSuggestions(env, line, cmd, prevWord, args[0:])
				}
				return getCommandSuggestions(env, line, cmd, prevWord, args)
			}
		}

		sug := prompt.Suggest{Text: name, Description: cmd.Short}
		if name != cmd.Use {
			sug.Description = fmt.Sprintf("Alias for `%s`. %s", cmd.Use, cmd.Short)
		}
		rootCompletions = append(rootCompletions, sug)
	}
	return rootCompletions
}

func getCommandSuggestions(env *Environment, line string, cmd *Command, prevWord string, args []string) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// Add args suggestions
	if len(prevWord) > 0 || cmd.EagerSuggestions {
		for _, sug := range cmd.Suggestions(env, args) {
			suggestions = append(suggestions, prompt.Suggest{Text: sug})
		}
	}

	// Add flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
		}
	})

	// Add matching flag suggestions
	if strings.HasPrefix(prevWord, "-") && strings.HasSuffix(prevWord, "=") {
		flagString := strings.TrimLeft(prevWord, "-")
		flagString = strings.TrimRight(flagString, "=")

		flag := cmd.Flags().Lookup(flagString)
		if flagSuggestions, ok := flag.Annotations[Suggestions]; ok {
			for _, sug := range flagSuggestions {
				suggestions = append(suggestions, prompt.Suggest{Text: sug})
			}
		}
	}
	return suggestions
}
