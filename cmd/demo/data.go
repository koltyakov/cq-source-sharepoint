package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/koltyakov/gosip/api"
	"github.com/schollz/progressbar/v3"
)

// Don't increase it too much as SharePoint will throttle quickly
var concurrency = 5

func seedManagers(sp *api.SP, managersNumber int) error {
	list := sp.Web().GetList("Lists/Managers")

	bar := progressbar.Default(int64(managersNumber), "Managers seeding...")
	err := runQueued(make([]int, managersNumber), concurrency, 1, func(_ int) error {
		var manager = map[string]any{
			"Title": gofakeit.Name(),
		}
		payload, _ := json.Marshal(manager)
		if _, err := list.Items().Add(payload); err != nil {
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
		return nil
	})
	_ = bar.Finish()
	return err
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
	err = runQueued(make([]int, customersNumber), concurrency, 1, func(_ int) error {
		var customer = map[string]any{
			"Title":         gofakeit.Company(),
			"RoutingNumber": gofakeit.AchRouting(),
			"Region":        regions[gofakeit.Number(0, 2)],
			"ManagerId":     managers[gofakeit.Number(0, len(managers)-1)],
			"Revenue":       gofakeit.Number(1000000, 1000000000),
		}
		payload, _ := json.Marshal(customer)
		if _, err := list.Items().Add(payload); err != nil {
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
		return nil
	})
	_ = bar.Finish()
	return err
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
	err = runQueued(make([]int, ordersNumber), concurrency, 1, func(_ int) error {
		var order = map[string]any{
			"Title":       gofakeit.AppName(),
			"CustomerId":  customers[gofakeit.Number(0, len(customers)-1)],
			"OrderNumber": gofakeit.AchAccount(),
			"OrderDate":   gofakeit.Date().Format("2006-01-02"),
			"Total":       gofakeit.Number(1000, 100000),
		}
		payload, _ := json.Marshal(order)
		if _, err := list.Items().Add(payload); err != nil {
			return fmt.Errorf("failed to create list item: %s", err)
		}
		_ = bar.Add(1)
		return nil
	})
	_ = bar.Finish()
	return err
}

// Runs a function on each item in a slice, with a maximum concurrency
func runQueued[T any](items []T, conc int, errLimit int, fn func(T) error) error {
	var errs []error
	slots := conc
	for _, item := range items {
		if len(errs) > errLimit {
			break
		}
		for slots == 0 {
			time.Sleep(10 * time.Microsecond)
		}
		slots = slots - 1
		go func(i T) {
			if err := fn(i); err != nil {
				errs = append(errs, err)
			}
			slots = slots + 1
		}(item)
	}
	for slots != conc {
		time.Sleep(10 * time.Microsecond)
	}
	return errors.Join(errs...)
}
