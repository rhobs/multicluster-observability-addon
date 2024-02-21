# Demo: Multicluster Observability AddOn

The following steps demonstrate how to use the `multicluster-observability-addon` to collect and forward observability signals (metrics, logs and traces) across a Red Hat Advanced Cluster Management (RHACM) managed fleet of OpenShift (OCP) clusters. Currently the `multicluster-observability-addon` is limited to only collection and forwarding of signals. Thus step 1 & 2 are dedicated to configure both RHACM and `LokiStack` as the storage solution on the RHACM hub cluster.

All steps are meant to be run on the hub cluster except when explicitely said.

## 1. Prerequisites: RHACM and OCP cluster fleet

For this demo you will need at least two OCP clusters (hosted on AWS) with one of them (the hub) having at least machines of flavor `m6a.4xlarge` in order to have enough resources for `LokiStack`. You will also need an S3 Bucket in the same region as the hub cluster.
 
1. Use the OpenShift Installer to create and setup two OCP cluster on AWS.
1. Install the `Advanced Cluster Management for Kubernetes` operator.
1. Create a `MultiClusterHub` resource.
1. Import each spoke cluster `RHACM` via the web console, using the commands option by running the commands on each spoke cluster.

## 2. Configure stores on the hub cluster

The following steps use Helm to install a set of RHACM `ConfigurationPolicies` that facilitate the configuration of the different stores that will be used in this demo.

_Hint:_ `cert-manager` operator is also installed in this step on the hub cluster. It will be used to demonstrate the ability to delegate PKI management for all tenants to a third-party tool.

1. Set the values in `demo/mcoa-demo/values.yaml`
1. Deploy the chart `helm upgrade --install mcoa-demo demo/mcoa-demo/`. This Helm chart will bootstrap configuration on the hub cluster to enabled it to receive signals from the spoke clustes.
1. Run `oc label --overwrite managedcluster/local-cluster cluster.open-cluster-management.io/clusterset=hub-mcoa-clusters` to label the `local-cluster` a.k.a. hub so that the policy applies to it.

### 2.1 Deploying a workload on the spoke cluster

Since we will want to observe all the 3 observability signals (metrics, logs and traces) in the end we will deploy a workload on the spoke clusters to generate these signals

1. On the spoke cluster run `helm upgrade --install demo-workload demo/demo-workload/` 

## 3. Install multicluster-observability-addon

1. Deploy the addon controller by running `make addon-deploy`.

## 4. Manage observability signals across spoke clusters

### 4.1 Signal configuration

Before enabling the AddOn on spoke clustes we first have define the configuration of each signal:

1. Set the values in `demo/addon-config/values.yaml`
1. Deploy the signal configuration with `helm upgrade --install addon-config demo/addon-config/`

### 4.2 Enable the addon for spoke clusters

The following chart will deploy the `ManagedClusterAddOn` resource that enables the AddOn on each spoke cluster.

1. Set the values in `demo/addon-install/values.yaml`.
1. Deploy it with `helm upgrade --install addon-install demo/addon-install/`. 

## 5. Validate with Grafana

Once the addon is installed on the spoke clusters and the respective forwarders have started to send their signal data you can
access this grafana instance to see the 3 signals in a single dashboard (roughly 2 minutes after installing the addon).

`oc -n grafana-operator get route grafana-route`

**Note: user: `root` password: `secret`**

There will be a pre-loaded dashboard called "MCOA" that will contain a view of the 3 signals regarding pods in the `openshift-console-operator` namespace.