package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/elastic/go-elasticsearch/client"
	"github.com/elastic/go-elasticsearch/api"
)

func init() {
}

func main() {
	gracefulStop := make(chan os.Signal, 1)
	go stop(gracefulStop)

	newIndex()
}

func newIndex() {
	c, err := client.New(client.WithHost("https://search-historic-gbpnzd-5ruj4svqcpuznwvhpsjv2jwveu.us-east-1.es.amazonaws.com:443"))
	if err != nil {
		log.Fatalf("%v connecting", err)
		return
	}

	getPayload(c)
}

func getPayload(c *client.Client) {
	resp, err := c.Get("healthcheck", "_doc", "1")
	if err != nil {
		log.Fatalf("%v got payload", err)
		return
	}
	decode(resp)
}

func decode(resp *api.GetResponse) {
	body, err := resp.DecodeBody()
	if err != nil {
		log.Fatalf("%v decoding payload", err)
		return
	}
	log.Printf("got payload: %v", body.StringToPrint())
}

func stop(gracefulStop chan os.Signal) {
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Stopping due to %s", <-gracefulStop)
	os.Exit(0)
}
