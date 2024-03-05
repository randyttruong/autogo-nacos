package name

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ReadName(url string) (string, error) {
	//使用Get方法获取服务器响应包数据
	resp, err := http.Get("http://" + url)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	//读取body内的内容
	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(content), nil
}
