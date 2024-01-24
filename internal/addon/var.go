package addon

import (
	"embed"
)

const (
	Name             = "multicluster-observability-addon"
	InstallNamespace = "open-cluster-management"

	McoaChartDir    = "manifests/charts/mcoa"
	MetricsChartDir = "manifests/charts/mcoa/charts/metrics"
	LoggingChartDir = "manifests/charts/mcoa/charts/logging"
	TracingChartDir = "manifests/charts/mcoa/charts/tracing"

	ConfigMapResource             = "configmaps"
	SecretResource                = "secrets"
	AddonDeploymentConfigResource = "addondeploymentconfigs"

	AdcMetricsDisabledKey = "metricsDisabled"
	AdcLoggingDisabledKey = "loggingDisabled"
	AdcTracingisabledKey  = "tracingDisabled"

	SignalLabelKey = "mcoa.openshift.io/signal"
	SignalMetrics  = "metrics"
	SignalLogging  = "logging"
	SignalTracing  = "tracing"
)

//go:embed manifests
//go:embed manifests/charts/mcoa
//go:embed manifests/charts/mcoa/templates/_helpers.tpl
//go:embed manifests/charts/mcoa/charts/logging/templates/_helpers.tpl
//go:embed manifests/charts/mcoa/charts/metrics/templates/_helpers.tpl
//go:embed manifests/charts/mcoa/charts/tracing/templates/_helpers.tpl
var FS embed.FS
