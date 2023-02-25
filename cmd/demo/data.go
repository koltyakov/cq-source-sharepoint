package main

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/koltyakov/gosip/api"
	"github.com/schollz/progressbar/v3"
)

func seedManagers(sp *api.SP, managersNumber int) error {
	list := sp.Web().GetList("Lists/Managers")

	bar := progressbar.Default(int64(managersNumber), "Managers seeding...")
	for i := 0; i < managersNumber; i++ {
		var manager = map[string]any{
			"Title": gofakeit.Name(),
		}
		payload, _ := json.Marshal(manager)
		if _, err := list.Items().Add(payload); err != nil {
			_ = bar.Finish()
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
	}
	_ = bar.Finish()

	return nil
}

func seedCustomers(sp *api.SP, customersNumber int) error {
	manResp, err := sp.Web().GetList("Lists/Managers").Items().Top(5000).Select("ID").Get()
	if err != nil {
		return fmt.Errorf("failed to get managers: %s", err)
	}
	var managers = make([]int, len(manResp.Data()))
	for i, manager := range manResp.Data() {
		managers[i] = manager.Data().ID
	}

	regions := []string{"AMER", "EMEA", "APAC"}

	list := sp.Web().GetList("Lists/Customers")

	bar := progressbar.Default(int64(customersNumber), "Customers seeding...")
	for i := 0; i < customersNumber; i++ {
		var customer = map[string]any{
			"Title":         gofakeit.Company(),
			"RoutingNumber": gofakeit.AchRouting(),
			"Region":        regions[gofakeit.Number(0, 2)],
			"ManagerId":     managers[gofakeit.Number(0, len(managers)-1)],
			"Revenue":       gofakeit.Number(1000000, 1000000000),
		}
		payload, _ := json.Marshal(customer)
		if _, err := list.Items().Add(payload); err != nil {
			_ = bar.Finish()
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
	}
	_ = bar.Finish()

	return nil
}

func seedOrders(sp *api.SP, ordersNumber int) error {
	custResp, err := sp.Web().GetList("Lists/Customers").Items().Top(5000).Select("ID").Get()
	if err != nil {
		return fmt.Errorf("failed to get customers: %s", err)
	}
	var customers = make([]int, len(custResp.Data()))
	for i, customer := range custResp.Data() {
		customers[i] = customer.Data().ID
	}

	list := sp.Web().GetList("Lists/Orders")

	bar := progressbar.Default(int64(ordersNumber), "Orders seeding...")
	for i := 0; i < ordersNumber; i++ {
		var order = map[string]any{
			"Title":       gofakeit.AppName(),
			"CustomerId":  customers[gofakeit.Number(0, len(customers)-1)],
			"OrderNumber": gofakeit.AchAccount(),
			"OrderDate":   gofakeit.Date().Format("2006-01-02"),
			"Total":       gofakeit.Number(1000, 100000),
		}
		payload, _ := json.Marshal(order)
		if _, err := list.Items().Add(payload); err != nil {
			_ = bar.Finish()
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
	}
	_ = bar.Finish()

	return nil
}
