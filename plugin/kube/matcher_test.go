package kube

import (
	"gotest.tools/assert"
	"testing"
)

func TestLabelMatcher_Match(t *testing.T) {
	var matcher client
	labels := make(map[string]string)
	labels["a"] = "b"
	labels["c"] = "d"

	sl := "a in (b),c=d"
	res, _ := matcher.IsLabelMatch(sl, labels)
	assert.Equal(t, true, res)

	sl = "a=bc=d"
	_, err := matcher.IsLabelMatch(sl, labels)
	assert.Equal(t, true, err != nil)

	sl = ""
	res, _ = matcher.IsLabelMatch(sl, labels)
	assert.Equal(t, true, res)

	var sl1 string
	res, _ = matcher.IsLabelMatch(sl1, labels)
	assert.Equal(t, true, res)
}
