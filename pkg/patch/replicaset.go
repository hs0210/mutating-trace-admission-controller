package patch

import (
	"encoding/json"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func mutateReplicaSet(raw []byte, patchAnnotations map[string]string) *v1beta1.AdmissionResponse {
	var replicaSet appv1.ReplicaSet
	err := json.Unmarshal(raw, &replicaSet)
	if err != nil {
		glog.Errorf("Could not unmarshal raw object: %v", err)
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	// FIXME: use temporary measures to avoid bugs(the infinite loop of replicaset when update deployment)
	if replicaSet.OwnerReferences != nil {
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}

	patchBytes, err := createPatch(replicaSet.Annotations, patchAnnotations)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	glog.V(3).Infof("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}
