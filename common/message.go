package common

import (
	"net/http"
)

// Code code
type Code string

func (c Code) String() string {
	if msg, ok := templates[c]; ok {
		return msg
	}
	return templates[ErrUnknown]
}

// all codes
const (
	// * plugin
	ErrPluginNotFound Code = "ErrPluginNotFound"
	ErrPluginInvalid       = "ErrPluginInvalid"

	// * request
	ErrRequestAccessDenied   = "ErrRequestAccessDenied"
	ErrRequestMethodNotFound = "ErrRequestMethodNotFound"
	ErrRequestParamInvalid   = "ErrRequestParamInvalid"
	// * resource
	ErrResourceNotFound        = "ErrResourceNotFound"
	ErrResourceAccessForbidden = "ErrResourceAccessForbidden"
	ErrResourceConflict        = "ErrResourceConflict"
	ErrResourceDeleteForbidden = "ErrResourceDeleteForbidden"
	ErrResourceHasBeenUsed     = "ErrResourceHasBeenUsed"
	ErrNodeNotReady            = "ErrNodeNotReady"
	ErrInvalidToken            = "ErrInvalidToken"

	// * volumes
	ErrVolumeType = "ErrVolumeType"
	// * unknown
	ErrUnknown = "UnknownError"
	// * application
	ErrAppNameConflict         = "ErrAppNameConflict"
	ErrVolumeNotFoundWhenMount = "ErrVolumeNotFoundWhenMount"
	ErrAppReferencedByNode     = "ErrAppReferencedByNode"
	// * node
	ErrNodeNumMaxLimit       = "ErrNodeNumMaxLimit"
	ErrNodeNumQueryException = "ErrNodeNumQueryException"

	// * config
	ErrConfigInUsed = "ErrConfigInUsed"
	// * register
	ErrRegisterQuotaNumOut     = "ErrRegisterQuotaNumOut"
	ErrRegisterDeleteRecord    = "ErrRegisterDeleteRecord"
	ErrRegisterDeleteCallback  = "ErrRegisterDeleteCallback"
	ErrRegisterPackage         = "ErrRegisterPackage"
	ErrRegisterRecordActivated = "ErrRegisterRecordActivated"
	// * db
	ErrDatabase  = "ErrDatabase"
	ErrUpdateCas = "ErrUpdateCas"
	// * k8s
	ErrK8S = "ErrK8S"
	// * ceph
	ErrCeph = "ErrCeph"
	// * io
	ErrIO = "ErrIO"
	// * template
	ErrTemplate = "ErrTemplate"
	// * function
	ErrFunction = "ErrFunction"
	// * resourceName
	ErrInvalidResourceName = "resourceName"
	ErrInvalidLabels       = "validLabels"
	ErrInvalidRequired     = "required"
	// * batchOp
	ErrInvalidArrayLength = "maxLength"
	// * fingerprintValue
	ErrInvalidFingerprintValue = "fingerprintValue"
	// * memory
	ErrInvalidMemory = "memory"
	// * duration
	ErrInvalidDuration = "duration"
	// * setcpus
	ErrInvalidSetcpus = "setcpus"
	// * nonBaetyl
	ErrInvalidName = "nonBaetyl"
	// * license
	ErrLicenseExpired      = "ErrLicenseExpired"
	ErrLicenseQuota        = "ErrLicenseQuota"
	ErrLicenseQuotaAcquire = "ErrLicenseQuotaAcquire"
	ErrLicenseQuotaRelease = "ErrLicenseQuotaRelease"
	// * third server error
	ErrThirdServer = "ErrThirdServer"
	// * object error
	ErrObjectOperationException = "ErrObjectOperationException"

	ErrResourceInvisible = "ErrResourceInvisible"
	ErrConvertConflict   = "ErrConvertConflict"

	ErrPubsubTimeout   = "ErrPubsubTimeout"
	ErrUpdateSubLabels = "ErrUpdateSubLabels"
)

