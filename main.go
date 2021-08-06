package main

import (
	"context"
	"image/png"
	"log"
	"os"

	render "github.com/ZashX/vktexbot/render"
)

func main() {
	ctx := context.Background()
	if err := Run(ctx); err != nil {
		log.Fatalf("unexpected result: %v", err)
	}

	return
}

func Run(ctx context.Context) error {
	r := render.New()

	text := `\documentclass[preview,border=3pt,3pt]{standalone}
	\begin{document}
	This is a LaTeX document.
	\end{document}`

	img, err := r.Rend(text)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
	return nil
}
