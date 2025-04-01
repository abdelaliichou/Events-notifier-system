package main

import (
	"middleware/example/internal/mq"
)

func main() {

	// starting Alerter consumer
	mq.StartStreamConsumer()

}
