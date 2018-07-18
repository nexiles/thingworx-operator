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

const (
	platformSettingsKey = "platform_settings"
	startupScriptKey    = "startup_script"
)

const platformSettingsTemplate = `{
  "PlatformSettingsConfig": {
	"BasicSettings": {
	  "EnableSystemLogging": true,
	  "BackupStorage": "/data/thingworx/backup",
	  "FileRepositoryRoot": "/data/thingworx/storage",
	  "Storage": "/data/thingworx/storage"
	},
	"AdministratorUserSettings": {
	  "InitialPassword": "${THINGWORX_ADMIN_PASSWORD}"
	}
  },
  "LicensingConnectionSettings": {
	"username": "${PTC_USER_NAME}",
	"password": "${PTC_USER_PASS}"
  },
  "PersistenceProviderPackageConfigs": {
	"PostgresPersistenceProviderPackage": {
	  "ConnectionInformation": {
		"jdbcUrl": "jdbc:postgresql://acid-thingworx-cluster:5432/thingworx",
		"password": "${THINGWORX_DATABASE_PASSWORD}",
		"username": "${THINGWORX_DATABASE_USER}"
	  }
	}
  }
}`

const startupScriptTemplate =`
mkdir -p $THINGWORX_PLATFORM_SETTINGS
touch $THINGWORX_PLATFORM_SETTINGS/.startup-script
cat >$THINGWORX_PLATFORM_SETTINGS/platform-settings.json <<EOL
{{ .platform_settings }}
EOL

export CATALINA_OPTS="-Djava.library.path=/usr/local/tomcat/webapps/Thingworx/WEB-INF/extensions"

cd $CATALINA_HOME
bin/catalina.sh run
`

func CreateConfigMap(twx *v1alpha1.Thingworx, secrets *v1.Secret) (cm *v1.ConfigMap, err error) {
	logrus.Infof("Creating ConfigMap for cluster: %s", twx.GetName())

	cm = &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: twx.Namespace,
		},
	}

	if len(twx.Spec.ConfigMapName) != 0 {
		cm.Name = twx.Spec.ConfigMapName
		err := sdk.Get(cm)
		if err != nil {
			return nil, fmt.Errorf("prepare config error: get configmap (%s) failed: %v", twx.Spec.ConfigMapName, err)
		}
	}

	cm.Name = configMapName(twx)
	cm.Labels = labelsForThingworx(twx.GetName())

	setDefaultConfig(twx, cm)
	updatePlatformSettings(twx, cm, secrets)

	addOwnerRefToObject(cm, asOwner(twx))

	err = sdk.Create(cm)
	if err != nil && !errors.IsAlreadyExists(err)  {
		return nil, fmt.Errorf("could not create ConfigMap (%s) %v", cm.Name, err)
	}

	logrus.Infof("Created ConfigMap: %s", cm.Name)
	return cm, nil
}

func setDefaultConfig(twx *v1alpha1.Thingworx, cm *v1.ConfigMap) {
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	setDefault(cm, "VERSION", "8.3.0")
	setDefault(cm, "DATABASE_URL", "jdbc:postgresql://acid-thingworx-cluster:5432/thingworx")
	setDefault(cm, "DATABASE_PASSWORD", "thingworx")
	setDefault(cm, "DATABASE_USER", "thingworx")
	setDefault(cm, "STORAGE_BASE", "/data/thingworx")
	setDefault(cm, "PLATFORM_SETTINGS", "/data/thingworx/platform")
	setDefault(cm, platformSettingsKey, platformSettingsTemplate)
	setDefault(cm, startupScriptKey, startupScriptTemplate)
}

func updatePlatformSettings(twx *v1alpha1.Thingworx, cm *v1.ConfigMap, secrets *v1.Secret) *v1.ConfigMap {
	var data, _ = copyMap(cm.Data)

	if renderConfigMapTemplate(cm, platformSettingsKey, data) != nil {
		return cm
	}

	if renderConfigMapTemplate(cm, startupScriptKey, data) != nil {
		return cm
	}

	return cm
}

