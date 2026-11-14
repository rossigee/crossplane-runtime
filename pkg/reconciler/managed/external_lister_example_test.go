/*
Copyright 2026 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package managed_test

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
)

// zoneClient is an ExternalClient for a kind scoped by "zone" that also
// implements resource.ExternalLister. Observe, Create, Update and Delete are
// left unimplemented; only discovery is shown here.
type zoneClient struct {
	managed.ExternalClientFns

	records map[string][]string // external names keyed by zone
}

func (c zoneClient) List(_ context.Context, _ resource.ProviderConfig, scope resource.ExternalListScope, _ string) (resource.ExternalListResult, error) {
	zone, ok := scope["zone"]
	if !ok {
		return resource.ExternalListResult{}, resource.ErrScopeRequired
	}
	var items []resource.ExternalListItem
	for _, name := range c.records[zone] {
		items = append(items, resource.ExternalListItem{
			ExternalName: name,
			Scope:        resource.ExternalListScope{"zone": zone},
		})
	}
	return resource.ExternalListResult{Items: items}, nil
}

func ExampleExternalLister() {
	var ec managed.ExternalClient = zoneClient{
		records: map[string][]string{
			"example.org": {"www", "mail"},
		},
	}

	// Discovery is optional, so callers check for it with a type assertion.
	lister, ok := ec.(resource.ExternalLister)
	if !ok {
		fmt.Println("provider does not support discovery")
		return
	}

	res, err := lister.List(context.Background(), &fake.ProviderConfig{}, resource.ExternalListScope{"zone": "example.org"}, "")
	if err != nil {
		fmt.Println("list failed:", err)
		return
	}
	for _, item := range res.Items {
		fmt.Println(item.ExternalName, item.Scope["zone"])
	}

	// Kinds that need a scope reject a call without one.
	if _, err := lister.List(context.Background(), &fake.ProviderConfig{}, nil, ""); err != nil {
		fmt.Println("without scope:", err)
	}

	// Output:
	// www example.org
	// mail example.org
	// without scope: external resource kind requires a scope
}
