package render

import (
	"image"

	pdf2png "github.com/brunsgaard/go-pdfium-render"
	"github.com/rwestlund/gotex"
)

type Render interface {
	Rend(text string) (*image.RGBA, error)
}

type render struct {
	options gotex.Options
}

func New() Render {
	opt := gotex.Options{
		Command: "/usr/bin/pdflatex",
		Runs:    1,
	}

	return &render{
		options: opt,
	}
}

func (r *render) Rend(text string) (*image.RGBA, error) {
	// textPreprocessor

	// text2pdf
	pdf, err := gotex.Render(text, r.options)
	if err != nil {
		return nil, err
	}

	// pdf2png
	pdf2png.InitLibrary()
	doc, err := pdf2png.NewDocument(&pdf)
	if err != nil {
		return nil, err
	}
	img := doc.RenderPage(0, 300)
	doc.Close()
	pdf2png.DestroyLibrary()

	return img, nil
}
