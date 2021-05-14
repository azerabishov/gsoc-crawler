package main


import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
	"log"
	"sync"
	"html/template"
	"encoding/json"
	// "os"
	"strings"
)


const (
	Base_url = "https://summerofcode.withgoogle.com/archive/2020/organizations/"
	Second_base_url = "https://summerofcode.withgoogle.com/"
	url = "https://summerofcode.withgoogle.com/archive/2020/organizations/6264664972853248/"
)

var (
	tmpl = template.Must(template.ParseFiles("index.tmpl"))
	store = newStore()
)
type Organization struct{
	Name, Url, Logo string
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
					organization.Name = n.FirstChild.Data
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

	fmt.Println(organization)
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

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", index)
	mux.HandleFunc("/handle", handleOrganizationRequest)
	hs := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
	}

	err := hs.ListenAndServe()
	if err != nil {
		log.Fatal("Listen and Serve: ", err)
	}

}


func handleOrganizationRequest(w http.ResponseWriter, r *http.Request) {
	// var wg sync.WaitGroup
	fmt.Println("store.Organizations")


	// store.fetchUrls(url)
	

	// data := store.Urls[3:(len(store.Urls)-15)]

	// var resultCh =  make(chan Organization, 200)


	// wg.Add(len(data))

	// for _, i := range data {
	// 	go fetchTechnologies(i, resultCh, &wg)
	// }


	// wg.Wait()
	
	// for i := 0; i < len(data); i++ {
	// 	store.Organizations = append(store.Organizations, <-resultCh)
	// }

	fmt.Println("store.Organizations")
	// fmt.Println(store.getOrganizations())
	b, err := json.Marshal(store.getOrganizations())
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(b))

    w.Write([]byte(string(b)))
	// fmt.Fprintf(w, store.Organizations)

}

func index(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, store.getOrganizations()); err != nil {
		log.Fatal(err)
	}
}


func (s *Store) getOrganizations() []Organization{
	return s.Organizations
}