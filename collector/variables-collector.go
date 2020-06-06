package collector

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cloud "github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud"
	"os"
	"regexp"
)

func getKeys(m map[string]*schema.Schema) []string {
	i := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func sourceConfigGet(source *map[string]*schema.Resource) (count int, jsonStr []byte, err error) {
	count = 0
	sourceRegexp, _ := regexp.Compile(`^tencentcloud_([\w_]*)$`)
	var jsonMap = make(map[string][]string)
	for k, v := range *source {
		sourceName := sourceRegexp.FindStringSubmatch(k)
		jsonMap[sourceName[1]] = getKeys(v.Schema)
		count++
	}

	jsonStr, err = json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("json marshal error!")
	}

	return
}

func write2file(content []byte, filepath string) (err error) {
	sourceFile, err := os.OpenFile("./current-configs/"+filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	if _, err = sourceFile.Write(content); err != nil {
		return err
	}
	return nil
}

func Run() {
	provider := cloud.Provider().(*schema.Provider)
	fmt.Println("==== collecting datasource config ====")
	datasourceCount, datasourceJson, err := sourceConfigGet(&provider.DataSourcesMap)
	if err != nil {
		fmt.Println("[CRITAL] Collecting datasource config error: ", err)
	} else {
		datasourceFile := "datasourceConfig.json"
		if err := write2file(datasourceJson, datasourceFile); err != nil {
			fmt.Println("[CRITAL] Writing datasource config file error: ", err)
		}
		println("[DEBUG] Total number of datasource is: ", datasourceCount)
	}

	fmt.Println("==== collecting resource config ====")
	resourceCount, resourceJson, err := sourceConfigGet(&provider.ResourcesMap)
	if err != nil {
		fmt.Println("[CRITAL] Collecting resource config error: ", err)
	} else {
		resourceFile := "resourceConfig.json"
		if err := write2file(resourceJson, resourceFile); err != nil {
			fmt.Println("[CRITAL] Writing resource config file error: ", err)
		}
		println("[DEBUG] Total number of resource is: ", resourceCount)
	}
}
