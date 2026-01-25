package web

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

var md = goldmark.New()
var policy = bluemonday.UGCPolicy()

func RenderMarkdown(src string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return policy.Sanitize(buf.String()), nil
}
