package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/koltyakov/cq-source-sharepoint/resources/lists"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
)

var pluginVersion = "v1.6.2"

func main() {
	siteURL := getSiteURL()
	strategy := getStrategy(siteURL)
	creds, err := getCreds(strategy)
	if err != nil {
		fmt.Printf("\033[31mInvalid strategy: %s\033[0m\n", err)
		return
	}
	sp, err := checkAuth(siteURL, strategy, creds)
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}

	version, _ := getPluginVersion()
	source := getSourceName()
	destination := getDestination()

	syncScenarios := getSyncScenarios()
	for _, scenario := range syncScenarios {
		if scenario == "lists" {
			listsConf, err := getListsConf(sp)
			if err != nil {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}
			fmt.Println(listsConf)
		}
	}

	spec := &SourceSpec{
		Name:         source,
		Registry:     "github",
		Path:         "koltyakov/sharepoint",
		Version:      version,
		Destinations: []string{destination},
		Spec: PluginSpec{
			Auth: AuthSpec{
				Strategy: strategy,
				Creds:    append([][]string{{"siteUrl", siteURL}}, creds...),
			},
		},
	}

	if err := spec.Save(source + ".yml"); err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		return
	}
}

func action[T any](message string, fn func() (T, error)) (T, error) {
	fmt.Printf("\033[33m%s\033[0m", message)
	data, err := fn()
	if err != nil {
		fmt.Print("\033[2K\r")
		return data, err
	}
	fmt.Print("\033[2K\r")
	return data, nil
}

func getSiteURL() string {
	siteURLQ := &survey.Input{
		Message: "SharePoint URL:",
		Help:    "Site absolute URL, e.g. https://contoso.sharepoint.com/sites/MySite",
	}

	var siteURL string
	_ = survey.AskOne(siteURLQ, &siteURL, survey.WithValidator(shouldBeURL))
	return siteURL
}

func getStrategy(siteURL string) string {
	strats, err := action("Resolving auth strategy...", func() ([]string, error) {
		return getStrategies(siteURL)
	})
	if err != nil {
		fmt.Printf("\033[31mError: %s\033[0m\n", err)
		strats = allStrats
	}

	strategyQ := &survey.Select{
		Message: "Auth method:",
		Options: strats,
		Help:    "See more at https://go.spflow.com/auth/overview",
		Description: func(value string, index int) string {
			return stratsConf[value].Desc
		},
	}

	var strategy string
	_ = survey.AskOne(strategyQ, &strategy)

	return strategy
}

func getCreds(strategy string) ([][]string, error) {
	s, ok := stratsConf[strategy]
	if !ok {
		return nil, fmt.Errorf("can't resolve strategy %s", strategy)
	}
	return s.Creds(), nil
}

func checkAuth(siteURL, strategy string, creds [][]string) (*api.SP, error) {
	auth, err := newAuthByStrategy(strategy)
	if err != nil {
		return nil, err
	}

	cnfg := map[string]string{"siteURL": siteURL}
	for _, c := range creds {
		cnfg[c[0]] = c[1]
	}
	credsBytes, _ := json.Marshal(cnfg)

	if err := auth.ParseConfig(credsBytes); err != nil {
		return nil, err
	}

	client := &gosip.SPClient{AuthCnfg: auth}
	sp := api.NewSP(client)

	web, err := action("Reaching site, checking auth...", sp.Web().Get)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\033[32mSuccess! Site title: \"%s\"\033[0m\n", web.Data().Title)

	return sp, nil
}

func getSourceName() string {
	var sourceName string
	sourceNameQ := &survey.Input{
		Message: "Source name:",
		Default: "sharepoint",
		Help:    "Source name to be used in the config file",
	}
	_ = survey.AskOne(sourceNameQ, &sourceName, survey.WithValidator(survey.Required))
	return sourceName
}

