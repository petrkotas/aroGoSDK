package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/petrkotas/aroGoSDK/pkg/client/services/redhatopenshift/mgmt/2020-04-30/redhatopenshift"
	"github.com/petrkotas/aroGoSDK/pkg/util/azureclient"
	mgmtredhatopenshift "github.com/petrkotas/aroGoSDK/pkg/util/azureclient/mgmt/redhatopenshift/2020-04-30/redhatopenshift"
)

func main() {
	log.Println("starting the cluster creation")

	// ocm service principal used to talk to azure
	ocmServicePrincipalClientID := ""
	ocmServicePrincipalClientSecret := ""
	ocmServicePrincipalTenantID := ""
	ocmServicePrincipalSubscriptionID := ""

	// AuthFile can be used to get authorizer in a function call
	// for user granted OCM Service Principal
	// The proposed methods offer authorizers based on env or file
	// which cannot be used when calling this in a function in a
	// multitenant environment
	// and yes, this is a workaround
	f := auth.FileSettings{}
	f.Values = map[string]string{}
	f.Values[auth.ClientID] = ocmServicePrincipalClientID
	f.Values[auth.ClientSecret] = ocmServicePrincipalClientSecret
	f.Values[auth.TenantID] = ocmServicePrincipalTenantID
	f.Values[auth.SubscriptionID] = ocmServicePrincipalSubscriptionID

	// Authorize to use against Resource Manager
	authorizer, err := f.ClientCredentialsAuthorizer(azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		panic(err)
	}

	// cluster basic configuration
	clusterName := "ocm-test-cluster"
	clusterLocation := "ueastus"
	clusterSubscriptionId := ""
	clusterResourceGroupName := ""

	// network configuration
	vnetName := ""
	masterSubnetName := ""
	workerSubnetName := ""
	masterSubnetResourceID := fmt.Sprintf("GET https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", clusterSubscriptionId, clusterResourceGroupName, vnetName, masterSubnetName)
	workerSubnetResourceID := fmt.Sprintf("GET https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s", clusterSubscriptionId, clusterResourceGroupName, vnetName, workerSubnetName)
	// default used by `az aro` client
	podCIDR := "10.128.0.0/14"
	serviceCIDR := "172.30.0.0/16"

	// customer created in cluster service principal
	clusterServicePrincipalClientID := ""
	clusterServicePrincipalClientSecret := ""

	ocp := redhatopenshift.OpenShiftCluster{
		OpenShiftClusterProperties: &redhatopenshift.OpenShiftClusterProperties{
			ClusterProfile: &redhatopenshift.ClusterProfile{
				Domain:          to.StringPtr(strings.ToLower(clusterName)),
				PullSecret:      to.StringPtr(""),
				ResourceGroupID: to.StringPtr(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", clusterSubscriptionId, "aro-"+strings.ToLower(clusterName))),
			},
			ServicePrincipalProfile: &redhatopenshift.ServicePrincipalProfile{
				ClientID:     to.StringPtr(clusterServicePrincipalClientID),
				ClientSecret: to.StringPtr(clusterServicePrincipalClientSecret),
			},
			NetworkProfile: &redhatopenshift.NetworkProfile{
				PodCidr:     to.StringPtr(podCIDR),
				ServiceCidr: to.StringPtr(serviceCIDR),
			},
			MasterProfile: &redhatopenshift.MasterProfile{
				VMSize:   redhatopenshift.StandardD8sV3,
				SubnetID: to.StringPtr(masterSubnetResourceID),
			},
			WorkerProfiles: &[]redhatopenshift.WorkerProfile{
				{
					Name:       to.StringPtr("worker"),
					VMSize:     redhatopenshift.VMSize1(redhatopenshift.StandardD4sV3),
					DiskSizeGB: to.Int32Ptr(128),
					Count:      to.Int32Ptr(3),
					SubnetID:   to.StringPtr(workerSubnetResourceID),
				},
			},
			ApiserverProfile: &redhatopenshift.APIServerProfile{
				Visibility: redhatopenshift.Public,
			},
			IngressProfiles: &[]redhatopenshift.IngressProfile{
				{
					Name:       to.StringPtr("default"),
					Visibility: redhatopenshift.Visibility1Public,
				},
			},
		},
		Location: to.StringPtr(clusterLocation),
	}

	client := mgmtredhatopenshift.NewOpenShiftClustersClient(&azureclient.PublicCloud, clusterSubscriptionId, authorizer)

	err = client.CreateOrUpdateAndWait(context.Background(), clusterResourceGroupName, clusterName, ocp)
	if err != nil {
		panic(err)
	}

	err = client.DeleteAndWait(context.Background(), clusterResourceGroupName, clusterName)
	if err != nil {
		panic(err)
	}

	log.Println(ocp)

	log.Println("all done")
}
