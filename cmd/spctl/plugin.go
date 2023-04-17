package main

import (
	"encoding/json"
	"net/http"
)

func getPluginVersion() (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.github.com/repos/koltyakov/cq-source-sharepoint/releases/latest", nil)
	if err != nil {
		return pluginVersion, err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return pluginVersion, err
	}

	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data["tag_name"].(string), nil
}