func getDestination() string {
	var destination string
	destinationNameQ := &survey.Input{
		Message: "Destination name:",
		Default: "postgres",
		Help:    "Destination name to be used in the config file",
	}
	_ = survey.AskOne(destinationNameQ, &destination, survey.WithValidator(survey.Required))
	return destination
}

var syncScenariosMap = map[string]string{
	"Lists and libraries":    "lists",
	"Content types rollup":   "content_types",
	"Search driven queries":  "search",
	"Managed metadata terms": "mmd",
	"User profiles (UPS)":    "profiles",
}

func getSyncScenarios() []string {
	syncScenariosQ := &survey.MultiSelect{
		Message: "Select subjects of sync:",
		Options: []string{
			"Lists and libraries",
			"Content types rollup",
			"Search driven queries",
			"Managed metadata terms",
			"User profiles (UPS)",
		},
	}

	var syncScenarios []string
	_ = survey.AskOne(syncScenariosQ, &syncScenarios, survey.WithValidator(survey.Required))

	for i, s := range syncScenarios {
		syncScenarios[i] = syncScenariosMap[s]
	}

	return syncScenarios
}

type ListConf struct {
	ID   string
	Spec lists.Spec
}

type listInfo struct {
	ID         string `json:"Id"`
	Title      string `json:"Title"`
	RootFolder struct {
		URL string `json:"ServerRelativeUrl"`
	} `json:"RootFolder"`
}

func getListsConf(sp *api.SP) ([]ListConf, error) {
	resp, err := action("Getting lists", func() (api.ListsResp, error) {
		return sp.Web().Lists().
			Select("Id,Title,RootFolder/ServerRelativeUrl").
			Expand("RootFolder").Get()
	})
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(sp.ToURL())
	basePath := u.Path + "/"

	data := resp.Data()
	ll := make([]string, len(data))
	llMap := make(map[string]listInfo)
	for i, l := range data {
		info := listInfo{}
		_ = json.Unmarshal(l.Normalized(), &info)
		listURI := strings.Replace(info.RootFolder.URL, basePath, "", 1)

		listKey := info.Title + " \033[90m(" + listURI + ")\033[0m"
		llMap[listKey] = info
		ll[i] = listKey
	}

	var listsToSync []string
	listsQ := &survey.MultiSelect{
		Message: "Select lists to sync:",
		Options: ll,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}
	_ = survey.AskOne(listsQ, &listsToSync, survey.WithValidator(survey.Required))

	for _, l := range listsToSync {
		info := llMap[l]
		if err := getListFieldsConf(sp, info.Title, info.RootFolder.URL); err != nil {
			return nil, err
		}
	}

	var listsConf []ListConf
	return listsConf, nil
}

func getListFieldsConf(sp *api.SP, name string, listURI string) error {
	resp, err := action("Getting fields for "+name, func() (api.FieldsResp, error) {
		return sp.Web().GetList(listURI).
			Fields().
			Filter("Hidden eq false and FieldTypeKind ne 12").
			Get()
	})
	if err != nil {
		return err
	}

	data := resp.Data()
	ignoreFields := []string{"AppAuthor", "AppEditor"}
	dd := []api.FieldResp{}
	for _, f := range data {
		if f.Data().TypeAsString == "Lookup" && f.Data().LookupList == "" {
			continue
		}

		if includes(ignoreFields, f.Data().EntityPropertyName) {
			continue
		}

		dd = append(dd, f)
	}

	fields := make([]string, len(dd))
	for i, f := range dd {
		fields[i] = f.Data().Title +
			" \033[90m(Prop: " + f.Data().EntityPropertyName +
			", Type: " + f.Data().TypeAsString + ")\033[0m"
	}

	var fieldsToSync []string
	fieldsQ := &survey.MultiSelect{
		Message: "Select fields to sync for " + name + ":",
		Options: fields,
		Filter: func(filter string, value string, index int) bool {
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}
	_ = survey.AskOne(fieldsQ, &fieldsToSync, survey.WithValidator(survey.Required))

	return nil
}

func includes(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
