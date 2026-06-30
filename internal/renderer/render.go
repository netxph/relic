package renderer

import (
	_ "embed"
	"io"
	"os"
	"text/template"
)

//go:embed templates/default.tmpl
var defaultTemplate string

// Render writes the release data as markdown to w using the given template string.
// Pass empty tmplStr to use the embedded default template.
func Render(data ReleaseData, tmplStr string, w io.Writer) error {
	if tmplStr == "" {
		tmplStr = defaultTemplate
	}
	t, err := template.New("release").Parse(tmplStr)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// LoadTemplate reads a template from a file path.
func LoadTemplate(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
