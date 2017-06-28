// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/pkg/errors"
	"github.com/SmartThingsOSS/frigga-go"
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

func (i instance) ASGName() string {
	return i.asgName
}

func (i instance) Name() string {
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
	region, ok := group.Region()
	if !ok {
		return nil, errors.New("only supports region-specific grouping")
	}

	// Fail if not cluster-specific
	cluster, ok := group.Cluster()
	if ok {
		// Fail if not region-specific
		return clusterRegionInstances(deploy.ClusterName(cluster), region, group, dep)
	}

	_, ok = group.Stack()
	if ok {
		return nil, errors.New("only supports app and cluster-specific grouping")
	}

	return appRegionInstances(group.App(), region, group, dep)
}

func clusterRegionInstances(cluster deploy.ClusterName, region string, group grp.InstanceGroup, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	result := make([]chaosmonkey.Instance, 0)


	cloudProvider, err := dep.CloudProvider(group.Account())
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cloud provider failed")
	}

	asgName, ids, err := dep.GetInstanceIDs(group.App(),deploy.AccountName(group.Account()), cloudProvider, deploy.RegionName(region), deploy.ClusterName(cluster))

	if err!=nil {
		return nil, err
	}

	for _, id := range ids {
		names, err := frigga.Parse(string(asgName))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse")
		}
		result = append(result,
			instance{appName: group.App(),
				accountName: group.Account(),
				regionName: region,
				stackName: names.Stack,
				clusterName: names.Cluster,
				asgName: string(asgName),
				id: string(id),
				cloudProvider: cloudProvider,
			})
	}

	return result, nil
}

func appRegionInstances(app string, region string, group grp.InstanceGroup, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	clusters, err := dep.GetClusterNames(app, deploy.AccountName(group.Account()))
	if err != nil {
		return nil, err
	}

	result := make([]chaosmonkey.Instance, 0)
	for _, cluster := range clusters {
		instances, err := clusterRegionInstances(cluster, region, group, dep)
		if err != nil {
			return nil, err
		}
		result = append(result, instances...)
	}

	return result, nil

}
