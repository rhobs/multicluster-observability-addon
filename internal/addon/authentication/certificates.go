package authentication

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/v2/kverrors"
	"github.com/rhobs/multicluster-observability-addon/internal/manifests"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdateRootCertificate(k8s client.Client) error {
	ctx := context.Background()

	err := checkCertManagerCRDs(ctx, k8s)
	if err != nil {
		return nil
	}

	objects := manifests.BuildAllRootCertificate()

	for _, obj := range objects {
		desired := obj.DeepCopyObject().(client.Object)
		mutateFn := manifests.MutateFuncFor(obj, desired, nil)

		op, err := ctrl.CreateOrUpdate(ctx, k8s, obj, mutateFn)
		if err != nil {
			klog.Error(err, "failed to configure resource")
			continue
		}

		msg := fmt.Sprintf("Resource has been %s", op)
		switch op {
		case ctrlutil.OperationResultNone:
			klog.Info(msg)
		default:
			klog.Info(msg)
		}
	}

	return nil
}

func checkCertManagerCRDs(ctx context.Context, k8s client.Client) error {
	crds := []string{"certificates.cert-manager.io", "issuers.cert-manager.io", "clusterissuers.cert-manager.io"}

	for _, crdName := range crds {
		key := client.ObjectKey{Name: crdName}
		crd := &apiextensions.CustomResourceDefinition{}
		err := k8s.Get(ctx, key, crd, &client.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return kverrors.New("cert-manager CRD is missing %s", crdName)
			}
			return err
		}
	}

	return nil
}
