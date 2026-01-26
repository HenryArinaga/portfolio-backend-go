// ./internal/web/markdown.go
package web

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

var md = goldmark.New()
var policy = bluemonday.UGCPolicy()

func RenderMarkdown(src string) (template.HTML, error) {
	src = strings.ReplaceAll(src, `\n`, "\n")

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}

	safe := policy.Sanitize(buf.String())
	return template.HTML(safe), nil
}
