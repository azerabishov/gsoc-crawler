package main


import (
	"fmt"
	"net/http"
	// "io/ioutil"
	"golang.org/x/net/html"
	"sync"
	"time"
	// "html/template"
)


const (
	Base_url = "https://summerofcode.withgoogle.com/archive/2020/organizations/"
	Second_base_url = "https://summerofcode.withgoogle.com/"
	url = "https://summerofcode.withgoogle.com/archive/2020/organizations/6264664972853248/"
)

type Organization struct{
	Name, Url string
	Technologies []string
}

type Store struct {
	Urls			[]string
	Organizations 	[]Organization
}


func (s *Store) fetchUrls(url string) {
	res, err := http.Get(Base_url)
	
	if err != nil{
		fmt.Println(err)
		return
	}
	
	doc, _ := html.Parse(res.Body)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					s.Urls = append(s.Urls, a.Val)
				}
			}
		}
	
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)


}



func fetchTechnologies(url string, c chan Organization, wg *sync.WaitGroup){
	res, _ := http.Get(Second_base_url + url)

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
					organization = Organization{Name: n.FirstChild.Data, Url: url}
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

func main()  {
	start := time.Now()

	store := newStore()
	
	var wg sync.WaitGroup

	store.fetchUrls(url)
	

	data := store.Urls[3:(len(store.Urls)-15)]

	var resultCh =  make(chan Organization, 200)


	wg.Add(len(data))

	for _, i := range data {
		go fetchTechnologies(i, resultCh, &wg)
	}


	wg.Wait()
	
	for i := 0; i < len(data); i++ {
		store.Organizations = append(store.Organizations, <-resultCh)
	}

	// mux := http.NewServeMux()

	// fs := http.FileServer(http.Dir("static"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fs))
	// mux.HandleFunc("/", index)

	// tmpl := template.Must(template.ParseFiles("index.tmpl"))

	// tmpl.Execute
	fmt.Println(store.Organizations)


	fmt.Println(time.Now().Sub(start))

}

// func index(w http.ResponseWriter, r *http.Request) {
// 	if err := indexTmpl.Execute(w, data.getOrganizations())
// }


// func (d *data) getOrganizations() {
	
// }