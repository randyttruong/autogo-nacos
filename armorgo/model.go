package extractrequest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Grpcs struct {
	Filename string `json:"filename"`
	Package  string `json:"package"`
	Rpcs     []struct {
		Path    string `json:"path"`
		Service string `json:"service"`
		Name    string `json:"name"`
	} `json:"rpcs"`
}

// 避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}

func ModelFromJson(filePath string) *Grpcs{
	fmt.Printf("The file path is :%s\n", filePath)

	//读取json
	fileData, err := ioutil.ReadFile(filePath)
	dropErr(err)
	//fmt.Println(string(fileData))
	fmt.Println("\n[AutoArmor]: model successfully from JSON "+filePath)

	// json.Unmarshal([]byte(JSON_DATA),JSON对应的结构体)
	res := &Grpcs{}
	_ = json.Unmarshal([]byte(fileData), &res)
	//fmt.Println(res.Rpcs[1].Name)
	return res
}

