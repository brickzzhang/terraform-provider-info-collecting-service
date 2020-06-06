package main

import (
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cloud "github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"varchecking-service/collector"
)

func GetDatasource(wr http.ResponseWriter, _ *http.Request) {
	provider := cloud.Provider().(*schema.Provider)
	resp, err := collector.SourceConfigGet(&provider.DataSourcesMap)
	if err != nil {
		wr.Write([]byte("GetDatasource error"))
		return
	}

	writeLen, err := wr.Write(resp)
	if err != nil || writeLen != len(resp) {
		wr.Write([]byte("GetDatasource error: failed to write resp"))
	}
}

func GetResource(wr http.ResponseWriter, _ *http.Request) {
	provider := cloud.Provider().(*schema.Provider)
	resp, err := collector.SourceConfigGet(&provider.ResourcesMap)
	if err != nil {
		wr.Write([]byte("GetResource error"))
		return
	}

	writeLen, err := wr.Write(resp)
	if err != nil || writeLen != len(resp) {
		wr.Write([]byte("GetResource error: failed to write resp"))
	}
}

func getDuplicatedName(queryKey string, sourceConfig []byte) (resp map[string][]string, err error) {
	var sourceMap = make(map[string]interface{})
	resp = make(map[string][]string)
	if err = json.Unmarshal(sourceConfig, &sourceMap); err != nil {
		return
	}

	sourceList := sourceMap[collector.SOURCELIST].(map[string]interface{})
	for k, v := range sourceList {
		vSlice := v.([]interface{})
		for i := 0; i<len(vSlice); i++ {
			if strings.Contains(strings.ToLower(vSlice[i].(string)), strings.ToLower(queryKey)) {
				resp[k] = append(resp[k], vSlice[i].(string))
			}
		}
	}

	return
}

func QueryName(wr http.ResponseWriter, r *http.Request) {
	// get query key from req body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		wr.Write([]byte("QueryName error: read req body error!"))
		return
	}
	var queryJson map[string]string
	if err := json.Unmarshal(body, &queryJson); err != nil {
		wr.Write([]byte("QueryName error: json unmarshal error!"))
		return
	}

	queryKey := queryJson["QueryKey"]
	if queryKey == "" {
		wr.Write([]byte("QueryName error: wrong querykey"))
		return
	}

	provider := cloud.Provider().(*schema.Provider)

	// get datasource config
	datasourceJson, err := collector.SourceConfigGet(&provider.DataSourcesMap)
	if err != nil {
		wr.Write([]byte("QueryName error: failed when getting datasource"))
		return
	}
	// get resource config
	resourceJson, err := collector.SourceConfigGet(&provider.ResourcesMap)
	if err != nil {
		wr.Write([]byte("QueryName error: failed when getting resource"))
		return
	}

	// check datasource duplicated name
	respDatesource, err := getDuplicatedName(queryKey, datasourceJson)
	if err != nil {
		wr.Write([]byte("QueryName error: failed when getDuplicatedName for Datasource"))
		return
	}
	// check resource duplicated name
	respResource, err := getDuplicatedName(queryKey, resourceJson)
	if err != nil {
		wr.Write([]byte("QueryName error: failed when getDuplicatedName for Resource"))
		return
	}

	// organize the resp
	var resp = make(map[string]interface{})
	resp[collector.DATASOURCE] = respDatesource
	resp[collector.RESOURCE] = respResource
	respJson, err := json.Marshal(resp)
	if err != nil {
		respJson = []byte("QueryName error: failed to json marshal respJson!")
	}

	writeLen, err := wr.Write(respJson)
	if err != nil || writeLen != len(respJson) {
		wr.Write([]byte("QueryName error: failed to write respJson"))
	}
}

func main() {
	http.HandleFunc("/GetDatasource", GetDatasource)
	http.HandleFunc("/GetResource", GetResource)
	http.HandleFunc("/QueryName", QueryName)
	err := http.ListenAndServe(":10001", nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}