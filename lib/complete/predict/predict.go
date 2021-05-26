// Package predict provides helper functions for completion predictors.
package predict

import "github.com/ipinfo/cli/lib/complete"

// Set predicts a set of predefined values.
type Set []string

func (p Set) Predict(_ string) (options []string) {
	return p
}

var (
	// Something is used to indicate that does not completes somthing. Such that other prediction
	// wont be applied.
	Something = Set{""}

	// Nothing is used to indicate that does not completes anything.
	Nothing = Set{}
)

// Or unions prediction functions, so that the result predication is the union of their
// predications.
func Or(ps ...complete.Predictor) complete.Predictor {
	return complete.PredictFunc(func(prefix string) (options []string) {
		for _, p := range ps {
			if p == nil {
				continue
			}
			options = append(options, p.Predict(prefix)...)
		}
		return
	})
}
