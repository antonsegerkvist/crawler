package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type edge struct {
	Depth  int
	Source string
	Target string
}

var maxDepth = 3
var visitedLinks = map[string]bool{}
var linkMap = map[string]*[]edge{}

func main() {
	performIteration("https://www.youtube.com/", 0)
	for key, val := range linkMap {
		for _, e := range *val {
			fmt.Println("---")
			fmt.Printf("source: %s\n", key)
			fmt.Printf("depth:  %d\n", e.Depth)
			fmt.Printf("target: %s\n", e.Target)
		}
	}
}

func performIteration(target string, depth int) {
	if depth >= maxDepth || visitedLinks[target] == true {
		return
	}
	visitedLinks[target] = true

	httpClient := http.Client{}

	response, err := httpClient.Get(target)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	defer response.Body.Close()

	node, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	links := &[]string{}
	traverseTree(node, links)

	for _, v := range *links {
		if e, ok := linkMap[target]; ok {
			*e = append(*e, edge{
				Depth:  depth,
				Source: target,
				Target: v,
			})
		} else {
			linkMap[target] = &[]edge{
				edge{
					Depth:  depth,
					Source: target,
					Target: v,
				},
			}
		}

		performIteration(v, depth+1)
	}
}

func traverseTree(node *html.Node, links *[]string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, v := range node.Attr {
			if v.Key == "href" {
				*links = append(*links, v.Val)
			}
		}
	}
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		traverseTree(next, links)
	}
}
