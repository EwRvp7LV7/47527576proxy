package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Replace(n *html.Node) {

	if n.Type == html.ElementNode {
		if n.Data == "a" {
			for na, a := range n.Attr {
				if a.Key == "href" {
					n.Attr[na].Val = strings.Replace(a.Val, "https://habr.com", "", -1)
				}
			}
		}
	}

	//вставлять знак "\u2122" после каждого слова длинной в 6 символов.
	if n.Type == html.TextNode && n.Parent.Data != "script" {
		re := regexp.MustCompile(`(?i)[\wа-я]{6,}`)
		n.Data = re.ReplaceAllString(n.Data, "$0\u2122")
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		Replace(child)
	}
}

type proxy struct{}

func (proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {

	log.Println(req.RemoteAddr, " ", req.Host, " ", req.Method, " ", req.URL.Path)

	//on docker it gets first (docker) port
	if !strings.Contains(req.Host, "localhost") {
		http.Error(wr, "its not habr", http.StatusBadRequest)
		return
	}

	client := &http.Client{}

	req1, _ := http.NewRequest(
		http.MethodGet,
		"https://habr.com"+req.URL.Path,
		nil,
	)

	resp, err := client.Do(req1)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}

	defer resp.Body.Close()

	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)

	re := regexp.MustCompile(`img|\Wjs|png|jpg|jpeg`)
	if re.MatchString(req.URL.Path) {
		io.Copy(wr, resp.Body)
		return
	}

	body, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	Replace(body)

	// //can send directly to wr, but need Content-Length
	// if err = html.Render(wr, body); err != nil {
	// 	log.Fatal(err)
	// }

	var buffer = new(bytes.Buffer)

	//can send directly to wr, but need Content-Length
	if err = html.Render(buffer, body); err != nil {
		log.Fatal(err)
	}

	wr.Header().Set("Content-Length", fmt.Sprint(buffer.Len()))

	wr.Write(buffer.Bytes())

}

func main() {

	handler := &proxy{}

	log.Println("Starting proxy server")
	log.Fatal("ListenAndServe:", http.ListenAndServe(":5000", handler))
	//log.Fatal("ListenAndServe:", http.ListenAndServeTLS(":443", "server.crt", "server.key", handler))

}
