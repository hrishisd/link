package main

import (
	"strings"
	"testing"
)

func makeLinkTest(t *testing.T, htmlString string, expectedLinks []Link) func() {
	expectedContains := func(link Link) bool {
		for _, other := range expectedLinks {
			if other == link {
				return true
			}
		}
		return false
	}

	linksAreExpected := func(links []Link) bool {
		if len(links) != len(expectedLinks) {
			return false
		}
		for _, link := range links {
			if !expectedContains(link) {
				return false
			}
		}
		return true
	}

	return func() {
		html := strings.NewReader(htmlString)
		links, err := TraverseHTML(html)
		if err != nil {
			t.Error("Error while parsing HTML with link")
		}
		if !linksAreExpected(links) {
			t.Errorf("expected %v not equal to actual %v", expectedLinks, links)
		}
	}
}

func TestTraverseBadHTML(t *testing.T) {
	htmlString := "bad html lalala"
	makeLinkTest(t, htmlString, nil)()
}

func TestTraverseNoLink(t *testing.T) {
	htmlString := `<html>
	<body>
	
	<h1>My First Heading</h1>
	
	<p>My first paragraph.</p>
	
	</body>
	</html>`
	makeLinkTest(t, htmlString, nil)()
}

func TestTraverseNoNesting(t *testing.T) {
	htmlString := `<body>
	<h1>Hello!</h1>
	<a href="/other-page">A link to another page</a>
  </body>
  </html>`
	makeLinkTest(t, htmlString, []Link{{"/other-page", "A link to another page"}})()
}

func TestTraverseSingleLinkNestedBody(t *testing.T) {
	htmlString := `<a href="/dog">
	<span>Something in a span</span>
	Text not in a span
	<b>Bold text!</b>
  </a>`
	expectedLinks := []Link{{Href: "/dog", Text: "Something in a span Text not in a span Bold text!"}}
	makeLinkTest(t, htmlString, expectedLinks)
}

func TestTraverseMultipleLinks(t *testing.T) {
	htmlString := `</head>
	<body>
	  <h1>Social stuffs</h1>
	  <div>
		<a href="https://www.twitter.com/joncalhoun">
		  Check me out on twitter
		  <i class="fa fa-twitter" aria-hidden="true"></i>
		</a>
		<a href="https://github.com/gophercises">
		  Gophercises is on <strong>Github</strong>!
		</a>
	  </div>
	</body>
	</html>`
	expectedLinks := []Link{{"https://www.twitter.com/joncalhoun", "Check me out on twitter"}, {"https://github.com/gophercises", "Gophercises is on Github!"}}
	makeLinkTest(t, htmlString, expectedLinks)()
}
