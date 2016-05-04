package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
)

var addr = flag.String("addr", ":8080", "http service address")

//var homeTempl = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != "/" {
	//		http.Error(w, "Not found", 404)
	//		return
	//	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var file = r.URL.Path[1:]
	if file == "" {
		file = "home.html"
	}
	if _, err := os.Stat(file); err != nil {
		http.Error(w, "Not found", 404)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var homeTempl = template.Must(template.ParseFiles(file))
	homeTempl.Execute(w, r.Host)
}

func main() {
	flag.Parse()
	go h.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
