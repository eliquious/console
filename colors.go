package console

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/gookit/color"
)

// ColorScheme sets the color scheme for the console
type ColorScheme struct {
	ScrollbarThumbColor          prompt.Color
	ScrollbarBGColor             prompt.Color
	PrefixTextColor              prompt.Color
	InputTextColor               prompt.Color
	DescriptionBGColor           prompt.Color
	DescriptionTextColor         prompt.Color
	SuggestionBGColor            prompt.Color
	SuggestionTextColor          prompt.Color
	SelectedSuggestionBGColor    prompt.Color
	SelectedSuggestionTextColor  prompt.Color
	SelectedDescriptionBGColor   prompt.Color
	SelectedDescriptionTextColor prompt.Color
}

// DefaultColorScheme is the default color scheme for clix.
var DefaultColorScheme = &ColorScheme{
	ScrollbarThumbColor:          prompt.Red,
	ScrollbarBGColor:             prompt.White,
	PrefixTextColor:              prompt.Red,
	InputTextColor:               prompt.White,
	DescriptionBGColor:           prompt.White,
	DescriptionTextColor:         prompt.DarkGray,
	SuggestionBGColor:            prompt.DarkGray,
	SuggestionTextColor:          prompt.White,
	SelectedSuggestionBGColor:    prompt.White,
	SelectedSuggestionTextColor:  prompt.DarkGray,
	SelectedDescriptionBGColor:   prompt.DarkGray,
	SelectedDescriptionTextColor: prompt.White,
}

// PrintInfo prints info with a green label
func PrintInfo(label string, format string, value ...interface{}) {
	fmt.Printf("%s: %s\n", color.LightGreen.Render(label), fmt.Sprintf(format, value...))
}
