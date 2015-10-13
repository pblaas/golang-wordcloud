package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"html/template"
	"net/http"
	"path"
)

type Wordcount struct {
	Word  string `bson:"_id"`
	Value int    `bson:"value"`
}

func wcloudpage(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("mgodb").C("word_count")

	query := c.Find(nil)
	var words []Wordcount
	if err := query.All(&words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "wcloud.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, words); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", HomeHandler)

	channels := r.Path("/channels").Subrouter()
	channels.Methods("GET").HandlerFunc(ChannelIndexHandler)

	channeldate := r.PathPrefix("/channels/{id}/{date}").Subrouter()
	channeldate.Methods("GET").HandlerFunc(ChannelShowHandlerDate)

	channel := r.PathPrefix("/channels/{id}").Subrouter()
	channel.Methods("GET").HandlerFunc(ChannelShowHandler)

	keyword := r.PathPrefix("/twitter/{keyword}").Subrouter()
	keyword.Methods("GET").HandlerFunc(TwitterKeywordShowHandler)

	fmt.Println("Wordcloud Server started!")
	http.ListenAndServe(":3002", r)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Basic information on the app.")
}

func ChannelIndexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Listing of wordcloud channels.")
}

func ChannelShowHandler(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("irclogs").C(id)

	query := c.Find(nil)
	var words []Wordcount
	if err := query.All(&words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "wcloud.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(rw, words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func TwitterKeywordShowHandler(rw http.ResponseWriter, r *http.Request) {
	keyword := mux.Vars(r)["keyword"]

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(keyword).C("word_count")

	query := c.Find(nil).Sort("-value").Limit(75)
	var words []Wordcount
	if err := query.All(&words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "wcloud.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(rw, words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func ChannelShowHandlerDate(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	date := mux.Vars(r)["date"]
	iddate := fmt.Sprint(id, "_", date)
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	//Optional. Switch the session to monotonic behavior
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("irclogs").C(iddate)

	query := c.Find(nil)
	var words []Wordcount
	if err := query.All(&words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "wcloud.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(rw, words); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}
