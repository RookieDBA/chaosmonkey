package command

import (
	"fmt"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/Netflix/chaosmonkey/spinnaker"
	"github.com/SmartThingsOSS/frigga-go"
	"os"
)

// DumpRegions lists the regions that a cluster is in
func DumpRegions(cluster, account string, spin spinnaker.Spinnaker) {

	names, err := frigga.Parse(cluster)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}

	regions, err := spin.GetRegionNames(names.App, deploy.AccountName(account), deploy.ClusterName(cluster))
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	for _, region := range regions {
		fmt.Println(region)
	}

}
