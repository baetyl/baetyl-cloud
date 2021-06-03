package common

import (
	"strings"

	"github.com/baetyl/baetyl-go/v2/spec/v1"
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

func UpdateSysAppByAccelerator(accelerator string, sysApps []string) []string {
	found := false
	index := 0
	for i, app := range sysApps {
		if strings.Contains(app, v1.BaetylGPUMetrics) {
			found = true
			index = i
			break
		}
	}
	if accelerator == v1.NVAccelerator {
		if !found {
			sysApps = append(sysApps, v1.BaetylGPUMetrics)
		}
	} else {
		if found {
			sysApps = append(sysApps[:index], sysApps[index+1:]...)
		}
	}
	return sysApps
}
