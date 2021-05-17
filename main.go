package main


import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"html/template"
	"encoding/json"
	"time"
	"gsoc-crawler/gsoc"
)


const (
	baseUrl = "https://summerofcode.withgoogle.com/"
)

var (
	tmpl = template.Must(template.ParseFiles("index.tmpl"))
	archives = []Archives{
		{Title: "GSoC 2020", Url: "archive/2020/organizations/"},
		{Title: "GSoC 2019", Url: "archive/2019/organizations/"},
		{Title: "GSoC 2018", Url: "archive/2018/organizations/"},
		{Title: "GSoC 2017", Url: "archive/2017/organizations/"},
		{Title: "GSoC 2016", Url: "archive/2016/organizations/"},
	}
)

type Archives struct {
	Title, Url string
}


type Url struct {
	Year string
}



func main()  {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/index", index)
	mux.HandleFunc("/", indexPage)
	mux.HandleFunc("/organizations", organizationHandler)

    err := http.ListenAndServe(":9090", mux)
	if err != nil {
		log.Fatal("Listen and Serve: ", err)
	}

}

func indexPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, archives); err != nil {
		log.Fatal(err)
	}
}


func organizationHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var store gsoc.Store
	var wg sync.WaitGroup

	
    decoder := json.NewDecoder(r.Body)
    var u Url
    err := decoder.Decode(&u)
    if err != nil {
        log.Println(err)
		return 
    }


	store.Urls, err = gsoc.FetchUrls(baseUrl + u.Year)

	if err != nil {
        log.Println(err)
		return 
    }


	var resultCh =  make(chan gsoc.Organization, len(store.Urls))

	wg.Add(len(store.Urls))

	for _, i := range store.Urls {
		go gsoc.FetchTechnologies(i, resultCh, &wg)
	}

	wg.Wait()
	
	for i := 0; i < len(store.Urls); i++ {
		store.Organizations = append(store.Organizations, <-resultCh)
	}

	b, err := json.Marshal(store.GetOrganizations())
    if err != nil {
        log.Println(err)
        return
    }
	fmt.Println(time.Now().Sub(start))

    w.Write([]byte(string(b)))
}



func index(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, archives); err != nil {
		log.Fatal(err)
		return
	}
}

