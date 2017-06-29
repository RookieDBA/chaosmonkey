// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/pkg/errors"
	"github.com/SmartThingsOSS/frigga-go"
	"log"
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


func isException(exs []chaosmonkey.Exception, account deploy.AccountName, cluster deploy.ClusterName, region deploy.RegionName) bool {
	names, err := frigga.Parse(string(cluster))
	if err != nil {
		log.Printf("ERROR Couldn't parse cluster name %s", cluster)
		return false
	}

	for _, ex := range exs {
		if ex.Matches(string(account), names.Stack, names.Detail, string(region)) {
			return true
		}
	}
	return false
}

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	r, ok := group.Region()
	if !ok {
		return nil, errors.New("only supports region-specific grouping")
	}

	region := deploy.RegionName(r)

	switch cfg.Grouping {
	case chaosmonkey.App:
		return appRegionInstances(group.App(), cfg.Exceptions, region, group, dep)
	case chaosmonkey.Stack:
		return nil, errors.New("stack-level grouping not yet implemented")
	case chaosmonkey.Cluster:
		cluster, ok := group.Cluster()
		if !ok {
			return nil, errors.Errorf("app %s is configured cluster-only but not cluster specified in group %s", group.App(), group)
		}

		if isException(cfg.Exceptions, deploy.AccountName(group.Account()), deploy.ClusterName(cluster), region) {
			return nil, nil
		}

		return clusterRegionInstances(deploy.ClusterName(cluster), region, group, dep)
	default:
		return nil, errors.New("only app/stack/cluster groupings supported")
	}

}

func clusterRegionInstances(cluster deploy.ClusterName, region deploy.RegionName, group grp.InstanceGroup, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {


	result := make([]chaosmonkey.Instance, 0)

	cloudProvider, err := dep.CloudProvider(group.Account())
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cloud provider failed")
	}

	asgName, ids, err := dep.GetInstanceIDs(group.App(),deploy.AccountName(group.Account()), cloudProvider, region, deploy.ClusterName(cluster))

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
				regionName: string(region),
				stackName: names.Stack,
				clusterName: names.Cluster,
				asgName: string(asgName),
				id: string(id),
				cloudProvider: cloudProvider,
			})
	}

	return result, nil
}

func appRegionInstances(app string, exs []chaosmonkey.Exception, region deploy.RegionName, group grp.InstanceGroup, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	account := deploy.AccountName(group.Account())
	clusters, err := dep.GetClusterNames(app, account)
	if err != nil {
		return nil, err
	}

	result := make([]chaosmonkey.Instance, 0)
	for _, cluster := range clusters {

		if isException(exs, account, cluster, region) {
			continue
		}

		instances, err := clusterRegionInstances(cluster, region, group, dep)
		if err != nil {
			return nil, err
		}
		result = append(result, instances...)
	}

	return result, nil

}
