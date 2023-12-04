// Package models 模型定义
package models

type NodeHostStats struct {
	Name        string
	Status      map[string]float64
	CPUStats    map[string]float64
	MemoryStats map[string]float64
	GPUStats    map[string]float64
	DiskStats   map[string]float64
	NetIOStats  map[string]float64
}

type NodeSysAppStats struct {
	Name        string
	Status      map[string]float64
	CPUStats    map[string]float64
	MemoryStats map[string]float64
}

type QPSStats struct {
	Name    string
	QPS     map[string]float64
	Success map[string]float64
	Fail    map[string]float64
}

type NodeUserAppStats struct {
	Name        string
	Status      map[string]float64
	CPUStats    map[string]float64
	MemoryStats map[string]float64
	QPSStats    map[string]QPSStats
}

type NodePromStats struct {
	Name         string
	Namespace    string
	SysAppStats  []NodeSysAppStats
	UserAppStats []NodeUserAppStats
	HostStats    []NodeHostStats
}
