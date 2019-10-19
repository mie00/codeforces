package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Example struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Examples []Example

func getText(in *html.Node) (string, error) {
	if in.Data != "div" {
		return "", fmt.Errorf("didn't get a div, got %s instead", in.Data)
	}
	if in.LastChild == nil {
		return "", fmt.Errorf("elements has no children")
	}
	c := in.LastChild
	if c.Data != "pre" {
		return "", fmt.Errorf("didn't get a pre as a child, got %s instead", c.Data)
	}
	c = c.FirstChild
	if c.Type != html.TextNode {
		return "", fmt.Errorf("didn't get text, got %d instead", c.Type)
	}
	return strings.TrimSpace(c.Data), nil
}

func nameToUrl(name string) string {
	return fmt.Sprintf("http://codeforces.com/problemset/problem/%s/%s", name[1:], name[0:1])
}

func extract(name string) (*Examples, error) {
	url := nameToUrl(name)
	fmt.Printf("using url: %q\n", url)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("didn't expect %s as a status", resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	inputs := []string{}
	outputs := []string{}
	var f func(*html.Node)
	var retErr error
	f = func(n *html.Node) {
		if retErr != nil {
			return
		}
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					for _, cn := range strings.Split(attr.Val, " ") {
						if cn == "input" {
							txt, err := getText(n)
							if err != nil {
								retErr = err
								return
							}
							inputs = append(inputs, txt)
						}
						if cn == "output" {
							txt, err := getText(n)
							if err != nil {
								retErr = err
								return
							}
							outputs = append(outputs, txt)
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if retErr != nil {
				return
			}
		}
	}
	f(doc)
	if len(inputs) != len(outputs) {
		return nil, fmt.Errorf("error processing data, mismatching input output length")
	}
	if len(inputs) == 0 {
		return nil, fmt.Errorf("couldn't find any examples")
	}
	examples := make(Examples, len(inputs))
	for i := 0; i < len(inputs); i++ {
		examples[i].Input = inputs[i]
		examples[i].Output = outputs[i]
	}
	return &examples, nil
}
func extractString(name string) (string, error) {
	examples, err := extract(name)
	if err != nil {
		return "", err
	}
	ret, err := json.Marshal(examples)
	return string(ret), err
}
