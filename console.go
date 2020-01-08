package console

import (
	"github.com/c-bata/go-prompt"
)

// OptionFunc changes the config for customization.
type OptionFunc func(conf *Config)

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
		prompt.OptionLivePrefix(env.ChangeLivePrefix),
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
			ASCIICode: []byte{0x1b, 0x7f},
			Fn:        prompt.DeleteWord,
		}),
		prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
			ASCIICode: []byte{27, 27, 91, 68},
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

// AddScope addsa a scope at the root level.
func (c *Console) AddScope(scope *Scope) {
	c.rootScope.AddSubScope(scope)
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
