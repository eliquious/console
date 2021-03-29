package console

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/dop251/goja"
	"github.com/eliquious/console/colors"
)

func NewEvalCommand() *Command {

	// console config for the eval prompt
	conf := &Config{
		Title:           "goja",
		Prefix:          "eval",
		ColorScheme:     DefaultColorScheme,
		TitleScreenFunc: func() {},
	}

	// dynamic prefix function
	var environ *Environment
	prefixFunc := func() (string, bool) {
		scopes := []string{}
		for index := 0; index < len(environ.ScopeStack); index++ {
			scopes = append(scopes, environ.ScopeStack[index].Name)
		}
		scopes = append(scopes, "eval")
		return strings.Join(scopes, ":") + environ.Prefix, true
	}

	// create a new JS virtual machine
	vm := goja.New()

	// executor function
	executor := func(line string) {
		line = strings.TrimSpace(line)
		if line == "pop" || line == "exit" {
			fmt.Println("Press Ctrl-D to exit")
			return
		}

		val, err := vm.RunString(line)
		if err != nil {
			fmt.Printf(colors.Red("error: ", err))
		}
		fmt.Println(val)
	}
	evalPrompt := newEvalPrompt(conf, prefixFunc, executor)

	command := &Command{
		Use:              "eval",
		Short:            "Launch JS interpreter",
		EagerSuggestions: true,
		Run: func(env *Environment, cmd *Command, args []string) error {
			environ = env

			evalPrompt.Run()
			return nil
		},
		IsBuiltIn: true,
	}
	return command
}

func newEvalPrompt(conf *Config, prefixFunc func() (string, bool), execFunc prompt.Executor) *prompt.Prompt {

	promptOpts := []prompt.Option{
		prompt.OptionTitle(conf.Title),
		prompt.OptionPrefix(conf.Prefix),
		prompt.OptionLivePrefix(prefixFunc),
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

	completer := func(prompt.Document) []prompt.Suggest { return []prompt.Suggest{} }
	return prompt.New(execFunc, completer, promptOpts...)
}
