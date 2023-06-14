package actions

import (
	"reflect"
	"testing"
)

var (
	policyContent = make([]byte, 10)
)

func TestPolicy(t *testing.T) {

	testPolicy := &Policy{
		Majority:      1,
		SuperMajority: 2,
	}

	PutPolicy(*testPolicy, &policyContent)
	p, _ := ParsePolicy(policyContent, 2)

	if !reflect.DeepEqual(p, testPolicy) {
		t.Error("Parse not working for actions Policy")
	}
}
