package addon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetSecrets(t *testing.T) {
	var (
		defaultNamespace = "open-cluster-management"
		clusterSecret    = &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "cluster",
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		}
		defaultSecret = &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: defaultNamespace,
			},
			Data: map[string][]byte{
				"bar": []byte("baz"),
			},
		}
	)

	for _, tc := range []struct {
		name                    string
		secretName              string
		configResourceNamespace string
		mcAddonNamespace        string
		expecedError            bool
		expectedSecret          *corev1.Secret
	}{
		{
			name:                    "secret in cluster namespace",
			mcAddonNamespace:        "cluster",
			configResourceNamespace: defaultNamespace,
			secretName:              "foo",
			expectedSecret:          clusterSecret,
		},
		{
			name:                    "default secret used",
			mcAddonNamespace:        "cluster-no-secret",
			configResourceNamespace: defaultNamespace,
			secretName:              "foo",
			expectedSecret:          defaultSecret,
		},
		{
			name:                    "default namespace not defined",
			configResourceNamespace: "",
			expecedError:            true,
		},
		{
			name:                    "no secret found",
			mcAddonNamespace:        "cluster",
			configResourceNamespace: defaultNamespace,
			secretName:              "bar",
			expecedError:            true,
		},
	} {
		fakeKubeClient := fake.NewClientBuilder().
			WithObjects(clusterSecret, defaultSecret).
			Build()

		target := Endpoint("foo")
		targetSecrets := map[Endpoint]string{
			target: tc.secretName,
		}
		secrets, err := GetSecrets(context.TODO(), fakeKubeClient, tc.configResourceNamespace, tc.mcAddonNamespace, targetSecrets)
		if tc.expecedError {
			require.Error(t, err)
			return
		}
		secret, ok := secrets[target]
		require.True(t, ok)
		require.NoError(t, err)
		require.Equal(t, tc.expectedSecret.Name, secret.Name)
		require.Equal(t, tc.expectedSecret.Namespace, secret.Namespace)
		require.Equal(t, tc.expectedSecret.Data, secret.Data)
	}
}
