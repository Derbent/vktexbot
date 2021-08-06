package render

import (
	"image/png"
	"os"
	"testing"
)

func TestPDFRender(t *testing.T) {
	r := New()

	text :=
		`\documentclass[preview,border=3pt,3pt]{standalone}
		\begin{document}
		This is a LaTeX document.
		\end{document}`

	img, err := r.Rend(text)
	if err != nil {
		t.Errorf("can't rend: %v", r)
		return
	}

	f, err := os.OpenFile("test_files/out.png", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("can't create file: %v", err)
	}
	defer f.Close()
	png.Encode(f, img)
}
