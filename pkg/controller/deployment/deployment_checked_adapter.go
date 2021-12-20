package deployment

import (
	"context"
	"fmt"

	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

var globalDeploymentChecker deploymentChecker

var knownCertManagerDeployments = map[string][]string{
	//https://github.com/jetstack/cert-manager/tree/master/deploy/charts/cert-manager
	//https://github.com/IBM/ibm-cert-manager-operator/blob/master/deploy/olm-catalog/ibm-cert-manager-operator/3.17.0/ibm-cert-manager-operator.v3.17.0.clusterserviceversion.yaml
	"openshift-operators": {"ibm-cert-mana`ger-operator", "cert-manager"},
	//https://github.com/komish/certmanagerdeployment-operator/blob/main/controllers/componentry/constants.go
	"cert-manager": {"cert-manager"},
}

func InitGlobalDeploymentChecker(deploymentGetter v1.DeploymentsGetter) {
	globalDeploymentChecker = newDeploymentChecker(deploymentGetter)
}

func StaticResourcesAdapterFunction(ctx context.Context, _ []string) (bool, error) {
	return getReadyStatus(ctx)
}

func GenericDeploymentControllerAdapterFunction(ctx context.Context) (bool, error) {
	return getReadyStatus(ctx)
}

func getReadyStatus(ctx context.Context) (bool, error) {
	deploymentsExist, err := globalDeploymentChecker.isAnyDeploymentPresent(ctx, knownCertManagerDeployments)
	if err != nil {
		return false, err
	}
	if deploymentsExist {
		return false, fmt.Errorf("backoff: one of the known Cert Manager Operators was found in the cluster: %v", knownCertManagerDeployments)
	}
	return true, nil
}
