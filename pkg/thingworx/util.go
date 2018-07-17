package thingworx

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
	"github.com/nexiles/thingworx-operator/pkg/apis/thingworx/v1alpha1"
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
