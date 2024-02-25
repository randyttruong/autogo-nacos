package extractrequest

import "testing"

func TestParseYaml(t *testing.T) {
	filePath := "D:\\APP\\Go\\src\\go\\armorgo\\resources/deployment_files/frontend.yaml"
	ParseYaml(filePath)
}
