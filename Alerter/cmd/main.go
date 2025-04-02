package main

import (
	"middleware/example/internal/mq"
)

// github repo : https://github.com/abdelaliichou/Events-notifier-system

func main() {

	// starting Alerter consumer
	mq.StartStreamConsumer()

}
