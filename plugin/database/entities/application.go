// Package entities 数据库存储基本结构与方法
package entities

import (
	"reflect"
	"time"

	"github.com/baetyl/baetyl-go/v2/errors"
	"github.com/baetyl/baetyl-go/v2/json"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"k8s.io/api/core/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

type Application struct {
	ID              int64     `db:"id"`
	Namespace       string    `db:"namespace"`
	Name            string    `db:"name"`
	Version         string    `db:"version"`
	Type            string    `db:"type"`
	Mode            string    `db:"mode"`
	System          bool      `db:"is_system"`
	Labels          string    `db:"labels"`
	Selector        string    `db:"selector"`
	NodeSelector    string    `db:"node_selector"`
	Description     string    `db:"description"`
	InitService     string    `db:"init_services"`
	Services        string    `db:"services"`
	Volumes         string    `db:"volumes"`
	CronStatus      int       `db:"cron_status"`
	CronTime        time.Time `db:"cron_time"`
	CreateTime      time.Time `db:"create_time"`
	UpdateTime      time.Time `db:"update_time"`
	HostNetwork     bool      `db:"host_network"`
	DNSPolicy       string    `db:"dns_policy"`
	Replica         int       `db:"replica"`
	Workload        string    `db:"workload"`
	JobConfig       string    `db:"job_config"`
	Ota             string    `db:"ota"` // alter table baetyl_application [ add ota varchar(2048) default '{}' not null comment 'ota信息';]
	AutoScaleCfg    string    `db:"autoScaleCfg"`
	PreserveUpdates bool      `db:"preserve_updates"`
}

func ToAppModel(app *Application) (*specV1.Application, error) {
	labels := map[string]string{}
	err := json.Unmarshal([]byte(app.Labels), &labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	services := []specV1.Service{}
	err = json.Unmarshal([]byte(app.Services), &services)
	if err != nil {
		return nil, errors.Trace(err)
	}

	initServices := []specV1.Service{}
	err = json.Unmarshal([]byte(app.InitService), &initServices)
	if err != nil {
		return nil, errors.Trace(err)
	}

	volumes := []specV1.Volume{}
	err = json.Unmarshal([]byte(app.Volumes), &volumes)
	if err != nil {
		return nil, errors.Trace(err)
	}

	jobConfig := specV1.AppJobConfig{}
	err = json.Unmarshal([]byte(app.JobConfig), &jobConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}

	ota := specV1.OtaInfo{}
	err = json.Unmarshal([]byte(app.Ota), &ota)
	if err != nil {
		return nil, errors.Trace(err)
	}

	autoScaleCfg := specV1.AutoScaleCfg{}
	err = json.Unmarshal([]byte(app.AutoScaleCfg), &autoScaleCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &specV1.Application{
		Namespace:         app.Namespace,
		Name:              app.Name,
		Version:           app.Version,
		Type:              app.Type,
		Mode:              app.Mode,
		System:            app.System,
		CreationTimestamp: app.CreateTime.UTC(),
		Labels:            labels,
		Selector:          app.Selector,
		NodeSelector:      app.NodeSelector,
		Description:       app.Description,
		InitServices:      initServices,
		Services:          services,
		Volumes:           volumes,
		CronStatus:        specV1.CronStatusCode(app.CronStatus),
		CronTime:          app.CronTime.UTC(),
		UpdateTime:        app.UpdateTime.UTC(),
		HostNetwork:       app.HostNetwork,
		DNSPolicy:         v1.DNSPolicy(app.DNSPolicy),
		Replica:           app.Replica,
		Workload:          app.Workload,
		JobConfig:         &jobConfig,
		Ota:               ota,
		AutoScaleCfg:      &autoScaleCfg,
		PreserveUpdates:   app.PreserveUpdates,
	}, nil
}

func ToAppListModel(app *Application) *models.AppItem {
	labels := map[string]string{}
	err := json.Unmarshal([]byte(app.Labels), &labels)
	if err != nil {
		return nil
	}
	services := []specV1.Service{}
	err = json.Unmarshal([]byte(app.Services), &services)
	if err != nil {
		return nil
	}

	volumes := []specV1.Volume{}
	err = json.Unmarshal([]byte(app.Volumes), &volumes)
	if err != nil {
		return nil
	}

	jobConfig := specV1.AppJobConfig{}
	err = json.Unmarshal([]byte(app.JobConfig), &jobConfig)
	if err != nil {
		return nil
	}

	ota := specV1.OtaInfo{}
	err = json.Unmarshal([]byte(app.Ota), &ota)
	if err != nil {
		return nil
	}

	autoScaleCfg := specV1.AutoScaleCfg{}
	err = json.Unmarshal([]byte(app.AutoScaleCfg), &autoScaleCfg)
	if err != nil {
		return nil
	}

	return &models.AppItem{
		Name:              app.Name,
		Namespace:         app.Namespace,
		Type:              app.Type,
		Mode:              app.Mode,
		Version:           app.Version,
		Labels:            labels,
		Selector:          app.Selector,
		NodeSelector:      app.NodeSelector,
		CreationTimestamp: app.CreateTime.UTC(),
		Description:       app.Description,
		System:            app.System,
		CronStatus:        specV1.CronStatusCode(app.CronStatus),
		CronTime:          app.CronTime.UTC(),
		HostNetwork:       app.HostNetwork,
		Replica:           app.Replica,
		Workload:          app.Workload,
		JobConfig:         &jobConfig,
		Ota:               ota,
		AutoScaleCfg:      &autoScaleCfg,
		PreserveUpdates:   app.PreserveUpdates,
	}
}

func FromAppModel(namespace string, app *specV1.Application) (*Application, error) {
	labels, err := json.Marshal(app.Labels)
	if err != nil {
		return nil, errors.Trace(err)
	}

	services, err := json.Marshal(app.Services)
	if err != nil {
		return nil, errors.Trace(err)
	}

	initServices, err := json.Marshal(app.InitServices)
	if err != nil {
		return nil, errors.Trace(err)
	}

	volumes, err := json.Marshal(app.Volumes)
	if err != nil {
		return nil, errors.Trace(err)
	}

	jobConfig, err := json.Marshal(app.JobConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}

	if app.CronTime.IsZero() {
		app.CronTime = time.Now()
	}

	ota, err := json.Marshal(app.Ota)
	if err != nil {
		return nil, errors.Trace(err)
	}

	autoScaleCfg, err := json.Marshal(app.AutoScaleCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &Application{
		Name:            app.Name,
		Namespace:       namespace,
		Version:         GenResourceVersion(),
		Type:            app.Type,
		Mode:            app.Mode,
		System:          app.System,
		Labels:          string(labels),
		Selector:        app.Selector,
		NodeSelector:    app.NodeSelector,
		Description:     app.Description,
		InitService:     string(initServices),
		Services:        string(services),
		Volumes:         string(volumes),
		CronStatus:      int(app.CronStatus),
		CronTime:        app.CronTime.UTC(),
		HostNetwork:     app.HostNetwork,
		DNSPolicy:       string(app.DNSPolicy),
		Replica:         app.Replica,
		Workload:        app.Workload,
		JobConfig:       string(jobConfig),
		Ota:             string(ota),
		AutoScaleCfg:    string(autoScaleCfg),
		PreserveUpdates: app.PreserveUpdates,
	}, nil
}

func EqualApp(app1, app2 *specV1.Application) bool {
	if !equalVolume(app1.Volumes, app2.Volumes) {
		return false
	}

	if !equalServices(app1.Services, app2.Services) {
		return false
	}

	if !equalServices(app1.InitServices, app2.InitServices) {
		return false
	}

	if app1.CronStatus != app2.CronStatus {
		return false
	}

	if !app1.CronTime.Equal(app2.CronTime) {
		return false
	}

	if len(app1.Labels) != len(app2.Labels) || (len(app1.Labels) != 0 && !reflect.DeepEqual(app1.Labels, app2.Labels)) {
		return false
	}

	if app1.Replica != app2.Replica {
		return false
	}

	if app1.HostNetwork != app2.HostNetwork {
		return false
	}

	if app1.PreserveUpdates != app2.PreserveUpdates {
		return false
	}

	if app1.DNSPolicy != app2.DNSPolicy {
		return false
	}

	if app1.Workload != app2.Workload {
		return false
	}

	return app1.Selector == app2.Selector &&
		app1.Description == app2.Description &&
		app1.NodeSelector == app2.NodeSelector &&
		reflect.DeepEqual(app1.JobConfig, app2.JobConfig) &&
		reflect.DeepEqual(app1.Ota, app2.Ota) && reflect.DeepEqual(app1.AutoScaleCfg, app2.AutoScaleCfg)
}

func equalVolume(vol1, vol2 []specV1.Volume) bool {
	if len(vol1) != len(vol2) {
		return false
	}

	if len(vol1) != 0 {
		for i := range vol1 {
			v1 := vol1[i]
			v2 := vol2[i]
			flag := (v1.Name == v2.Name) &&
				((v1.Secret != nil && v2.Secret != nil && v1.Secret.Name == v2.Secret.Name && v1.Secret.Version == v2.Secret.Version) || (v1.Secret == nil && v2.Secret == nil)) &&
				((v1.Config != nil && v2.Config != nil && v1.Config.Name == v2.Config.Name && v1.Config.Version == v2.Config.Version) || (v1.Config == nil && v2.Config == nil)) &&
				((v1.HostPath != nil && v2.HostPath != nil && v1.HostPath.Path == v2.HostPath.Path) || (v1.HostPath == nil && v2.HostPath == nil)) &&
				((v1.EmptyDir != nil && v2.EmptyDir != nil && v1.EmptyDir.Medium == v2.EmptyDir.Medium && v1.EmptyDir.SizeLimit == v2.EmptyDir.SizeLimit) ||
					(v1.EmptyDir == nil && v2.EmptyDir == nil))
			if !flag {
				return false
			}
		}
	}
	return true
}

func equalServices(svc1, svc2 []specV1.Service) bool {
	if len(svc1) != len(svc2) {
		return false
	}

	if len(svc1) != 0 {
		for i := range svc1 {
			s1 := svc1[i]
			s2 := svc2[i]
			if flag := (s1.Name == s2.Name) && (s1.Hostname == s2.Hostname) && (s1.Image == s2.Image) && (s1.Replica == s2.Replica) && (s1.HostNetwork == s2.HostNetwork) &&
				(s1.Runtime == s2.Runtime) && reflect.DeepEqual(s1.SecurityContext, s2.SecurityContext) && (s1.LivenessProbe == s2.LivenessProbe) &&
				(s1.ReadinessProbe == s2.ReadinessProbe) && (s1.ImagePullPolicy == s2.ImagePullPolicy) &&
				reflect.DeepEqual(s1.JobConfig, s2.JobConfig) && reflect.DeepEqual(s1.FunctionConfig, s2.FunctionConfig); !flag {
				return false
			}

			if len(s1.VolumeMounts) != len(s2.VolumeMounts) || (len(s1.VolumeMounts) != 0 && !reflect.DeepEqual(s1.VolumeMounts, s2.VolumeMounts)) {
				return false
			}

			if len(s1.Ports) != len(s2.Ports) || (len(s1.Ports) != 0 && !reflect.DeepEqual(s1.Ports, s2.Ports)) {
				return false
			}

			if len(s1.Devices) != len(s2.Devices) || (len(s1.Devices) != 0 && !reflect.DeepEqual(s1.Devices, s2.Devices)) {
				return false
			}

			if s1.WorkingDir != s2.WorkingDir {
				return false
			}

			if len(s1.Args) != len(s2.Args) || (len(s1.Args) != 0 && !reflect.DeepEqual(s1.Args, s2.Args)) {
				return false
			}

			if len(s1.Command) != len(s2.Command) || (len(s1.Command) != 0 && !reflect.DeepEqual(s1.Command, s2.Command)) {
				return false
			}

			if len(s1.Env) != len(s2.Env) || (len(s1.Env) != 0 && !reflect.DeepEqual(s1.Env, s2.Env)) {
				return false
			}

			if len(s1.Labels) != len(s2.Labels) || (len(s1.Labels) != 0 && !reflect.DeepEqual(s1.Labels, s2.Labels)) {
				return false
			}

			if len(s1.Functions) != len(s2.Functions) || (len(s1.Functions) != 0 && !reflect.DeepEqual(s1.Functions, s2.Functions)) {
				return false
			}
			if s1.Resources != nil && s2.Resources != nil {
				if len(s1.Resources.Limits) != len(s2.Resources.Limits) || (len(s1.Resources.Limits) != 0 && !reflect.DeepEqual(s1.Resources.Limits, s2.Resources.Limits)) {
					return false
				}

				if len(s1.Resources.Requests) != len(s2.Resources.Requests) || (len(s1.Resources.Requests) != 0 && !reflect.DeepEqual(s1.Resources.Requests, s2.Resources.Requests)) {
					return false
				}
			} else if s1.Resources == nil && s2.Resources == nil {
				continue
			} else {
				return false
			}
		}
	}
	return true
}
