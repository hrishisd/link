package main

import (
	"strings"
	"testing"
)

func makeLinkTestNoError(t *testing.T, htmlString string, expectedLink Link) func(*testing.T) {
	return func(t *testing.T) {
		html := strings.NewReader(htmlString)
		link, err := TraverseHTML(html)
		if err != nil {
			t.Error("Error while parsing HTML with link")
		}
		if link != expectedLink {
			t.Errorf("expected %v not equal to actual %v", expectedLink, link)
		}
	}
}

func TestTraverseBadHTML(t *testing.T) {
	html := strings.NewReader("bad html")
	_, err := TraverseHTML(html)
	if err == nil {
		t.Error("No error while parsing bad html")
	}
}

func TestTraverseNoLink(t *testing.T) {
	html := strings.NewReader(`<html>
	<body>
	
	<h1>My First Heading</h1>
	
	<p>My first paragraph.</p>
	
	</body>
	</html>`)

	_, err := TraverseHTML(html)
	if err == nil {
		t.Error("No error while parsing bad html")
	}
}

func TestTraverseNoNesting(t *testing.T) {
	makeLinkTestNoError(t, `<body>
	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
  </body>
  </html>`, Link{Href: "/other-page", Text: "A link to another page"})(t)
}

func TestTraverseNested(t *testing.T) {
	makeLinkTestNoError(t, `<a href="/dog">
	<span>Something in a span</span>
	Text not in a span
	<b>Bold text!</b>
  </a>`, Link{
		Href: "/dog",
		Text: "Something in a span Text not in a span Bold text!",
	})(t)
}
