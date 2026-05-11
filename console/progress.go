package console

import (
	"github.com/schollz/progressbar/v3"
)

type Progress struct {
	bar *progressbar.ProgressBar
}

func NewProgress(description string, total int) *Progress {
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionClearOnFinish(),
	)
	return &Progress{bar: bar}
}

func (p *Progress) Advance() {
	p.bar.Add(1)
}

func (p *Progress) Finish() {
	p.bar.Finish()
}

func (p *Progress) Run(fn func(*Progress)) {
	fn(p)
	p.Finish()
}
