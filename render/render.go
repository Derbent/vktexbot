package render

import (
	"fmt"
	"image"

	"github.com/ZashX/vktexbot/gotex"
	pdf2png "github.com/brunsgaard/go-pdfium-render"
)

type Render interface {
	Rend(text string) (*image.RGBA, error)
}

type render struct {
	options  gotex.Options
	template string
}

func New() Render {
	opt := gotex.Options{
		Command: "/usr/bin/pdflatex",
		Runs:    1,
	}

	template := `
		\documentclass[preview,border=20pt,20pt]{standalone}
		\setlength{\textwidth}{8.3cm}
		\usepackage{amsmath,amsthm,amssymb,amsfonts,mathtools,mathtext,physics}
		\usepackage[T1,T2A]{fontenc}
		\usepackage[utf8]{inputenc}
		\usepackage[english,russian]{babel}
		\usepackage{listings}
		\usepackage{xcolor}

		\definecolor{codegreen}{rgb}{0,0.6,0}
		\definecolor{codegray}{rgb}{0.5,0.5,0.5}
		\definecolor{codepurple}{rgb}{0.58,0,0.82}
		\definecolor{backcolour}{rgb}{0.95,0.95,0.92}

		\lstdefinestyle{mystyle}{
			backgroundcolor=\color{backcolour},
			commentstyle=\color{codegreen},
			keywordstyle=\color{magenta},
			numberstyle=\tiny\color{codegray},
			stringstyle=\color{codepurple},
			basicstyle=\ttfamily\footnotesize,
			breakatwhitespace=false,
			breaklines=true,
			captionpos=b,
			keepspaces=true,
			numbers=left,
			numbersep=5pt,
			showspaces=false,
			showstringspaces=false,
			showtabs=false,
			tabsize=2
		}

        \lstset{style=mystyle}

		\begin{document}
		%v
		\end{document}`

	return &render{
		options:  opt,
		template: template,
	}
}

func (r *render) Rend(text string) (*image.RGBA, error) {
	// textPreprocessor

	// text2pdf
	pdf, err := gotex.Render(fmt.Sprintf(r.template, text), r.options)
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
