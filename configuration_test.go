package configurationmanager

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "",
			run:  func(t *testing.T) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
