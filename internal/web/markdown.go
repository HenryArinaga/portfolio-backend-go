// ./internal/web/markdown.go
package web

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
	),
)
var policy = bluemonday.UGCPolicy()

func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, `\r\n`, "\n")
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "    ")
	return s
}

func RenderMarkdown(src string) (template.HTML, error) {
	src = normalizeWhitespace(src)

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}

	safe := policy.Sanitize(buf.String())
	return template.HTML(safe), nil
}
