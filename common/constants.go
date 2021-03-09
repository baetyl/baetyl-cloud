package common

type Resource string
type VolumeType string
type State string
type Security string
type Event string
type SystemApplication string
type CertType string
type SyncMode string
type ModuleType string

const (
	// TimeFormat
	TimeFormat = "2006-01-02T15:04:05Z"
	// KeyContextNamespace the key of namespace in context
	KeyContextNamespace = "namespace"

	// ResourceName resource name
	ResourceName = "resourceName"
	// ResourceVersion resource version
	ResourceVersion = "resourceVersion"
	// ResourceInvisible resource invisible
	ResourceInvisible = "resource-invisible"
	// Application application resource
	Application Resource = "application"
	// App alias of application resource
	APP Resource = "app"
	// Configuration configuration resource
	Configuration Resource = "configuration"
	// Config config resource
	Config Resource = "config"
	// Secret secret resource
	Secret Resource = "secret"
	// Certificate certificate resource
	Certificate Resource = "certificate"
	// Registry registry resource
	Registry Resource = "registry"
	// Deprecated
	// Deployment deployment resource
	Deployment Resource = "deployment"
	// Node node resource
	Node Resource = "node"
	// Shadow shadow resource
	Shadow Resource = "shadow"
	// NodeDesire nodedesire resource
	NodeDesire Resource = "nodedesire"
	// NodeReport nodereport resource
	NodeReport Resource = "nodereport"
	// Batch batch resource
	Batch Resource = "batch"
	// Index index resource
	Index Resource = "index"
	// !deprecated
	// DefaultConfigDir default host dir of config
	DefaultConfigDir = "var/db/baetyl"
	// DefaultLibDir lib dir of the service by default
	DefaultLibDir = "var/lib/baetyl"
	// DefaultBinDir the file path of master binary
	DefaultBinDir = "bin"
	// DefaultBinFile the file of master binary
	DefaultBinFile = "baetyl"
	// DefaultActiveWebPort the port of web active
	DefaultActiveWebPort = 30007 // The range of valid ports is 30000-32767
	// AppConfFileName application config file name
	AppConfFileName = "application.yml"
	// DefaultSetupFile
	DefaultSetupFile = "setup.sh"
	// DefaultOffActivationFile
	DefaultOffActivationFile = "offlineActivation.sh"
	// DefaultOnActivationFile
	DefaultOnActivationFile = "onlineActivation.sh"
	// DefaultMasterConfDir master config dir by default
	DefaultMasterConfDir = "etc/baetyl"
	// DefaultMasterConfFile master config file by default
	DefaultMasterConfFile = "conf.yml"
	// DefaultMasterPlistFile
	DefaultMasterPlistFile = "baetyl.plist"
	// DefaultMasterDaemonFile
	DefaultMasterDaemonFile = "baetyl.service"
	// DefaultCoreFrequency
	DefaultCoreFrequency = "20"
	// DefaultCoreAPIPort
	DefaultCoreAPIPort = "30050"
	// ComposeVersion compose version
	ComposeVersion = "3"
	// Bind bind
	Bind VolumeType = "bind"
	// Volume volume
	Volume VolumeType = "volume"

	DesiredApplications    = "apps"
	DesiredSysApplications = "sysapps"

	// Receive receive
	Receive State = "RECEIVE"
	// Deploy deploy
	Deploy State = "DEPLOYED"
	// Updating updating
	Updating State = "UPDATING"
	// Updated updated
	Updated State = "UPDATED"
	// Timeout timeout
	Timeout State = "TIMEOUT"
	// Failure failure
	Failure State = "FAILURE"

	// Secret type
	SecretRegistry = "registry"
	SecretConfig   = "config"
	MaxRetryNum    = 20

	// LabelKeyFunction tag of function
	LabelKeyFunction = "baetyl-function"
	LabelNodeName    = "baetyl-node-name"
	LabelAppName     = "baetyl-app-name"
	LabelSystem      = "baetyl-cloud-system"
	LabelBatch       = "baetyl-batch"
)

const (
	BaetylCloud      = "baetyl-cloud"
	BaetylCloudGroup = "cloud.baetyl.io"
	Description      = "description"
	UpdateTimestamp  = "updateTimestamp"
	Metadata         = "matadata"
	PkiCertID        = "pkiCertID"
	NodeSelector     = "nodeSelector"

	AnnotationDescription     = BaetylCloudGroup + "/" + Description
	AnnotationUpdateTimestamp = BaetylCloudGroup + "/" + UpdateTimestamp
	AnnotationMetadata        = BaetylCloudGroup + "/" + Metadata
	AnnotationPkiCertID       = BaetylCloudGroup + "/" + PkiCertID
	AnnotationNodeSelector    = BaetylCloudGroup + "/" + NodeSelector
)

const (
	None                    Security = "None"
	Token                   Security = "Token"
	Cert                    Security = "Cert"
	Dongle                  Security = "Dongle"
	DefaultActiveTime                = 1167580800 //todo avoid mysql ERROR 1292. "2007-01-01 00:00:00" s
	DefaultSN                        = "agent-sn"
	DefaultSNPath                    = "/var/lib/baetyl/sn"
	DefaultSNFile                    = "fingerprint.txt"
	DefaultInputField                = "sn"
	DefaultAgentConfigMount          = "/etc/baetyl"
	FingerprintSN                    = 0x1
	FingerprintInput                 = 0x2
	FingerprintHostName              = 0x4
	FingerprintBootID                = 0x8
	FingerprintSystemUUID            = 0x10
	FingerprintMachineID             = 0x20

	FilenameYamlService       = "service.yml"
	PageSize                  = 20
	Online                    = "online"
	Offline                   = "offline"
	Success             Event = "success"
	Nothing             Event = "nothing"
	Failed              Event = "failed"
	Inactivated               = 0x0
	Activated                 = 0x1
	DisableWhitelist          = 0x0
	EnableWhitelist           = 0x1
)

const (
	Sync     CertType = "sync"
	Internal CertType = "internal"

	ObjectSource = "object-source"
)

const (
	AWSS3ReadPermission    = "public-read"
	AWSS3WritePermission   = "public-read-write"
	AWSS3PrivatePermission = "private"
)

const (
	ContainerApp = "container"
	FunctionApp  = "function"
)

const (
	ConfigObjectPrefix = "_object_"
)

const (
	UnpackTypeZip = "zip"
)

var (
	// Todo use global.go
	Cache          map[string]string
	FingerprintMap = map[int]string{
		0x1:  "sn",
		0x2:  "input",
		0x4:  "hostName",
		0x8:  "bootID",
		0x10: "systemUUID",
		0x20: "machineID",
	}
)

const (
	ReportMeta = "reportMeta"
	DesireMeta = "desireMeta"
	NodeProps  = "nodeprops"
	NodeInfo   = "node"
	NodeStats  = "nodestats"
)

const (
	TypeUser           ModuleType = "user"
	TypeUserRuntime    ModuleType = "runtime_user"
	TypeSystem         ModuleType = "system"
	TypeSystemOptional ModuleType = "opt_system"
)

const (
	TaskNamespaceDelete = "namespace-delete"
)
