package redhatopenshift

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"github.com/Azure/go-autorest/autorest"

	mgmtredhatopenshift20200430 "github.com/petrkotas/aroGoSDK/pkg/client/services/redhatopenshift/mgmt/2020-04-30/redhatopenshift"
	"github.com/petrkotas/aroGoSDK/pkg/util/azureclient"
)

// OperationsClient is a minimal interface for azure OperationsClient
type OperationsClient interface {
	OperationsClientAddons
}

type operationsClient struct {
	mgmtredhatopenshift20200430.OperationsClient
}

var _ OperationsClient = &operationsClient{}

// NewOperationsClient creates a new OperationsClient
func NewOperationsClient(environment *azureclient.AROEnvironment, subscriptionID string, authorizer autorest.Authorizer) OperationsClient {
	var client mgmtredhatopenshift20200430.OperationsClient
	client = mgmtredhatopenshift20200430.NewOperationsClientWithBaseURI(environment.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = authorizer

	return &operationsClient{
		OperationsClient: client,
	}
}
