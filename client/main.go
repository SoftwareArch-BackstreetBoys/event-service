package main

import (
	"client/configs"
	"client/services"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	creds := insecure.NewCredentials()
	cc, err := grpc.Dial(configs.EnvPort(), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	eventClient := services.NewEventServiceClient(cc)
	eventService := services.NewEventService(eventClient)

	// err = eventService.CreateEvent("Test3", "this is the second test", "when", "where", 99, 1, "club_id string", "ping")
	// err = eventService.GetEvent("66e5df38d30fd5c1bf9fb829")
	// err = eventService.DeleteEvent("66e5df38d30fd5c1bf9fb829")
	// err = eventService.UpdateEvent("66e5dff1d30fd5c1bf9fb82a", "Test2", "this is the first test, no it is updated!", "when", "where", 99, 1, "club_id string")
	// err = eventService.ListEvents()
	// err = eventService.JoinEvent("66e5dff1d30fd5c1bf9fb82a", "65e4dff1d20fd5c1bf0fb55a")
	err = eventService.LeaveEvent("66e5dff1d30fd5c1bf9fb82a", "65e4dff1d20fd5c1bf0fb55a")
	if err != nil {
		log.Fatal(err)
	}
}
