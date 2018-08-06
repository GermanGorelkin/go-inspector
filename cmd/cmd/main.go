package main

import (
	"github.com/germangorelkin/go-inspector/inspector"
	"log"
	"net/url"
	"os"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is not set.")
	}
	instance := os.Getenv("INSTANCE")
	if instance == "" {
		log.Fatal("INSTANCE is not set.")
	}

	log.Print(apiKey)
	log.Print(instance)
	//return

	inst, err := url.Parse(instance)
	if err != nil {
		log.Panic(err)
	}
	cfg := inspector.ClintConf{
		APIKey:   apiKey,
		Instance: inst,
	}
	c := inspector.NewClient(cfg)

	f, err := os.Open("C://Users//gg//Desktop//hotfield_4e9a6589b0c0ad6a1533543806228.jpg")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	u, err := c.Image.Uploads(f, f.Name())
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v", u)
}
