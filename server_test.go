package server

import (
	"github.com/microbay/microbay/model"
	"testing"
)

func generateTestModel() []*model.Resource {
	resources := make([]*model.Resource, 2)
	micros := make([]model.Micro, 5)
	micros[0] = model.Micro{"http://localhost:9000", 1}
	micros[1] = model.Micro{"http://localhost:9001", 33}
	micros[2] = model.Micro{"http://localhost:9002", 2}
	micros[3] = model.Micro{"http://localhost:9003", 17}
	micros[4] = model.Micro{"http://localhost:9004", 40}
	resources[0] = &model.Resource{
		Micros: micros,
	}
	micros = make([]model.Micro, 2)
	micros[0] = model.Micro{"http://localhost:9000", 0}
	micros[1] = model.Micro{"http://localhost:9001", 1000}
	resources[1] = &model.Resource{
		Micros: micros,
	}
	return resources
}

func TestBootstrapLoadBalancer(t *testing.T) {
	resources := generateTestModel()
	bootstrapLoadBalancer(resources)
	resource := resources[0]
	if resource.Backends == nil {
		t.Error(
			"expected", "instance of backend.Backends",
			"got", nil,
		)
		t.FailNow()
	}
	if resource.Backends.Len() != 93 {
		t.Error(
			"expected", 93,
			"got", resource.Backends.Len(),
		)
		t.FailNow()
	}
	resource = resources[1]
	if resource.Backends == nil {
		t.Error(
			"expected", "instance of backend.Backends",
			"got", nil,
		)
		t.FailNow()
	}
	if resource.Backends.Len() != 1000 {
		t.Error(
			"expected", 1000,
			"got", resource.Backends.Len(),
		)
		t.FailNow()
	}

	b := resource.Backends
	i := 0
	for e := b.Choose(); e != nil; e = b.Choose() {
		if i > 1001 {
			break
		}
		if e.String() == "http://localhost:9000" {
			t.Error(
				"expected", "http://localhost:9001",
				"got", "http://localhost:9000",
			)
			t.FailNow()
		}
		i++
	}
}
