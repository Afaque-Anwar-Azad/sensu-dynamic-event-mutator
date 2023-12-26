package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the mutator plugin config.
type Config struct {
	sensu.PluginConfig
	ApiUrl         string
	ApiKey         string
	Labels         map[string]string
	Annotations    map[string]string
	AddLabels      bool
	AddAnnotations bool
	AddAll         bool
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-dynamic-event-mutator",
			Short:    "mutator for dynamically updating sensu event labels and annotations based on sensu entity",
			Keyspace: "sensu.io/plugins/sensu-dynamic-event-mutator/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "api-url",
			Env:       "SENSU_API_URL",
			Argument:  "api-url",
			Shorthand: "a",
			Default:   "http://127.0.0.1:8080",
			Usage:     "Sensu API URL",
			Value:     &plugin.ApiUrl,
		},
		{
			Path:      "api-key",
			Env:       "SENSU_API_KEY",
			Argument:  "api-key",
			Shorthand: "k",
			Default:   "",
			Secret:    true,
			Usage:     "Sensu API Key",
			Value:     &plugin.ApiKey,
		},
		{
			Path:      "",
			Env:       "",
			Argument:  "add-labels",
			Shorthand: "",
			Default:   false,
			Usage:     "Mutates event.Entity.Labels based on entity.Labels",
			Value:     &plugin.AddLabels,
		},
		{
			Path:      "",
			Env:       "",
			Argument:  "add-annotations",
			Shorthand: "",
			Default:   false,
			Usage:     "Mutates event.Entity.Annotations based on entity.Annotations",
			Value:     &plugin.AddAnnotations,
		},
		{
			Path:      "",
			Env:       "",
			Argument:  "add-all",
			Shorthand: "",
			Default:   false,
			Usage:     "Mutates event.Entity.Labels based on entity.Labels and event.Entity.Annotations based on entity.Annotations",
			Value:     &plugin.AddAll,
		},
	}
)

func main() {
	mutator := sensu.NewGoMutator(&plugin.PluginConfig, options, checkArgs, executeMutator)
	mutator.Execute()
}

func checkArgs(event *types.Event) error {
	if len(os.Getenv("SENSU_API_URL")) > 0 {
		plugin.ApiUrl = os.Getenv("SENSU_API_URL")
	}
	if len(os.Getenv("SENSU_API_KEY")) > 0 {
		plugin.ApiKey = os.Getenv("SENSU_API_KEY")
	}
	if len(plugin.ApiKey) == 0 || len(plugin.ApiUrl) == 0 {
		return fmt.Errorf("--api-url and --api-key or enviornment variables $SENSU_API_URL and $SENSU_API_KEY are required!")
	}

	entityMetadata := getEntityMetadata(event)

	if plugin.AddLabels {
		fmt.Printf("Adding labels from \"entity.Labels\"\n")
		plugin.Labels = addLabels(entityMetadata)
	}

	if plugin.AddAnnotations {
		fmt.Printf("Adding annotations from \"entity.Annotations\"\n")
		plugin.Annotations = addAnnotations(entityMetadata)
	}

	if plugin.AddAll {
		fmt.Printf("Adding labels from \"entity.Labels\" and Adding annotations from \"entity.Annotations\"\n")
		plugin.Labels = addLabels(entityMetadata)
		plugin.Annotations = addAnnotations(entityMetadata)
	}
	return nil
}

func handleError(message string, err error) {
	fmt.Printf("%s: %s\n", message, err)
	os.Exit(1)
}

func getEntityMetadata(event *types.Event) map[string]interface{} {

	namespace, entityName := "default", event.Entity.Name
	url := fmt.Sprintf("%s/api/core/v2/namespaces/%s/entities/%s", os.Getenv("SENSU_API_URL"), namespace, entityName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		handleError("Error creating request", err)
	}
	req.Header.Add("Authorization", "Key "+os.Getenv("SENSU_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		handleError("Error making request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError("Error reading response body", err)
	}

	var entityInfo map[string]json.RawMessage
	err = json.Unmarshal(body, &entityInfo)
	if err != nil {
		handleError("Error parsing JSON response", err)
	}

	var entityMetadata map[string]interface{}
	err = json.Unmarshal(entityInfo["metadata"], &entityMetadata)
	if err != nil {
		handleError("Error extracting metadata from response", err)
	}
	return entityMetadata
}

func addLabels(entityMetadata map[string]interface{}) map[string]string {

	labels, ok := entityMetadata["labels"].(map[string]interface{})
	if !ok {
		handleError("Error extracting labels from metadata", fmt.Errorf("labels field not found"))
	}
	jsonLabels, err := json.Marshal(labels)
	if err != nil {
		handleError("Error encoding labels to JSON", err)
	}

	var entityLabels map[string]string
	err = json.Unmarshal([]byte(jsonLabels), &entityLabels)
	if err != nil {
		handleError("Error decoding labels from JSON", err)
	}
	return entityLabels
}

func addAnnotations(entityMetadata map[string]interface{}) map[string]string {

	annotations, ok := entityMetadata["annotations"].(map[string]interface{})
	if !ok {
		handleError("Error extracting Annotations from metadata", fmt.Errorf("labels field not found"))
	}
	jsonAnnotations, err := json.Marshal(annotations)
	if err != nil {
		handleError("Error encoding Annotations to JSON", err)
	}
	var entityAnnotations map[string]string
	err = json.Unmarshal([]byte(jsonAnnotations), &entityAnnotations)
	if err != nil {
		handleError("Error decoding Annotations from JSON", err)
	}
	return entityAnnotations
}

func executeMutator(event *types.Event) (*types.Event, error) {
	if plugin.AddLabels {
		event.Entity.Labels = plugin.Labels
	}

	if plugin.AddAnnotations {
		event.Entity.Annotations = plugin.Annotations
	}

	if plugin.AddAll {
		event.Entity.Labels = plugin.Labels
		event.Entity.Annotations = plugin.Annotations
	}

	return event, nil
}
