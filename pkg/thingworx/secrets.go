package thingworx

import (
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

func CreateSecrets(twx *v1alpha1.Thingworx) (secret *v1.Secret, err error) {
	logrus.Infof("CreateSecrets: creating secrets for: `%s`", twx.GetName())

	secret = &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: twx.Namespace,
			Name: secretsName(twx),
			Labels: labelsForThingworx(twx.GetName()),
		},
		Data: map[string][]byte{
			TWX_ADMIN_PASSWORD_KEY: RandomPassword(PasswordLength),
			TWX_ADMIN_USER_KEY: []byte(TWX_ADMIN_USER),
		},
	}

	addOwnerRefToObject(secret, asOwner(twx))

	err = sdk.Create(secret)
	if err != nil && !errors.IsAlreadyExists(err)  {
		return nil, fmt.Errorf("could not create Secret (%s): %v", secret.Name, err)
	}

	logrus.Infof("Created Secret: %s", secret.Name)
	return secret, nil
}
