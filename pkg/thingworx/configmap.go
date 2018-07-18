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

func CreateConfigMap(twx *v1alpha1.Thingworx) (err error) {
	logrus.Infof("Creating ConfigMap for cluster: %s", twx.GetName())

	cm := &v1.ConfigMap{
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
			return fmt.Errorf("prepare config error: get configmap (%s) failed: %v", twx.Spec.ConfigMapName, err)
		}
	}

	cm.Name = configMapName(twx)
	cm.Labels = labelsForThingworx(twx.GetName())

	setDefaultConfig(twx, cm)
	updatePlatformSettings(twx, cm)

	addOwnerRefToObject(cm, asOwner(twx))

	err = sdk.Create(cm)
	if err != nil && !errors.IsAlreadyExists(err)  {
		logrus.Errorf("Could not create ConfigMap (%s): %v", cm.Name, err)
	}

	logrus.Infof("Created ConfigMap: %s", cm.Name)
	return nil
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

func updatePlatformSettings(twx *v1alpha1.Thingworx, cm *v1.ConfigMap) *v1.ConfigMap {
	// TODO: Actually replace things here

	if renderConfigMapTemplate(cm, platformSettingsKey) != nil {
		return cm
	}
	if renderConfigMapTemplate(cm, startupScriptKey) != nil {
		return cm
	}

	return cm
}

func configMapName(twx *v1alpha1.Thingworx) string {
	return fmt.Sprintf("%s.configmap.thingworxes.%s", twx.GetName(), v1alpha1.GroupName)
}
