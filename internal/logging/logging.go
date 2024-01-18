package logging

import (
	"encoding/json"

	loggingv1 "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	corev1 "k8s.io/api/core/v1"
)

func buildSubscriptionChannel(resources Resources) string {
	adoc := resources.AddOnDeploymentConfig
	if adoc == nil || len(adoc.Spec.CustomizedVariables) == 0 {
		return defaultLoggingVersion
	}

	for _, keyvalue := range adoc.Spec.CustomizedVariables {
		if keyvalue.Name == subscriptionChannelValueKey {
			return keyvalue.Value
		}
	}
	return defaultLoggingVersion
}

func buildSecrets(resources Resources) ([]SecretValue, error) {
	secretsValue := []SecretValue{}
	for _, secret := range resources.Secrets {
		dataJSON, err := json.Marshal(secret.Data)
		if err != nil {
			return secretsValue, err
		}
		secretValue := SecretValue{
			Name: secret.Name,
			Data: string(dataJSON),
		}
		secretsValue = append(secretsValue, secretValue)
	}
	return secretsValue, nil
}

func buildClusterLogForwarderSpec(resources Resources) (*loggingv1.ClusterLogForwarderSpec, error) {
	clf := resources.ClusterLogForwarder
	for _, secret := range resources.Secrets {
		if err := templateWithSecret(&clf.Spec, secret); err != nil {
			return nil, err
		}
	}

	return &clf.Spec, nil
}

func templateWithSecret(spec *loggingv1.ClusterLogForwarderSpec, secret corev1.Secret) error {
	clfOutputName, ok := secret.Annotations[annotationTargetOutputName]
	if !ok {
		return nil
	}
	// TODO(JoaoBraveCoding) Validate that clfOutputName actually exists
	// TODO(JoaoBraveCoding) Validate secret

	for k, output := range spec.Outputs {
		if output.Name == clfOutputName {
			output.Secret = &loggingv1.OutputSecretSpec{
				Name: secret.Name,
			}
			spec.Outputs[k] = output
		}
	}

	return nil
}
