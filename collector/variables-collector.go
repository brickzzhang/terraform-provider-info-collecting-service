package collector

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func SourceConfigGet(source *map[string]*schema.Resource) (jsonStr []byte, err error) {
	count := 0
	sourceRegexp, _ := regexp.Compile(`^tencentcloud_([\w_]*)$`)
	var resp = make(map[string]interface{})
	var sourceMap = make(map[string][]string)
	for k, v := range *source {
		sourceName := sourceRegexp.FindStringSubmatch(k)
		sourceMap[sourceName[1]] = getKeys(v.Schema)
		count++
	}

	resp["TotalNum"] = count
	resp["SourceList"] = sourceMap

	jsonStr, err = json.Marshal(resp)
	if err != nil {
		fmt.Println("json marshal error!")
	}

	return
}
