package complete

// Command is an object that can be used to create complete options for a go
// executable that does not have a good binding to the `Completer` interface,
// or to use a Go program as complete binary for another executable.
type Command struct {
	// Sub is map of sub commands of the current command. The key refer to the
	// sub command name, and the value is it's command descriptive struct.
	Sub map[string]*Command

	// Flags is a map of flags that the command accepts. The key is the flag
	// name, and the value is it's predictions. In a chain of sub commands, no
	// duplicate flags should be defined.
	Flags map[string]Predictor

	// Args are extra arguments that the command accepts, those who are given
	// without any flag before. In any chain of sub commands, only one of them
	// should predict positional arguments.
	Args Predictor
}

// Complete runs the completion of the described command.
func (c *Command) Complete(name string) {
	Complete(name, c)
}

func (c *Command) SubCmdList() []string {
	subs := make([]string, 0, len(c.Sub))
	for sub := range c.Sub {
		subs = append(subs, sub)
	}
	return subs
}

func (c *Command) SubCmdGet(cmd string) Completer {
	if c.Sub[cmd] == nil {
		return nil
	}
	return c.Sub[cmd]
}

func (c *Command) FlagList() []string {
	flags := make([]string, 0, len(c.Flags))
	for flag := range c.Flags {
		flags = append(flags, flag)
	}
	return flags
}

func (c *Command) FlagGet(flag string) Predictor {
	return PredictFunc(func(prefix string) (options []string) {
		f := c.Flags[flag]
		if f == nil {
			return nil
		}
		return f.Predict(prefix)
	})
}

func (c *Command) ArgsGet() Predictor {
	return PredictFunc(func(prefix string) (options []string) {
		if c.Args == nil {
			return nil
		}
		return c.Args.Predict(prefix)
	})
}
