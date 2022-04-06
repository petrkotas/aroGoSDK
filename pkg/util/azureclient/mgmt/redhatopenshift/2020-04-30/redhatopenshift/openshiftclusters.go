package redhatopenshift

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest"
	mgmtredhatopenshift20200430 "github.com/petrkotas/aroGoSDK/pkg/client/services/redhatopenshift/mgmt/2020-04-30/redhatopenshift"
	"github.com/petrkotas/aroGoSDK/pkg/util/azureclient"
)

// OpenShiftClustersClient is a minimal interface for azure OpenshiftClustersClient
type OpenShiftClustersClient interface {
	ListCredentials(ctx context.Context, resourceGroupName string, resourceName string) (result mgmtredhatopenshift20200430.OpenShiftClusterCredentials, err error)
	Get(ctx context.Context, resourceGroupName string, resourceName string) (result mgmtredhatopenshift20200430.OpenShiftCluster, err error)
	OpenShiftClustersClientAddons
}

type openShiftClustersClient struct {
	mgmtredhatopenshift20200430.OpenShiftClustersClient
}

var _ OpenShiftClustersClient = &openShiftClustersClient{}

// NewOpenShiftClustersClient creates a new OpenShiftClustersClient
func NewOpenShiftClustersClient(environment *azureclient.AROEnvironment, subscriptionID string, authorizer autorest.Authorizer) OpenShiftClustersClient {
	var client mgmtredhatopenshift20200430.OpenShiftClustersClient
	client = mgmtredhatopenshift20200430.NewOpenShiftClustersClientWithBaseURI(environment.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = authorizer
	client.PollingDelay = 10 * time.Second
	client.PollingDuration = 2 * time.Hour

	return &openShiftClustersClient{
		OpenShiftClustersClient: client,
	}
}
