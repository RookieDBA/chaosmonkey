// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/pkg/errors"
)

type instance struct{
	appName string
	accountName string
	regionName string
	stackName string
	clusterName string
	asgName string
	id string
	cloudProvider string
}

func (i instance) AppName() string {
	return i.appName
}

func (i instance) AccountName() string {
	return i.accountName
}

func (i instance) RegionName() string {
	return i.regionName
}

func (i instance) StackName() string {
	return i.stackName
}

func (i instance) ClusterName() string {
	return i.clusterName
}

func (i instance) ID() string {
	return i.id
}

func (i instance) CloudProvider() string {
	return i.cloudProvider
}

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	// Fail if not cluster-specific
	cluster, ok := group.Cluster()
	if !ok {
		return nil, errors.New("only supports cluster-specific grouping")
	}

	// Fail if not region-specific

	region, ok := group.Region()
	if !ok {
		return nil, errors.New("only supports region-specific grouping")
	}

	result := make([]chaosmonkey.Instance, 0)


	provider, err := dep.CloudProvider(group.Account())
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cloud provider failed")
	}

	ids, err := dep.GetInstanceIDs(group.App(), group.Account(), region, cluster)



	if err!=nil {
		return nil, err
	}

	for _, id := range ids {
		result = append(result,
			instance{appName: group.App(),
				accountName: group.Account(),
				regionName: region,
				stackName: ???,
				clusterName: ???,
				asgName: ???,
				id: id,
				cloudProvider, group.

			})
	}



}