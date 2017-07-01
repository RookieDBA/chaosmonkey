package command

import (
	"github.com/Netflix/chaosmonkey/spinnaker"
	"github.com/SmartThingsOSS/frigga-go"
	"fmt"
	"os"
	"github.com/Netflix/chaosmonkey/deploy"
)

// DumpRegions lists the regions that a cluster is in
func DumpRegions(cluster, account string, spin spinnaker.Spinnaker) {

	names, err := frigga.Parse(cluster)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}

	provider, err := spin.CloudProvider(account)
	if err != nil {
		fmt.Printf("ERROR: Could not retrieve provider for account: %s. Reason: %v\n", account, err)
		os.Exit(1)
	}


	regions, err := spin.GetRegionNames(names.App, deploy.AccountName(account), provider, deploy.ClusterName(cluster))
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	for _, region := range regions {
		fmt.Println(region)
	}

}