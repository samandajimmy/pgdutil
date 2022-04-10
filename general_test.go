package pgdutil_test

import (
	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"repo.pegadaian.co.id/ms-pds/modules/pgdutil"
)

type tmp1 struct {
	Cif         string `json:"cif"`
	ProductCode string `json:"product_code"`
}

var _ = Describe("General", func() {

	Describe("InterfaceToMap", func() {
		It("expected to return as expected", func() {
			usedCif := gofakeit.Regex("[1234567890]{10}")
			obj := tmp1{
				Cif:         usedCif,
				ProductCode: "02",
			}

			mappedObj := pgdutil.InterfaceToMap(obj)
			expectedResult := map[string]interface{}{
				"cif":          usedCif,
				"product_code": "02",
			}

			Expect(mappedObj).To(Equal(expectedResult))
		})
	})
})
