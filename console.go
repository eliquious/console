package console

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/c-bata/go-prompt"
)

// Config is the app configuration.
type Config struct {
	Title           string
	Prefix          string
	MaxSuggestions  uint16
	ColorScheme     *ColorScheme
	TitleScreenFunc func()
}

// New creates a new Console.
func New(name string, opts ...OptionFunc) *Console {
	env := NewEnvironment("> ")
	rootScope := NewScope(name, "")

	// setup built-in commands
	addBuiltInCommands(rootScope)

	env.Push(rootScope)

	conf := &Config{
		Title:           "console",
		Prefix:          "> ",
		MaxSuggestions:  8,
		ColorScheme:     DefaultColorScheme,
		TitleScreenFunc: func() {},
	}

	for _, opt := range opts {
		opt(conf)
	}

	promptOpts := []prompt.Option{
		prompt.OptionTitle(conf.Title),
		prompt.OptionPrefix(conf.Prefix),
		prompt.OptionLivePrefix(env.LivePrefix),
		prompt.OptionMaxSuggestion(conf.MaxSuggestions),

		// Text colors
		prompt.OptionScrollbarThumbColor(conf.ColorScheme.ScrollbarThumbColor),
		prompt.OptionScrollbarBGColor(conf.ColorScheme.ScrollbarBGColor),
		prompt.OptionPrefixTextColor(conf.ColorScheme.PrefixTextColor),
		prompt.OptionInputTextColor(conf.ColorScheme.InputTextColor),
		prompt.OptionDescriptionBGColor(conf.ColorScheme.DescriptionBGColor),
		prompt.OptionDescriptionTextColor(conf.ColorScheme.DescriptionTextColor),
		prompt.OptionSuggestionBGColor(conf.ColorScheme.SuggestionBGColor),
		prompt.OptionSuggestionTextColor(conf.ColorScheme.SuggestionTextColor),
		prompt.OptionSelectedSuggestionBGColor(conf.ColorScheme.SelectedSuggestionBGColor),
		prompt.OptionSelectedSuggestionTextColor(conf.ColorScheme.SelectedSuggestionTextColor),
		prompt.OptionSelectedDescriptionBGColor(conf.ColorScheme.SelectedDescriptionBGColor),
		prompt.OptionSelectedDescriptionTextColor(conf.ColorScheme.SelectedDescriptionTextColor),

		// Key bindings for meta key
		prompt.OptionSwitchKeyBindMode(prompt.EmacsKeyBind),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 127},
			Fn:        prompt.DeleteWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 0x08},
			Fn:        prompt.DeleteWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			// ASCIICode: []byte{27, 27, 91, 68},
			ASCIICode: []byte{0x1b, 98},
			Fn:        prompt.GoLeftWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{0x1b, 102},
			Fn:        prompt.GoRightWord,
		}),
	}

	p := prompt.New(env.ExecutorFunc, env.CompletorFunc, promptOpts...)
	return &Console{
		config:    conf,
		env:       env,
		rootScope: rootScope,
		prompt:    p,
	}
}

// Console runs the prompt and manages the environment.
type Console struct {
	config    *Config
	env       *Environment
	prompt    *prompt.Prompt
	rootScope *Scope
}

// AddScope adds a scope at the root level.
func (c *Console) AddScope(scope *Scope) {
	c.rootScope.AddSubScope(scope)
}

// AddCommand adds a command at the root level.
func (c *Console) AddCommand(cmd *Command) {
	c.rootScope.AddCommand(cmd)
}

// Environment returns the console environment
func (c *Console) Environment() *Environment {
	return c.env
}

// Run runs the console
func (c *Console) Run() {
	c.config.TitleScreenFunc()
	c.prompt.Run()
}

func addBuiltInCommands(scope *Scope) {

	scope.AddCommand(&Command{
		Use:   "env",
		Short: "env lists all the environment variables for the commands",
		Run: func(env *Environment, cmd *Command, args []string) error {
			keys := env.Configuration.AllKeys()
			sort.Strings(keys)

			maxLen := getMaxLength(keys)
			for index := 0; index < len(keys); index++ {
				fmt.Printf("%s   %v\n", padRight(keys[index], " ", maxLen), env.Configuration.Get(keys[index]))
			}
			return nil
		},
		IsBuiltIn:       true,
		ShouldPropagate: true,
	})
	scope.AddCommand(&Command{
		Use:              "get",
		Short:            "Gets a current env var",
		EagerSuggestions: true,
		Suggestions: func(env *Environment, args []string) []string {
			if len(args) < 2 {
				keys := env.Configuration.AllKeys()
				sort.Strings(keys)
				return keys
			}
			return []string{}
		},
		Run: func(env *Environment, cmd *Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires 1 argument")
			}
			fmt.Printf("%s   %v\n", args[0], env.Configuration.Get(args[0]))
			return nil
		},
		IsBuiltIn:       true,
		ShouldPropagate: true,
	})

	scope.AddCommand(&Command{
		Use:              "set",
		Short:            "Sets an env var",
		EagerSuggestions: true,
		Suggestions: func(env *Environment, args []string) []string {
			if len(args) < 2 {
				keys := env.Configuration.AllKeys()
				sort.Strings(keys)
				return keys
			}
			return []string{}
		},
		Run: func(env *Environment, cmd *Command, args []string) error {
			if len(args) != 2 {
				return errors.New("requires 2 arguments")
			}
			env.Configuration.Set(args[0], args[1])
			// fmt.Printf("%s   %v\n", args[0], env.Configuration.Get(args[0]))
			return nil
		},
		IsBuiltIn:       true,
		ShouldPropagate: true,
	})
	scope.AddCommand(&Command{
		Use:     "exit",
		Aliases: []string{"pop"},
		Short:   "Exit pops a scope from the environment. Exits console if at the root scope.",
		Run: func(env *Environment, cmd *Command, args []string) error {
			env.Pop()
			return nil
		},
		IsBuiltIn:       true,
		ShouldPropagate: true,
	})

	scope.AddCommand(&Command{
		Use:   "quit",
		Short: "Exits the console regardless of scope",
		Run: func(env *Environment, cmd *Command, args []string) error {
			os.Exit(0)
			return nil
		},
		IsBuiltIn:       true,
		ShouldPropagate: true,
	})
}
