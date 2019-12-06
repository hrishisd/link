package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Link is a href tag with its corresponding text
type Link struct {
	Href string
	Text string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func delExtraSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func textFromNode(node *html.Node) string {
	switch node.Type {
	case html.TextNode:
		return node.Data
	case html.CommentNode:
		return ""
	default:
		return textFromChildrenOf(node)
	}
}

func textFromChildrenOf(node *html.Node) string {
	res := ""
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res += " " + textFromNode(c)
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

func parseLinkFromNode(node *html.Node) (Link, error) {
	href, err := hrefFromNode(node)
	if err != nil {
		return Link{}, err
	}
	text := delExtraSpace(textFromChildrenOf(node))
	return Link{href, text}, nil
}

func traverseHTMLNode(node *html.Node) (Link, error) {
	// fmt.Println(node.Data)
	if node.Type == html.ElementNode && node.Data == "a" {
		return parseLinkFromNode(node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		link, err := traverseHTMLNode(child)
		if err == nil {
			return link, err
		}
	}
	return Link{}, errors.New("No link found in html")
}

// TraverseHTML parses links from an io.Reader
func TraverseHTML(htmlString io.Reader) (Link, error) {
	root, err := html.Parse(htmlString)
	if err != nil {
		return Link{}, err
	}
	return traverseHTMLNode(root)
}

func main() {
	f, err := os.Open("ex2.html")
	defer f.Close()
	check(err)
	link, err := TraverseHTML(f)
	fmt.Println(link)
}
