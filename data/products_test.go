package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{Name: "xxx", Price: 2.0, SKU: "abc-abc-abc"}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
