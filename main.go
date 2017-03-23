package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/go-redis/redis"
)

type RedisCredentials struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

// struct for reading env
type VCAPServices struct {
	Redis []struct {
		Credentials RedisCredentials `json:"credentials"`
	} `json:"a9s-redis32"`
}

type KeyValue struct {
	Key   string
	Value string
}

// template store
var templates map[string]*template.Template

// fill template store
func initTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	templates["new"] = template.Must(template.ParseFiles("templates/new.html", "templates/base.html"))
}

func fetchCredentials() (RedisCredentials, error) {
	// no new read of the env var, the reason is the receiver loop
	var s VCAPServices
	err := json.Unmarshal([]byte(os.Getenv("VCAP_SERVICES")), &s)
	if err != nil {
		log.Println(err)
		return RedisCredentials{}, err
	}

	return s.Redis[0].Credentials, nil
}

func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, _ := templates[name]
	err := tmpl.ExecuteTemplate(w, template, viewModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewClient() (*redis.Client, error) {
	credentials, err := fetchCredentials()
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", credentials.Host, credentials.Port),
		Password: credentials.Password,
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	log.Printf("pong: %v ; err = %v\n", pong, err)

	return client, err
}

// create KV pair
func createKeyValue(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.PostFormValue("key")
	value := r.PostFormValue("value")

	http.Redirect(w, r, "/", 302)

	// insert key value into service
	client, err := NewClient()
	if err != nil {
		log.Printf("Failed to create connection: %v", err)
		return
	}
	err = client.Set(key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %v and value %v ; err = %v", key, value, err)
		return
	}
}

func newKeyValue(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new", "base", nil)
}

func renderKeyValues(w http.ResponseWriter, r *http.Request) {
	keyStore := make([]KeyValue, 0)

	client, err := NewClient()
	if err != nil {
		log.Printf("Failed to create connection: %v\n", err)
	} else {
		log.Printf("Collecting keys.\n")
		// collect keys
		keys, err := client.Keys("*").Result()
		if err != nil {
			log.Printf("Failed to fetch keys, err = %v\n", err)
		}
		for _, key := range keys {
			value, err := client.Get(key).Result()
			if err != nil {
				log.Printf("Failed to fetch value for key %v, err = %v\n", key, err)
			} else {
				keyStore = append(keyStore, KeyValue{Key: key, Value: value})
			}
		}
	}

	renderTemplate(w, "index", "base", keyStore)
}

func main() {
	initTemplates()

	port := "9000"
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "9000"
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir(path.Join(dir, "public")))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", renderKeyValues)
	http.HandleFunc("/key-values/new", newKeyValue)
	http.HandleFunc("/key-values/create", createKeyValue)

	log.Printf("Listening on :%v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
