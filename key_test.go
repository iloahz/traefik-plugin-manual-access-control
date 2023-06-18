package traefikpluginmanualaccesscontrol

import "testing"

func TestGenerateToken(t *testing.T) {
	key, err := NewKey("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	token := key.GenerateToken()
	t.Log(token)
}

func TestValidateToken(t *testing.T) {
	key, err := NewKey("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	token := key.GenerateToken()
	valid := key.ValidateToken(token)
	if !valid {
		t.Fatal("token not valid")
	}
}

func TestInvalidToken(t *testing.T) {
	key1, err := NewKey("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	key2, err := NewKey("YPsntsoaD3yzmj1OQMrLk51xUCGq5kw6c1in7Xffx0s=")
	if err != nil {
		t.Fatal(err)
	}
	token := key1.GenerateToken()
	valid := key2.ValidateToken(token)
	if valid {
		t.Fatal("token is valid")
	}
}
