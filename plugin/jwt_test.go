package plugin

import "testing"

func TestGenerateToken(t *testing.T) {
	j, err := NewJWT("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	token := j.GenerateToken()
	t.Log(token)
}

func TestValidateToken(t *testing.T) {
	j, err := NewJWT("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	token := j.GenerateToken()
	claims, err := j.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)
}

func TestInvalidToken(t *testing.T) {
	j1, err := NewJWT("8HvPVByYKFAt16+qG5/ZDgV11iVEDPOxVrM+caF81jA=")
	if err != nil {
		t.Fatal(err)
	}
	j2, err := NewJWT("YPsntsoaD3yzmj1OQMrLk51xUCGq5kw6c1in7Xffx0s=")
	if err != nil {
		t.Fatal(err)
	}
	token := j1.GenerateToken()
	_, err = j2.ValidateToken(token)
	if err == nil {
		t.Fatal("token is valid")
	}
}
