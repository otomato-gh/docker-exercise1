// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Post - Our struct for all posts
type Post struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type RedisInstance struct {
	RInstance *redis.Client
}

var Posts []Post

func initDB() *RedisInstance {
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_addr := redis_host + ":" + redis_port
	client := redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Can't connect to redis- %s", err)
	}

	for _, p := range Posts {
		json, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
		}
		client.Set(p.Id, json, 0).Err()
		if err != nil {
			log.Println(err)
		}
	}
	return &RedisInstance{RInstance: client}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My Blog!")
	log.Println("Endpoint Hit: homePage")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
	log.Println("healthcheck")
}

func (c *RedisInstance) returnAllPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: returnAllPosts")

	var allPosts []Post
	iter := c.RInstance.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		var post Post
		log.Println("Got " + iter.Val())
		bytes, _ := c.RInstance.Get(iter.Val()).Bytes()
		log.Println("bytes " + string(bytes))
		err := json.Unmarshal(bytes, &post)
		if err != nil {
			log.Fatalf("Error getting posts - %s", err)
		}
		allPosts = append(allPosts, post)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(allPosts)
}

func (c *RedisInstance) returnSinglePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var post Post
	bytes, _ := c.RInstance.Get(key).Bytes()
	err := json.Unmarshal(bytes, &post)
	if err != nil {
		log.Fatalf("Error getting posts - %s", err)
	}

	json.NewEncoder(w).Encode(post)

}

func (c *RedisInstance) createNewPost(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Post struct
	// append this to our Posts array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Post
	json.Unmarshal(reqBody, &post)
	// update our global Posts array to include
	// our new Post
	myjson, err := json.Marshal(post)
	if err != nil {
		log.Println(err)
	}
	err = c.RInstance.Set(post.Id, myjson, 0).Err()
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(post)
}

func (c *RedisInstance) deletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := c.RInstance.Del(id).Err()
	if err != nil {
		log.Println(err)
	}

}

func handleRequests() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. env PORT=8081 go run main.go")
	}

	redisHandler := initDB()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/healthz", healthz)
	myRouter.HandleFunc("/posts", redisHandler.returnAllPosts)
	myRouter.HandleFunc("/post", redisHandler.createNewPost).Methods("POST")
	myRouter.HandleFunc("/post/{id}", redisHandler.deletePost).Methods("DELETE")
	myRouter.HandleFunc("/post/{id}", redisHandler.returnSinglePost)
	log.Printf("App listening on port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func main() {
	Posts = []Post{
		Post{Id: "1", Title: "Hello", Desc: "Post Description", Content: "Post Content"},
		Post{Id: "2", Title: "Hello 2", Desc: "Post Description", Content: "Post Content"},
	}

	handleRequests()
}
