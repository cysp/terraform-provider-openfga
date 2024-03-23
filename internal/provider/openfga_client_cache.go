package provider

import (
	openfgaClient "github.com/openfga/go-sdk/client"
)

type openfgaClientCache struct {
	apiUrl           string
	clientsByStoreId map[string]*openfgaClient.OpenFgaClient
}

func NewOpenfgaClientCache(apiUrl string) *openfgaClientCache {
	return &openfgaClientCache{
		apiUrl:           apiUrl,
		clientsByStoreId: make(map[string]*openfgaClient.OpenFgaClient),
	}
}

func (c *openfgaClientCache) GetGlobalClient() (*openfgaClient.OpenFgaClient, error) {
	return c.GetClientForStore("")
}

func (c *openfgaClientCache) GetClientForStore(storeId string) (*openfgaClient.OpenFgaClient, error) {
	if client, ok := c.clientsByStoreId[storeId]; ok {
		return client, nil
	}

	client, err := openfgaClient.NewSdkClient(&openfgaClient.ClientConfiguration{
		ApiUrl:  c.apiUrl,
		StoreId: storeId,
	})
	if err != nil {
		return nil, err
	}

	c.clientsByStoreId[storeId] = client

	return client, nil
}
