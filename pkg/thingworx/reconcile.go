package thingworx

import (
	"github.com/sirupsen/logrus"
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/core/v1"
)

func Reconcile(twx *v1alpha1.Thingworx) (err error) {
	logrus.Infof("Reconcile: cluster name: %s", twx.GetName())

	twx = twx.DeepCopy()

	changed := twx.SetDefaults()

	if changed {
		logrus.Infof("Reconcile: defaults changed. Updating resource: %v", twx)
		return sdk.Update(twx)
	}

	if twx.Status.Phase == v1alpha1.ClusterPhaseInitial {
		logrus.Infof("Reconcile: initial phase.")
		var secrets *v1.Secret
		// var cm *v1.ConfigMap

		// Create the Secrets for this cluster
		secrets, err = CreateSecrets(twx)
		if err != nil {
			return
		}

		// Create the ConfigMap resource for this cluster.
		_, err = CreateConfigMap(twx, secrets)
		if err != nil {
			return
		}

	}

	return
}

