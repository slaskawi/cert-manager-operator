package deployment

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type deploymentsGetterFunc func(ctx context.Context, namespace, deploymentName string) error

type deploymentChecker struct {
	deploymentsGetter     v1.DeploymentsGetter
	deploymentsGetterFunc deploymentsGetterFunc
}

func newDeploymentChecker(deploymentsGetter v1.DeploymentsGetter) deploymentChecker {
	ret := deploymentChecker{
		deploymentsGetter: deploymentsGetter,
	}
	ret.deploymentsGetterFunc = ret.deploymentsGetterImpl
	return ret
}

func (d *deploymentChecker) isAnyDeploymentPresent(ctx context.Context, namespacesAndDeployments map[string][]string) (bool, error) {
	for namespace, deployments := range namespacesAndDeployments {
		for _, deployment := range deployments {
			ret, err := d.isDeploymentPresent(ctx, namespace, deployment)
			if err != nil {
				return false, err
			}
			if ret == true {
				return true, nil
			}
		}
	}
	return false, nil
}

func (d *deploymentChecker) isDeploymentPresent(ctx context.Context, namespace, deploymentName string) (bool, error) {
	err := d.deploymentsGetterFunc(ctx, namespace, deploymentName)
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func (d *deploymentChecker) deploymentsGetterImpl(ctx context.Context, namespace, deploymentName string) error {
	_, err := d.deploymentsGetter.Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	return err
}
