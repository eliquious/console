package console

import (
	"github.com/c-bata/go-prompt"
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
	DescriptionBGColor:           prompt.LightGray,
	DescriptionTextColor:         prompt.DarkGray,
	SuggestionBGColor:            prompt.DarkGray,
	SuggestionTextColor:          prompt.White,
	SelectedSuggestionBGColor:    prompt.LightGray,
	SelectedSuggestionTextColor:  prompt.DarkGray,
	SelectedDescriptionBGColor:   prompt.DarkGray,
	SelectedDescriptionTextColor: prompt.LightGray,
}
