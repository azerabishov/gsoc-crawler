package main


import (
	"fmt"
	"net/http"
	// "io/ioutil"
	"golang.org/x/net/html"
	"sync"
	"time"
)


const (
	Base_url = "https://summerofcode.withgoogle.com/archive/2020/organizations/"
	Second_base_url = "https://summerofcode.withgoogle.com/"
	url = "https://summerofcode.withgoogle.com/archive/2020/organizations/6264664972853248/"
)

type Organizations struct{
	Name, Url string
	Technologies []string
}



func getOrganizatiosUrl(url string)  ([]string, error){
	var urls []string

	res, err := http.Get(Base_url)
	
	if err != nil {
		return nil, err
	}

	doc, _ := html.Parse(res.Body)

	getUrl(doc, &urls)

	return urls, nil

}


func getUrl(n *html.Node, urls *[]string){
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				*urls = append(*urls, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getUrl(c, urls)
	}

}

func getTechnology(url string, c chan Organizations, wg *sync.WaitGroup){
	res, _ := http.Get(Second_base_url + url)

	doc, _ := html.Parse(res.Body)

	var umumi Organizations
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Val == "organization__tag organization__tag--technology" {
					umumi.Technologies = append(umumi.Technologies, n.FirstChild.Data)
					break
				}
			}
		}else if n.Type == html.ElementNode && n.Data == "h3" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					umumi = Organizations{Name: n.FirstChild.Data, Url: url}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	c<-umumi
	wg.Done()

	// return umumi, nil
}


func main()  {

	start := time.Now()
	var result []Organizations
	var wg sync.WaitGroup

	data, err := getOrganizatiosUrl(url)
	
	if err != nil{
		fmt.Println(err)
		return
	}
	data = data[3:(len(data)-15)]

	var resultCh =  make(chan Organizations, 200)


	wg.Add(len(data))

	for _, i := range data {
		go getTechnology(i, resultCh, &wg)
	}



	wg.Wait()
	
	for i := 0; i < len(data); i++ {
		result = append(result, <-resultCh)
	}


	fmt.Println(result)


	fmt.Println(time.Now().Sub(start))

}