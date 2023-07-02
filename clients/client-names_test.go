package clients

import (
	"testing"
)

func TestGenerateName(t *testing.T) {
	for i := 0; i < 8; i++ {
		t.Log(GenerateName())
	}
}

// wanted-viper
// fleet-rocket
// accepted-dynamite
// valued-forge
// tops-owl
// informed-owlwoman
// evolving-killer-shrike
// literate-mystique
