package plugin

import (
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"

	"github.com/baetyl/baetyl-cloud/v2/models"
)

//go:generate mockgen -destination=../mock/plugin/node.go -package=plugin github.com/baetyl/baetyl-cloud/v2/plugin Node

type Node interface {
	GetNode(namespace, name string) (*v1.Node, error)
	CreateNode(namespace string, node *v1.Node) (*v1.Node, error)
	UpdateNode(namespace string, node *v1.Node) (*v1.Node, error)
	DeleteNode(namespace, name string) error
	ListNode(namespace string, listOptions *models.ListOptions) (*models.NodeList, error)
}
