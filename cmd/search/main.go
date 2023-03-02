package main

import (
	"encoding/json"
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

	results, err := sp.Search().PostQuery(&api.SearchQuery{
		QueryText: "*",
		RowLimit:  1,
	})

	if err != nil {
		log.Fatalf("failed to get web: %s", err)
	}

	j, _ := json.MarshalIndent(results.Data().PrimaryQueryResult.RelevantResults.Table.Rows, "", "  ")
	log.Println(string(j))
}
