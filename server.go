package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	
	"sync"
	
	"time"
)

var (
	urlmap = make(map[string]string)
	mu sync.Mutex
)

func main(){
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/",homeHandler)
	http.HandleFunc("/shorten",shortenHandler)
	http.HandleFunc("/r/",redirectHandler)
	fmt.Println("server run http://localhost:8080")
	http.ListenAndServe(":8080",nil)


	
}

func homeHandler(w http.ResponseWriter,r *http.Request){
	tmpl := template.Must(template.ParseFiles("index.html"))

	mu.Lock()
	defer mu.Unlock()
	tmpl.Execute(w,urlmap)
}

func shortenHandler(w http.ResponseWriter,r *http.Request){
	if r.Method != "POST" {
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}

	r.ParseForm()
	originalURL := r.FormValue("url")
	code := randomString(6)

	mu.Lock()
	urlmap[code] = originalURL
	mu.Unlock()

	http.Redirect(w,r,"/",http.StatusSeeOther)
}
func redirectHandler(w http.ResponseWriter,r *http.Request){

	code := r.URL.Path[len("/r/"):]

	mu.Lock()
	url,ok := urlmap[code]
	mu.Unlock()

	if ok {
		http.Redirect(w,r,url,http.StatusSeeOther)
	}else{
		http.NotFound(w,r)
	}
}

func randomString(n int ) string{
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune,n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
		
	}
	return string(b)
}