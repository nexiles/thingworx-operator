package thingworx

import (
	"github.com/sirupsen/logrus"
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
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
	}

	return
}

