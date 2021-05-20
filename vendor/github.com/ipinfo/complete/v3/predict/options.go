package predict

import (
	"fmt"
	"github.com/ipinfo/complete/v3"
	"strings"
)

// Option provides prediction through options pattern.
//
// Usage:
//
//  func(o ...predict.Option) {
//  	cfg := predict.Options(o)
//  	// use cfg.Predict...
//  }
type Option func(*Config)

// OptValues allows to set a desired set of valid values for the flag.
func OptValues(values ...string) Option {
	return OptPredictor(Set(values))
}

// OptPredictor allows to set a custom predictor.
func OptPredictor(p complete.Predictor) Option {
	return func(o *Config) {
		if o.Predictor != nil {
			panic("predictor set more than once.")
		}
		o.Predictor = p
	}
}

// OptCheck enforces the valid values on the predicted flag.
func OptCheck() Option {
	return func(o *Config) {
		if o.check {
			panic("check set more than once")
		}
		o.check = true
	}
}

// Config stores prediction options.
type Config struct {
	complete.Predictor
	check bool
}

// Options return a config from a list of options.
func Options(os ...Option) Config {
	var op Config
	for _, f := range os {
		f(&op)
	}
	return op
}

func (c Config) Predict(prefix string) []string {
	if c.Predictor != nil {
		return c.Predictor.Predict(prefix)
	}
	return nil
}

// Check checks that value is one of the predicted values, in case
// that the check field was set.
func (c Config) Check(value string) error {
	if !c.check || c.Predictor == nil {
		return nil
	}
	predictions := c.Predictor.Predict(value)
	if len(predictions) == 0 {
		return nil
	}
	for _, vv := range predictions {
		if value == vv {
			return nil
		}
	}
	return fmt.Errorf("not in allowed values: %s", strings.Join(predictions, ","))
}
