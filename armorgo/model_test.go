package extractrequest

import (
	"testing"
)

func TestModelFromJson(t *testing.T) {
	filePath := "/usr/local/go/src/autoarmor/armorgo/resources/rpc_info.json"
	ModelFromJson(filePath)
}
