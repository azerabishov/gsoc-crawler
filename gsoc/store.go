package gsoc


import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"sync"
	"strings"
)


const baseUrl = "https://summerofcode.withgoogle.com/"

var (
	store = newStore()
)


type Organization struct{
	Title, Url, Logo string
	Technologies []string
}


type Store struct {
	Urls			[]string
	Organizations 	[]Organization
}



func FetchUrls(url string)  ([]string, error){
	res, err := http.Get(url)
	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	doc, _ := html.Parse(res.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if(n.Attr[0].Val == "organization-card__link"){
						store.Urls = append(store.Urls, a.Val)
					}
				}
			}
		}

	
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return store.Urls, nil;
}



func FetchTechnologies(url string, c chan Organization, wg *sync.WaitGroup){
	res, _ := http.Get(baseUrl + url)

	doc, _ := html.Parse(res.Body)
	var organization Organization
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Val == "organization__tag organization__tag--technology" {
					organization.Technologies = append(organization.Technologies, n.FirstChild.Data)
					break
				}
			}
		}else if n.Type == html.ElementNode && n.Data == "h3" {
			for _, a := range n.Attr {
				if a.Key == "class" {
					organization.Title = n.FirstChild.Data
					organization.Url = url
				}
			}
		}else if n.Type == html.ElementNode && n.Data == "org-logo" {
			for _, a := range n.Attr {
				if a.Key == "data" {
					organization.Logo = strings.Split(a.Val, "'")[3];

				}
			}
		}
		

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	c<-organization
	wg.Done()
}




func newStore() *Store {
	return &Store{
		Urls: []string{},
		Organizations: []Organization{},
		}
}



func (s *Store) GetOrganizations() []Organization{
	return s.Organizations
}