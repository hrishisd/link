package link

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link is a href tag with its corresponding text
type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML document and return a
// slice of links parsed from it
func Parse(htmlReader io.Reader) ([]Link, error) {
	docRoot, err := html.Parse(htmlReader)
	if err != nil {
		return nil, err
	}
	return parseHTMLRec(docRoot), nil
}

func parseHTMLRec(node *html.Node) []Link {
	var links []Link
	if node.Type == html.ElementNode && node.Data == "a" {
		link, err := parseLinkFromNode(node)
		if err == nil {
			links = append(links, link)
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		childLinks := parseHTMLRec(child)
		links = append(links, childLinks...)
	}
	return links
}

func parseLinkFromNode(node *html.Node) (Link, error) {
	href, err := hrefFromNode(node)
	if err != nil {
		return Link{}, err
	}
	text := delExtraSpace(textFromChildrenOf(node))
	return Link{href, text}, nil
}

func textFromNode(node *html.Node) string {
	switch node.Type {
	case html.TextNode:
		return node.Data
	case html.ElementNode:
		return textFromChildrenOf(node)
	default:
		return ""
	}
}

func textFromChildrenOf(node *html.Node) string {
	res := ""
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res += textFromNode(c)
	}
	return res
}

func hrefFromNode(node *html.Node) (string, error) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("Didn't find href in node")
}

func delExtraSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")

}
