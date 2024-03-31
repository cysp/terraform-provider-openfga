package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestProtocol6ProviderServerSchemaVersion(t *testing.T) {
	p := New("test")()

	ps, err := providerserver.NewProtocol6WithError(p)()
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	schema, err := ps.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	assert.NotNil(t, schema.Provider)
	assert.EqualValues(t, 0, schema.Provider.Version)
}
