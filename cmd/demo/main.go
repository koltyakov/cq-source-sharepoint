package main

import (
	"fmt"
	"log"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	"github.com/koltyakov/gosip/auth"
)

func main() {
	authCnfg, err := auth.NewAuthFromFile("./config/private.json")
	if err != nil {
		log.Fatalf("failed to create auth config: %s", err)
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	if _, err := sp.Web().Get(); err != nil {
		log.Fatalf("failed to get web: %s", err)
	}

	// Provision lists
	for _, listModel := range listsModel {
		_ = DropList(sp, listModel.Title) // drop list if exists

		fmt.Printf("Ensuring list \"%s\"\n", listModel.Title)
		created, err := EnsureList(sp, listModel.Title, listModel.URI)
		if err != nil {
			log.Fatalf("failed to ensure list: %s", err)
		}
		if created {
			fmt.Printf("List \"%s\" was created\n", listModel.Title)
		}
		for _, fieldModel := range listModel.Fields {
			fmt.Printf("Ensuring list's \"%s\" field \"%s\"\n", listModel.Title, fieldModel.Name)
			created, err := EnsureListField(sp, listModel.Title, fieldModel.Name, fieldModel.SchemaXML)
			if err != nil {
				log.Fatalf("failed to ensure list \"%s\" field \"%s\": %s", listModel.Title, fieldModel.Name, err)
			}
			if created {
				fmt.Printf("List's \"%s\" field \"%s\" was created\n", listModel.Title, fieldModel.Name)
			}
		}
	}

	// Seed data
	if err := seedManagers(sp, 20); err != nil {
		log.Fatalf("failed to seed managers: %s", err)
	}

	if err := seedCustomers(sp, 500); err != nil {
		log.Fatalf("failed to seed customers: %s", err)
	}

	if err := seedOrders(sp, 10000); err != nil {
		log.Fatalf("failed to seed orders: %s", err)
	}

	fmt.Println("Done")
}
