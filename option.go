package console

// OptionFunc changes the config for customization.
type OptionFunc func(conf *Config)

// WithTitleScreen adds a title screen function which runs before the initial prompt.
func WithTitleScreen(fn func()) OptionFunc {
	return func(conf *Config) {
		conf.TitleScreenFunc = fn
	}
}
