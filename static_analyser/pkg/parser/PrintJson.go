package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	t "static_analyser/pkg/types"

)

func PrintJson(mani t.TCPManifest) {
	b, err := json.MarshalIndent(mani, "", " ")
	if err != nil {
		fmt.Println("error:", err)
	}
	//生成json文件
	err = ioutil.WriteFile(mani.Service+".json", b, 0777)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	var data interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println("data", data)
}
