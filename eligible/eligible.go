// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/pkg/errors"
)

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) ([]*deploy.Instance, error) {
	// Fail if not cluster-specific
	if _, ok := group.Cluster(); !ok {
		return nil, errors.New("only supports cluster-specific grouping")
	}

	return nil, errors.New("not yet impelemented")
}