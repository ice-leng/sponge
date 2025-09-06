package rails

import (
	"fmt"
	"testing"
)

func TestDecodeSignedCookie_MissingSecret(t *testing.T) {
	result, err := DecodeSignedCookie("b1870c9c2d472d577b91a25f3ae9daa626725afffa70876d2fd9e004720e9a4f822bdcf0ddc07f3c54ae110d9ff852d5b5f648be56a275338f028287f90e8a85",
		"fpkxE8E9Xksk0W2YXDwXAUhluSaIjMfaKzhII7cAzlwU+hG+7p6nNld+JCa7JyA18Zcl+TvDFJiFS5vRh46PRj6LhmUuxti5PdMH2oPM7UiyllHVcveJcm2ucqZokgx6cMCtrcXfAg+2D3L74JlYvJ9iy6M2mpA1oDCg5jfosvMm8GD0QZfh/DSLjqlZdMUA9S/hcjhak20sG5ZOsq/E9jMnH3DYQoMCxa1oaa+pGcZOcjAkxMFx0FkKjvCGbw9iRO/J0Y8XBBuOrNVBp4U+Zyz4U739RvlO3cG7Odk9s3MCUC+WRw8juIkJ9EMUWJwmIc5uJILZimSdVfwh+Qoj7lEZzwdGw6pFTA91pYpGeUuC1sxnLmIQCUYeoamevPwfFa/tN+eAWZuLq2iAlGWQUf70ECUakrGef6k5JME=--Fgbc3j45HzLzebZK--FzkwNBBEImauLsbCzdz/TA==",
		"_coreui_pro_rails_starter_session")
	if err != nil {
		t.Fatalf("expected error for missing secretKeyBase")
	}
	fmt.Println(result)
}
