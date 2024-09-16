package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLabels(t *testing.T) {
	assert := assert.New(t)

	nodegroupLabel := "azrebalance-disabled"

	ng0 := Nodegroup{labels: map[string]string{nodegroupLabel: "true"}}
	err := ng0.ParseLabels(nodegroupLabel)
	assert.Nil(err, nil, "Error should be nil")
	assert.True(ng0.LabelBool, "Parsed label should be true")

	ng1 := Nodegroup{labels: map[string]string{nodegroupLabel: "false"}}
	err = ng1.ParseLabels(nodegroupLabel)
	assert.Nil(err, nil, "Error should be nil")
	assert.False(ng1.LabelBool, "Parsed label should be false")

	ng2 := Nodegroup{labels: map[string]string{nodegroupLabel: "false"}}
	err = ng2.ParseLabels("another-label")
	assert.EqualError(err, ErrLabelNotPresent, "Label not found error")

	ng3 := Nodegroup{labels: map[string]string{nodegroupLabel: "fals"}}
	err = ng3.ParseLabels(nodegroupLabel)
	assert.Error(err, "Should not be able to parse label value to bool")
}
