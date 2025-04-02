package main

import (
	"Scheduler/models"
	"Scheduler/mq"
	"Scheduler/webservice"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhashkevych/scheduler"
	"os"
	"os/signal"
	"time"
)

// github repo : https://github.com/abdelaliichou/Events-notifier-system

func main() {

	// Starting NATS JetStream
	mq.InitStream()

	schedulingFunctionCall()

}

func schedulingFunctionCall() {

	ctx := context.Background()
	sc := scheduler.NewScheduler()
	sc.Add(ctx, FetchingFromConfig, time.Second*10)

	// cette partie sert à maintenir le programme en vie, tant qu'il n'a pas reçu de signal os.Interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	sc.Stop()
}

func FetchingFromConfig(_ context.Context) {

	// Make the HTTP GET request
	var body []byte
	body = webservice.HttpRequest(models.CONFIG_PATH, false)

	// Parse the JSON data into a slice of Resource structs
	var resources []models.Resource
	err := json.Unmarshal(body, &resources)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed data
	fmt.Println("\nGetting all resources from Config : ", models.CONFIG_PATH)
	for _, resource := range resources {
		fmt.Printf("ID: %s, UcaID: %d, Name: %s\n", resource.Id, resource.UcaID, resource.Name)
	}

	// fetching the data from the uca server
	FetchingFromUCA(resources, false)

}

func FetchingFromUCA(resources []models.Resource, show bool) {

	if show {
		fmt.Println("\nURL with all resourceIDs : ", models.UCA_URL("8", resources))
	}

	var allEvents []models.Event
	for _, resource := range resources {

		var customResources []models.Resource
		customResources = append(customResources, resource)
		url := models.UCA_URL("8", customResources)
		if show {
			fmt.Printf("\nURL of resource %d id : %s\n", resource.UcaID, url)
		}

		// doing the request with this particular resourceID
		ucaResp := webservice.HttpRequest(url, false)
		specificResourceEvents := models.ParsingEvents(ucaResp, resource.Id, false)
		allEvents = append(allEvents, specificResourceEvents...)

	}

	if len(allEvents) == 0 {
		fmt.Println("NO EVENTS FROM UCA !")
		return
	}

	fmt.Println("ALL EVENTS FROM UCA :")
	models.DisplayEvents(allEvents)

	// Group events by UID & merge resource IDs
	groupedEvents := models.GroupEventsByUID(allEvents)

	fmt.Println("ALL EVENTS FROM UCA AFTER GROUPING :")
	models.DisplayEvents(groupedEvents)

	// Send into MQ
	mq.SendEventsToMQ(groupedEvents)
}
