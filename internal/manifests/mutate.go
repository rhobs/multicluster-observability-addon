package manifests

import (
	"reflect"

	"github.com/ViaQ/logerr/v2/kverrors"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MutateFuncFor returns a mutate function based on the existing resource's concrete type.
// It currently supports the following types and will return an error for other types:
//
//   - Secret
func MutateFuncFor(existing, desired client.Object, depAnnotations map[string]string) controllerutil.MutateFn {
	return func() error {
		existingAnnotations := existing.GetAnnotations()
		if err := mergeWithOverride(&existingAnnotations, desired.GetAnnotations()); err != nil {
			return err
		}
		existing.SetAnnotations(existingAnnotations)

		existingLabels := existing.GetLabels()
		if err := mergeWithOverride(&existingLabels, desired.GetLabels()); err != nil {
			return err
		}
		existing.SetLabels(existingLabels)

		if ownerRefs := desired.GetOwnerReferences(); len(ownerRefs) > 0 {
			existing.SetOwnerReferences(ownerRefs)
		}

		switch existing.(type) {
		case *corev1.Secret:
			s := existing.(*corev1.Secret)
			wantS := desired.(*corev1.Secret)
			mutateSecret(s, wantS)
			existingAnnotations := s.GetAnnotations()
			if err := mergeWithOverride(&existingAnnotations, depAnnotations); err != nil {
				return err
			}
			s.SetAnnotations(existingAnnotations)

		case *certmanagerv1.CertificateRequest:
			cr := existing.(*certmanagerv1.CertificateRequest)
			wantCr := desired.(*certmanagerv1.CertificateRequest)
			mutateCertificateRequest(cr, wantCr)

		default:
			t := reflect.TypeOf(existing).String()
			return kverrors.New("missing mutate implementation for resource type", "type", t)
		}
		return nil
	}
}

func mergeWithOverride(dst, src interface{}) error {
	err := mergo.Merge(dst, src, mergo.WithOverride)
	if err != nil {
		return kverrors.Wrap(err, "unable to mergeWithOverride", "dst", dst, "src", src)
	}
	return nil
}

func mutateSecret(existing, desired *corev1.Secret) {
	existing.Annotations = desired.Annotations
	existing.Labels = desired.Labels
	existing.Data = desired.Data
}

func mutateCertificateRequest(existing, desired *certmanagerv1.CertificateRequest) {
	existing.Annotations = desired.Annotations
	existing.Labels = desired.Labels
	existing.Spec = desired.Spec
}
