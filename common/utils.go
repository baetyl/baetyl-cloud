package common

import (
	"strings"

	uuid2 "github.com/google/uuid"
)

// TODO: use uuid v4

// UUIDPrune generate uuid without '-'
func UUIDPrune() string {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(uuid.String(), "-", "")
}

// UUID generate uuid
func UUID() string {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}

func CompareNumericalString(a, b string) int {
	lengthA := len(a)
	lengthB := len(b)
	if lengthA > lengthB {
		return 1
	} else if lengthA < lengthB {
		return -1
	} else {
		return strings.Compare(a, b)
	}

}

func AddSystemLabel(labels map[string]string, infos map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}

	for name, value := range infos {
		labels[name] = value
	}

	return labels
}