var templates = map[Code]string{
	// * plugin
	ErrPluginNotFound: "The plugin{{if .name}} ({{.name}}){{end}} is not found.",
	ErrPluginInvalid:  "The plugin {{.name}} is invalid, not implement all interfaces of {{.kind}}.",
	// * request
	ErrRequestAccessDenied:   "访问已被拒绝。\nThe request access is denied.",
	ErrRequestMethodNotFound: "访问了不存在的api接口。\nThe request method is not found.",
	ErrRequestParamInvalid:   "非法的请求参数。\nThe request parameter is invalid.{{if .error}} ({{.error}}){{end}}",
	// * resource
	ErrResourceNotFound:        "访问不存在的资源。\nThe {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} is not found{{if .namespace}} in namespace({{.namespace}}){{end}}.",
	ErrResourceAccessForbidden: "The {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} can not be accessed{{if .namespace}} in namespace({{.namespace}}){{end}}.",
	ErrResourceConflict:        "The {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} already exist.",
	ErrResourceHasBeenUsed:     "该资源名称已被占用，请更换命名。The {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} has been used.",
	ErrResourceDeleteForbidden: "The {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} can not be deleted{{if .namespace}} in namespace({{.namespace}}){{end}}",
	// * volumes
	ErrVolumeType: "The volume{{if .name}} ({{.name}}){{end}} type should be{{if .type}} ({{.type}}){{end}}.",
	// * unknown
	ErrUnknown: "There is a unknown error{{if .error}} ({{.error}}){{end}}. If the attempt to retry does not work, please contact us.",
	// * application
	ErrAppNameConflict:         "A naming conflict occurs when you try to create/update app.{{if .where}} where={{.where}}.{{end}}{{if .name}} name={{.name}}.{{end}}",
	ErrVolumeNotFoundWhenMount: "The mount volume name{{if .name}}({{.name}}){{end}} can't find in the Volumes[].",
	ErrNodeNotReady:            "The node {{if .name}}({{.name}} ){{end}}is not ready, please retry later.",
	ErrAppReferencedByNode:     "The {{if .name}}({{.name}}){{end}} app is still referenced by a node.",
	// * node
	ErrNodeNumMaxLimit:       "节点个数已达上线，请联系相关人员申请更高节点限额。\nThe number of nodes reaches the maximum limit",
	ErrNodeNumQueryException: "The number of nodes is null",
	// * config
	ErrConfigInUsed: "该配置名称已被占用，请更换配置名称。\nThe config name {{if .name}}({{.name}}){{end}} in used.",
	// * register
	ErrRegisterQuotaNumOut:     "Number reached the upper limit {{if .num}}({{.num}}){{end}}",
	ErrRegisterDeleteRecord:    "Batch {{if .name}}({{.name}}){{end}} delete failed, record not null.",
	ErrRegisterDeleteCallback:  "Callback {{if .name}}({{.name}}){{end}} is used, cannot delete.",
	ErrRegisterPackage:         "Problem with package.{{if .error}} ({{.error}}){{end}}",
	ErrRegisterRecordActivated: "The record is activated.",
	// * db
	ErrDatabase: "Problem with database operation.{{if .error}} ({{.error}}){{end}}",
	// * k8s
	ErrK8S: "Problem with k8s operation.{{if .error}} ({{.error}}){{end}}",
	// * Ceph
	ErrCeph: "Problem with Ceph operation. {{if .error}} ({{.error}}){{end}}",
	// * IO
	ErrIO: "Problem with IO operation. {{if .error}} ({{.error}}){{end}}",
	// * Template
	ErrTemplate: "Problem with Template parse. {{if .error}} ({{.error}}){{end}}",
	// * function(cfc, aws lambda)
	ErrFunction: "Problem occurred when importing a function.{{if .error}} ({{.error}}){{end}}",

	ErrInvalidResourceName: "非法的资源名称，请根据以下提示重新命名。\nThe field ({{if .resourceName}}{{.resourceName}}{{end}}) beginning and ending with an alphanumeric character ([a-z0-9]) with dashes (-) or the string which is consist of no more than 63 characters.{{if .error}} {{.error}}{{end}}.",
	ErrInvalidLabels: "非法的Label名称，请根据以下提示重新命名。\nThe field ({{if .validLabels}}{{.validLabels}}{{end}}) must contains labels which can be an empty string or a string which is consist of no more than 63 alphanumeric characters, '-', '_', and must start and end with an alphanumeric character, " +
		"the key prefix is optional. If specified, the prefix must be a DNS subdomain: a series of DNS labels separated by dots (.), not longer than 253 characters in total, followed by a slash (/)",
	ErrInvalidRequired:         "{{if .error}}{{.error}}{{end}}",
	ErrInvalidFingerprintValue: "The field ({{if .fingerprintValue}}{{.fingerprintValue}}{{end}}) beginning and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-) or the string which is consist of no more than 63 characters",
	ErrInvalidMemory:           "The ({{if .memory}}{{.memory}}{{end}}) setting must be a positive integer, optionally followed by a corresponding unit (k|m|g|t|p)",
	ErrInvalidDuration:         "The ({{if .duration}}{{.duration}}{{end}}) must be a positive integer, optionally followed by a corresponding time unit (s|m|h)",
	ErrInvalidSetcpus: "The ({{if .setcpus}}{{.setcpus}}{{end}}) must be a comma-separated list or hyphen-separated range of CPUs a container can use, " +
		"a valid value might be 0-3 (to use the first, second, third, and fourth CPU) or 1,3 (to use the second and fourth CPU)",
	ErrInvalidName: "非法名称，baetyl命名已被系统占用，请重新命名。\nThe field ({{if .nonBaetyl}}{{.nonBaetyl}}{{end}}) cannot contain baetyl (case insensitive)",

	// * Token auth for init server
	ErrInvalidToken: "非法的token导致鉴权失败，请联系系统管理人员。\nThe token is invalid",

	// * License
	ErrLicenseExpired:      "The license {{if .name}}({{.name}}){{end}} has expired. {{if .error}} ({{.error}}){{end}}",
	ErrLicenseQuota:        "Check {{if .name}}({{.name}}){{end}} quota failed, the limited number is {{if .limit}}({{.limit}}){{end}}",
	ErrLicenseQuotaAcquire: "Check {{if .name}}({{.name}}){{end}} quota acquire failed, the acquire number is {{if .number}}({{.number}}){{end}}",
	ErrLicenseQuotaRelease: "Check {{if .name}}({{.name}}){{end}} quota release failed, the acquire number is {{if .number}}({{.number}}){{end}}",

	// * third server error
	ErrThirdServer: "Third server {{if .name}}({{.name}}){{end}} error.{{if .error}} ({{.error}}){{end}}",

	ErrObjectOperationException: "Problem with {{if .source}}({{.source}}){{end}} object operation.{{if .error}} ({{.error}}){{end}}",

	ErrInvalidArrayLength: "The length of the array exceeds the limit",

	ErrResourceInvisible: "The {{if .type}}({{.type}}) {{end}}resource{{if .name}} ({{.name}}){{end}} is not visible.",
	ErrConvertConflict:   "Problem with converting {{if .name}} ({{.name}}){{end}}.{{if .error}} ({{.error}}){{end}}",

	ErrPubsubTimeout:   "Publish or subscribe message timeout. {{if .error}} ({{.error}}){{end}}",
	ErrUpdateSubLabels: "Failed to update sub node labels. {{if .error}} ({{.error}}){{end}}",
}

func getHTTPStatus(c Code) int {
	switch c {
	case ErrResourceNotFound, ErrRequestMethodNotFound:
		return http.StatusNotFound
	case ErrRequestAccessDenied:
		return http.StatusUnauthorized
	case ErrResourceHasBeenUsed:
		return http.StatusForbidden
	case ErrUnknown:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
