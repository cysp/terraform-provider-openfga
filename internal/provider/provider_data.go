package provider

import (
	openfgaClient "github.com/openfga/go-sdk/client"
)

type OpenfgaProviderData struct {
	client *openfgaClient.OpenFgaClient
}
