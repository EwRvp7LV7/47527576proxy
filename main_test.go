package main

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

type testpair struct {
	raw   string
	ready string
}

var tests = []testpair{
	{"<html><head></head><body><p>лооООооонг</p></body></html>",
		"<html><head></head><body><p>лооООооонг\u2122</p></body></html>"},
	{"<html><head></head><body><p>loooong</p></body></html>",
		"<html><head></head><body><p>loooong\u2122</p></body></html>"},
	{"<html><head></head><body><li><a href=\"https://habr.com/questions/48611185/\">short</a></li><li></li></body></html>",
		"<html><head></head><body><li><a href=\"/questions/48611185/\">short</a></li><li></li></body></html>"},
}

func innerHTML(n *html.Node) string {
	var b bytes.Buffer
	html.Render(&b, n)
	return b.String()
}

func TestSplitText(t *testing.T) {
	// fmt.Print("error1")
	// t.Error("error")
	for _, pair := range tests {

		body, _ := html.Parse(strings.NewReader(pair.raw))

		Replace(body)
		innerHTML := innerHTML(body)

		if innerHTML != pair.ready {
			t.Error(
				"For", pair.raw,
				"expected", pair.ready,
				"got", innerHTML,
			)
		}
	}
}
