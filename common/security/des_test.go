package security

import (
	"testing"
)

func Test3DesEncrypt(t *testing.T) {
	cardNo := "6236216615490000006"
	encrypted := "YbAkr7OkZzY8Z3KCJYTQtwlTaY96GDFf"
	key := "012345678901234567890123"

	res, err := TripleEcbDesEncryptBase64(cardNo, key)
	if err != nil {
		t.Fail()
	}

	t.Log("encrypted:", res, "expected:", encrypted)

	res, err = TripleEcbDesDecryptBase64(encrypted, key)
	if err != nil {
		t.Fail()
	}

	t.Log("origin:", res)
}
