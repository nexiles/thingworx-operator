package stub

import (
	"context"

	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/nexiles/thingworx-operator/pkg/thingworx"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) (err error) {
	logrus.Infof("Handle kind %s", event.Object.GetObjectKind())

	switch o := event.Object.(type) {
	case *v1alpha1.Thingworx:

		if event.Deleted {
			logrus.Debugf("Ignoring 'Deleted' event.")
			return nil
		}

		err = thingworx.Reconcile(o)
	}
	return
}


