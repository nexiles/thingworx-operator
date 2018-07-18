package thingworx

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
	"k8s.io/api/core/v1"
	"bytes"
	"text/template"
	"fmt"
)

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(o metav1.Object, r metav1.OwnerReference) {
	o.SetOwnerReferences(append(o.GetOwnerReferences(), r))
}

// labelsForVault returns the labels for selecting the resources
// belonging to the given vault name.
func labelsForThingworx(name string) map[string]string {
	return map[string]string{
		"app": "thingworx",
		"thingworx_cluster": name}
}

// asOwner returns an owner reference set as the vault cluster CR
func asOwner(v *v1alpha1.Thingworx) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: api.SchemeGroupVersion.String(),
		Kind:       v1alpha1.ThingworxKind,
		Name:       v.Name,
		UID:        v.UID,
		Controller: &trueVar,
	}
}

func setDefault(cm *v1.ConfigMap, key string, dflt string) *v1.ConfigMap{

	if _, exists := cm.Data[key]; !exists {
		cm.Data[key] = dflt
	}

	return cm
}

func renderConfigMapTemplate(cm *v1.ConfigMap, key string) (err error) {

	if _, exists := cm.Data[key]; !exists {
		return fmt.Errorf("renderConfigMapTemplate: key missing: %s", key)
	}

	var t1 = template.Must(template.New(key).Parse(cm.Data[key]))

	buf := bytes.NewBufferString("")
	t1.Execute(buf, cm.Data)
	cm.Data[key] = buf.String()

	return
}