package actions

import(
	"testing"
)

func TestPolicy(t *testing.T){
	p := ParsePolicy([PutPolicy(*policy)], int(content[]))
	if p == nil {
		t.Error("Could not")
		return
	}
}