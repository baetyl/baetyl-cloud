package kube

import (
	"github.com/jinzhu/copier"
	kl "k8s.io/apimachinery/pkg/labels"
)

func (c *client) IsLabelMatch(labelSelector string, labels map[string]string) (bool, error) {
	selector, err := kl.Parse(labelSelector)

	if err != nil {
		return false, err
	}

	labelSet := kl.Set{}
	copier.Copy(&labelSet, &labels)

	return selector.Matches(labelSet), nil
}
