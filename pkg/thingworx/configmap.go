package thingworx

import (
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
	"github.com/sirupsen/logrus"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/apimachinery/pkg/api/errors"
)

func CreateConfigMap(twx *v1alpha1.Thingworx) (err error) {
	logrus.Infof("Creating ConfigMap for cluster: %s", twx.GetName())

	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: twx.Namespace,
			Name: config_map_name(twx),
			Labels: labelsForThingworx(twx.GetName()),
		},
	}

	addOwnerRefToObject(cm, asOwner(twx))

	err = sdk.Create(cm)
	if err != nil && !errors.IsAlreadyExists(err)  {
		logrus.Errorf("Could not create ConfigMap (%s): %v", cm.Name, err)
	}

	return nil
}

func config_map_name(twx *v1alpha1.Thingworx) string {
	return fmt.Sprintf("%s.configmap.thingworxes.%s", twx.GetName(), v1alpha1.GroupName)
}
