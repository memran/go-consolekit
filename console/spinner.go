package console

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	spin   *spinner.Spinner
	prefix string
}

func NewSpinner(description string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = description + " "
	return &Spinner{spin: s, prefix: description}
}

func (s *Spinner) Start() {
	s.spin.Start()
}

func (s *Spinner) Success(text string) {
	s.spin.Stop()
	if text != "" {
		NewCLIRenderer().Success(text)
	}
}

func (s *Spinner) Error(text string) {
	s.spin.Stop()
	NewCLIRenderer().Error(text)
}
