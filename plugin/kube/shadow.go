package kube

import (
	"fmt"
	"strings"
	"time"

	"github.com/baetyl/baetyl-go/v2/json"
	"github.com/baetyl/baetyl-go/v2/log"
	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/baetyl/baetyl-go/v2/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/models"
	"github.com/baetyl/baetyl-cloud/v2/plugin/kube/apis/cloud/v1alpha1"
)

func (c *client) Get(tx interface{}, namespace, name string) (*models.Shadow, error) {
	defer utils.Trace(c.log.Debug, "shadow Get")()
	nodeDesire, err := c.customClient.CloudV1alpha1().NodeDesires(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nodeReport, err := c.customClient.CloudV1alpha1().NodeReports(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	shd := buildShadow(nodeDesire.Namespace, nodeDesire.Name, nodeDesire.CreationTimestamp.Time.UTC())
	shd.Desire = fromDesire(nodeDesire)
	shd.Report = fromReport(nodeReport)
	return shd, nil
}

func (c *client) Create(tx interface{}, shadow *models.Shadow) (*models.Shadow, error) {
	namespace := shadow.Namespace
	name := shadow.Name
	desire, err := toDesire(shadow)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(c.log.Debug, "shadow Create")()
	nd, err := c.customClient.CloudV1alpha1().NodeDesires(namespace).Create(c.ctx, desire, metav1.CreateOptions{})
	if err != nil {
		d, err := c.customClient.CloudV1alpha1().NodeDesires(namespace).Get(c.ctx, name, metav1.GetOptions{})
		if err == nil && d != nil {
			desire.ResourceVersion = d.ResourceVersion
			desire.Labels = d.Labels
			desire, err = c.customClient.CloudV1alpha1().NodeDesires(namespace).Update(c.ctx, desire, metav1.UpdateOptions{})
			if err != nil {
				return nil, err
			}
		}
	} else {
		desire = nd
	}

	report, err := toReport(shadow)
	if err != nil {
		return nil, err
	}

	nr, err := c.customClient.CloudV1alpha1().NodeReports(namespace).Create(c.ctx, report, metav1.CreateOptions{})
	if err != nil {
		r, err := c.customClient.CloudV1alpha1().NodeReports(namespace).Get(c.ctx, name, metav1.GetOptions{})
		if err == nil && r != nil {
			report.ResourceVersion = r.ResourceVersion
			report.Labels = r.Labels
			report, err = c.customClient.CloudV1alpha1().NodeReports(namespace).Update(c.ctx, report, metav1.UpdateOptions{})
			if err != nil {
				return nil, err
			}
		}
	} else {
		report = nr
	}

	shd := buildShadow(shadow.Namespace, shadow.Name, desire.CreationTimestamp.Time.UTC())
	shd.Report = fromReport(report)
	shd.Desire = fromDesire(desire)
	return shd, nil
}

func (c *client) List(namespace string, nodeList *models.NodeList) (*models.ShadowList, error) {

	option := metav1.ListOptions{
		LabelSelector: generatorLabelSelector(nodeList),
	}

	defer utils.Trace(c.log.Debug, "shadow List")()
	deisres, err := c.customClient.CloudV1alpha1().NodeDesires(namespace).List(c.ctx, option)
	if err != nil {
		return nil, err
	}
	reports, err := c.customClient.CloudV1alpha1().NodeReports(namespace).List(c.ctx, option)
	if err != nil {
		return nil, err
	}

	return toShadowListModel(nodeList, toDesireMap(deisres), toReportMap(reports)), nil
}

func (c *client) ListShadowByNames(tx interface{}, namespace string, names []string) ([]*models.Shadow, error) {
	if names == nil || len(names) < 1 {
		return nil, nil
	}
	var shadows []*models.Shadow
	for _, name := range names {
		shd, err := c.Get(tx, namespace, name)
		if err != nil {
			return nil, err
		}
		shadows = append(shadows, shd)
	}
	return shadows, nil
}

func (c *client) Delete(namespace, name string) error {
	defer utils.Trace(c.log.Debug, "shadow Delete")()
	err := c.customClient.CloudV1alpha1().NodeDesires(namespace).Delete(c.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		common.LogDirtyData(err,
			log.Any("type", common.NodeDesire),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}
	err = c.customClient.CloudV1alpha1().NodeReports(namespace).Delete(c.ctx, name, metav1.DeleteOptions{})
	if err != nil {
		common.LogDirtyData(err,
			log.Any("type", common.NodeReport),
			log.Any("namespace", namespace),
			log.Any("name", name),
			log.Any("operation", "delete"))
	}
	return err
}

func (c *client) UpdateDesire(tx interface{}, shadow *models.Shadow) error {
	desire, err := toDesire(shadow)
	if err != nil {
		return err
	}
	defer utils.Trace(c.log.Debug, "shadow UpdateDesire")()
	d, err := c.customClient.CloudV1alpha1().NodeDesires(shadow.Namespace).Get(c.ctx, desire.Name, metav1.GetOptions{})
	if err != nil {
		log.L().Error("get node desire error", log.Error(err))
		return err
	}
	desire.ResourceVersion = d.ResourceVersion
	desire.Labels = d.Labels
	desire, err = c.customClient.CloudV1alpha1().NodeDesires(shadow.Namespace).Update(c.ctx, desire, metav1.UpdateOptions{})
	if err != nil {
		log.L().Error("update node desire error", log.Error(err))
		return err
	}
	return nil
}

func (c *client) UpdateDesires(tx interface{}, shadows []*models.Shadow) error {
	if shadows == nil || len(shadows) < 1 {
		return nil
	}
	for _, shadow := range shadows {
		err := c.UpdateDesire(nil, shadow)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (c *client) UpdateReport(shadow *models.Shadow) (*models.Shadow, error) {
	report, err := toReport(shadow)
	if err != nil {
		return nil, err
	}
	defer utils.Trace(c.log.Debug, "shadow UpdateReport")()
	r, err := c.customClient.CloudV1alpha1().NodeReports(shadow.Namespace).Get(c.ctx, shadow.Name, metav1.GetOptions{})
	if err != nil {
		log.L().Error("get node report error", log.Error(err))
		return nil, err
	}
	report.ResourceVersion = r.ResourceVersion
	report.Labels = r.Labels
	report, err = c.customClient.CloudV1alpha1().NodeReports(shadow.Namespace).Update(c.ctx, report, metav1.UpdateOptions{})
	if err != nil {
		log.L().Error("update node report error", log.Error(err))
		return nil, err
	}
	shd := buildShadow(shadow.Namespace, shadow.Name, report.CreationTimestamp.Time.UTC())
	shd.Report = fromReport(report)
	return shd, nil
}

func toDesire(shadow *models.Shadow) (*v1alpha1.NodeDesire, error) {
	nodeDesire := &v1alpha1.NodeDesire{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeDesire",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      shadow.Name,
			Namespace: shadow.Namespace,
			Labels: map[string]string{
				common.LabelNodeName: shadow.Name,
			},
		},
	}

	if shadow.Desire == nil {
		shadow.Desire = models.BuildEmptyApps()
	}

	desire, err := json.Marshal(shadow.Desire)
	if err != nil {
		log.L().Error("node desire marshal exception", log.Error(err))
		return nil, err
	}

	nodeDesire.Spec.Desire = desire

	return nodeDesire, nil
}

func toReport(shadow *models.Shadow) (*v1alpha1.NodeReport, error) {
	report := &v1alpha1.NodeReport{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeReport",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      shadow.Name,
			Namespace: shadow.Namespace,
			Labels: map[string]string{
				common.LabelNodeName: shadow.Name,
			},
		},
	}

	if shadow.Report == nil {
		shadow.Report = models.BuildEmptyApps()
	}

	r, err := json.Marshal(shadow.Report)
	if err != nil {
		log.L().Error("node desire marshal exception", log.Error(err))
		return nil, err
	}

	report.Status.Report = r

	return report, nil
}

func fromDesire(desire *v1alpha1.NodeDesire) specV1.Desire {
	d := specV1.Desire{}
	if desire != nil && desire.Spec.Desire != nil {
		err := json.Unmarshal(desire.Spec.Desire, &d)
		if err != nil {
			log.L().Error("desire unmarshal exception", log.Error(err))
			d = models.BuildEmptyApps()
		}
	}

	return d
}

func fromReport(report *v1alpha1.NodeReport) specV1.Report {
	r := specV1.Report{}
	if report != nil && report.Status.Report != nil {
		err := json.Unmarshal(report.Status.Report, &r)
		if err != nil {
			log.L().Error("report unmarshal exception", log.Error(err))
			r = models.BuildEmptyApps()
		}
	}

	return r
}

func buildShadow(namespace, name string, createTime time.Time) *models.Shadow {
	shd := &models.Shadow{
		Name:              name,
		Namespace:         namespace,
		CreationTimestamp: createTime,
	}
	return shd
}

func toShadowModel(desire *v1alpha1.NodeDesire, report *v1alpha1.NodeReport) *models.Shadow {
	shadow := models.NewShadow(desire.Namespace, desire.Name)
	shadow.CreationTimestamp = desire.CreationTimestamp.UTC()

	if desire != nil && desire.Spec.Desire != nil {
		err := json.Unmarshal(desire.Spec.Desire, &shadow.Desire)
		if err != nil {
			log.L().Error("desire unmarshal exception", log.Error(err))
		}
	}
	if report != nil && report.Status.Report != nil {
		err := json.Unmarshal(report.Status.Report, &shadow.Report)
		if err != nil {
			log.L().Error("report unmarshal exception", log.Error(err))
		}
	}

	return shadow
}

func toShadowListModel(list *models.NodeList, desires map[string]v1alpha1.NodeDesire, reports map[string]v1alpha1.NodeReport) *models.ShadowList {
	result := &models.ShadowList{
		Items: make([]models.Shadow, 0, len(desires)),
	}
	for _, node := range list.Items {
		d := desires[node.Name]
		r := reports[node.Name]
		s := toShadowModel(&d, &r)
		result.Items = append(result.Items, *s)
	}
	result.Total = len(list.Items)
	return result
}

func toDesireMap(desires *v1alpha1.NodeDesireList) map[string]v1alpha1.NodeDesire {
	desireMap := make(map[string]v1alpha1.NodeDesire)
	for idx := range desires.Items {
		desireMap[desires.Items[idx].Name] = desires.Items[idx]
	}

	return desireMap
}

func toReportMap(reports *v1alpha1.NodeReportList) map[string]v1alpha1.NodeReport {
	reportMap := make(map[string]v1alpha1.NodeReport)
	for idx := range reports.Items {
		reportMap[reports.Items[idx].Name] = reports.Items[idx]
	}

	return reportMap
}

func generatorLabelSelector(nodeList *models.NodeList) string {
	names := make([]string, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		names = append(names, node.Name)
	}

	return fmt.Sprintf("%s in ( %s )", common.LabelNodeName, strings.Join(names, ","))
}
