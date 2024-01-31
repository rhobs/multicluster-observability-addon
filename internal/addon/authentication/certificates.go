package authentication

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/v2/kverrors"
	"github.com/go-logr/logr"
	"github.com/rhobs/multicluster-observability-addon/internal/manifests"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdateRootCertificate(k8s client.Client, log logr.Logger) error {
	ctx := context.Background()

	err := checkCertManagerCRDs(ctx, k8s)
	if err != nil {
		return err
	}

	objects := manifests.BuildAllRootCertificate()

	var errCount int32
	for _, obj := range objects {
		l := log.WithValues(
			"object_name", obj.GetName(),
			"object_kind", obj.GetObjectKind(),
		)

		desired := obj.DeepCopyObject().(client.Object)
		mutateFn := manifests.MutateFuncFor(obj, desired, nil)

		op, err := ctrl.CreateOrUpdate(ctx, k8s, obj, mutateFn)
		if err != nil {
			l.Error(err, "failed to configure resource")
			errCount++
			continue
		}

		msg := fmt.Sprintf("resource has been %s", op)
		switch op {
		case ctrlutil.OperationResultNone:
			l.Info(msg)
		default:
			l.Info(msg)
		}
	}
	if errCount > 0 {
		return kverrors.New("failed to configure root certificate resources")
	}

	return nil
}

func checkCertManagerCRDs(ctx context.Context, k8s client.Client) error {
	for _, crdName := range certManagerCRDs {
		key := client.ObjectKey{Name: crdName}
		crd := &apiextensionsv1.CustomResourceDefinition{}
		err := k8s.Get(ctx, key, crd, &client.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return kverrors.New("cert-manager CRD is missing", "name", crdName)
			}
			return err
		}
	}

	return nil
}
