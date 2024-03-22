package provider

import (
	openfgaClient "github.com/openfga/go-sdk/client"
)

type OpenfgaProviderData struct {
	clientCache *openfgaClientCache
}

func (d *OpenfgaProviderData) GetGlobalClient() (*openfgaClient.OpenFgaClient, error) {
	return d.clientCache.GetGlobalClient()
}

func (d *OpenfgaProviderData) GetClientForStore(storeId string) (*openfgaClient.OpenFgaClient, error) {
	return d.clientCache.GetClientForStore(storeId)
}
