package helm

import (
	"context"
	"testing"

	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"

	"github.com/rhobs/multicluster-observability-addon/internal/addon"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	"open-cluster-management.io/addon-framework/pkg/addonmanager/addontesting"
	"open-cluster-management.io/addon-framework/pkg/agent"
	addonapiv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	_ = operatorsv1.AddToScheme(scheme.Scheme)
	_ = operatorsv1alpha1.AddToScheme(scheme.Scheme)
	_ = addonapiv1alpha1.AddToScheme(scheme.Scheme)
	_ = apiextensionsv1.AddToScheme(scheme.Scheme)
)

func Test_Mcoa_Disable_Charts(t *testing.T) {
	var (
		managedCluster        *clusterv1.ManagedCluster
		managedClusterAddOn   *addonapiv1alpha1.ManagedClusterAddOn
		addOnDeploymentConfig *addonapiv1alpha1.AddOnDeploymentConfig
	)

	managedCluster = addontesting.NewManagedCluster("cluster-1")
	managedClusterAddOn = addontesting.NewAddon("test", "cluster-1")

	managedClusterAddOn.Status.ConfigReferences = []addonapiv1alpha1.ConfigReference{
		{
			ConfigGroupResource: addonapiv1alpha1.ConfigGroupResource{
				Group:    "addon.open-cluster-management.io",
				Resource: "addondeploymentconfigs",
			},
			ConfigReferent: addonapiv1alpha1.ConfigReferent{
				Namespace: "open-cluster-management-observability",
				Name:      "multicluster-observability-addon",
			},
		},
	}

	addOnDeploymentConfig = &addonapiv1alpha1.AddOnDeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multicluster-observability-addon",
			Namespace: "open-cluster-management-observability",
		},
		Spec: addonapiv1alpha1.AddOnDeploymentConfigSpec{
			CustomizedVariables: []addonapiv1alpha1.CustomizedVariable{},
		},
	}

	fakeKubeClient := fake.NewClientBuilder().
		WithScheme(scheme.Scheme).
		WithObjects(addOnDeploymentConfig).
		Build()

	agentAddon, err := addonfactory.NewAgentAddonFactory(addon.Name, addon.FS, addon.McoaChartDir).
		WithGetValuesFuncs(GetValuesFunc(context.TODO(), fakeKubeClient)).
		WithAgentRegistrationOption(&agent.RegistrationOption{}).
		WithScheme(scheme.Scheme).
		BuildHelmAgentAddon()
	if err != nil {
		klog.Fatalf("failed to build agent %v", err)
	}

	objects, err := agentAddon.Manifests(managedCluster, managedClusterAddOn)
	require.NoError(t, err)
	require.Empty(t, objects)
}

func Test_Mcoa_Disable_Chart_Hub(t *testing.T) {
	var (
		managedCluster        *clusterv1.ManagedCluster
		managedClusterAddOn   *addonapiv1alpha1.ManagedClusterAddOn
		addOnDeploymentConfig *addonapiv1alpha1.AddOnDeploymentConfig
	)

	managedCluster = addontesting.NewManagedCluster("cluster-1")
	managedCluster.Labels = map[string]string{
		"local-cluster": "true",
	}
	managedClusterAddOn = addontesting.NewAddon("test", "cluster-1")

	managedClusterAddOn.Status.ConfigReferences = []addonapiv1alpha1.ConfigReference{
		{
			ConfigGroupResource: addonapiv1alpha1.ConfigGroupResource{
				Group:    "addon.open-cluster-management.io",
				Resource: "addondeploymentconfigs",
			},
			ConfigReferent: addonapiv1alpha1.ConfigReferent{
				Namespace: "open-cluster-management",
				Name:      "multicluster-observability-addon",
			},
		},
	}

	addOnDeploymentConfig = &addonapiv1alpha1.AddOnDeploymentConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "multicluster-observability-addon",
			Namespace: "open-cluster-management-observability",
		},
		Spec: addonapiv1alpha1.AddOnDeploymentConfigSpec{
			CustomizedVariables: []addonapiv1alpha1.CustomizedVariable{
				{
					Name:  addon.KeyPlatformLogsCollection,
					Value: string(addon.ClusterLogForwarderV1),
				},
			},
		},
	}

	fakeKubeClient := fake.NewClientBuilder().
		WithScheme(scheme.Scheme).
		WithObjects(addOnDeploymentConfig).
		Build()

	loggingAgentAddon, err := addonfactory.NewAgentAddonFactory(addon.Name, addon.FS, addon.McoaChartDir).
		WithGetValuesFuncs(GetValuesFunc(context.TODO(), fakeKubeClient)).
		WithAgentRegistrationOption(&agent.RegistrationOption{}).
		WithScheme(scheme.Scheme).
		BuildHelmAgentAddon()
	if err != nil {
		klog.Fatalf("failed to build agent %v", err)
	}

	objects, err := loggingAgentAddon.Manifests(managedCluster, managedClusterAddOn)
	require.NoError(t, err)
	require.Empty(t, objects)
}
