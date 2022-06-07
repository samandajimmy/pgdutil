package pgdutil_test

import (
	"errors"
	"net/http"

	"github.com/samandajimmy/pgdutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type tmp struct {
	CIF string `json:"cif" validate:"required"`
}

var _ = Describe("Handler", func() {
	var iHandler pgdutil.IHandler
	var stHandler pgdutil.Handler
	var e pgdutil.DummyEcho
	var reqpl map[string]interface{}
	var pl interface{}
	var err error

	JustBeforeEach(func() {
		e = pgdutil.NewDummyEcho(http.MethodPost, "/", reqpl)
		iHandler = pgdutil.NewHandler(&stHandler)
	})

	BeforeEach(func() {
		stHandler = pgdutil.Handler{}
	})

	Describe("Validate", func() {
		JustBeforeEach(func() {
			err = iHandler.Validate(e.Context, pl)
		})

		Context("error echo binding", func() {
			BeforeEach(func() {
				reqpl = map[string]interface{}{"cif": "11223344"}
				pl = &[]struct{ Field string }{}
			})

			It("expect to return error", func() {
				Expect(err.Error()).To(Equal("code=400, message=Unmarshal type error: expected=[]struct { Field string }, got=object, field=, offset=1, internal=json: cannot unmarshal object into Go value of type []struct { Field string }"))
			})
		})

		Context("error echo validate", func() {
			BeforeEach(func() {
				reqpl = map[string]interface{}{"isError": true}
				pl = map[string]interface{}{"cif": "11223344"}
			})

			It("expect to return error", func() {
				Expect(err).To(Equal(pgdutil.ErrInternalServerError))
			})
		})

		Context("succeeded", func() {
			BeforeEach(func() {
				reqpl = map[string]interface{}{"cif": "11223344"}
				pl = tmp{CIF: "1122334455"}
			})

			It("expect to return nil error", func() {
				Expect(err).To(BeNil())
			})

			It("expect to reset struct handler", func() {
				Expect(stHandler).To(Equal(pgdutil.Handler{}))
			})
		})
	})

	Describe("ShowResponse", func() {
		var respData tmp
		var mockResp pgdutil.Response
		var errInput error
		var errs pgdutil.ResponseErrors

		JustBeforeEach(func() {
			err = iHandler.ShowResponse(e.Context, respData, errInput, errs)
		})

		BeforeEach(func() {
			respData = tmp{}
			mockResp = pgdutil.Response{}
			errInput = nil
			errs = pgdutil.ResponseErrors{}
		})

		Context("error input is nil", func() {
			Context("response errs is nil", func() {
				BeforeEach(func() {
					respData = tmp{
						CIF: "jimmy",
					}
					mockResp = pgdutil.Response{
						Code:    "00",
						Status:  "Success",
						Message: "Data Berhasil Dikirim",
						Data:    respData,
					}
				})

				It("expect to have response data ", func() {
					Expect(stHandler.Response).To(Equal(mockResp))
				})

				It("expect to have nil error ", func() {
					Expect(err).To(BeNil())
				})
			})

			Context("response errs is not nil", func() {
				BeforeEach(func() {
					errs = pgdutil.ResponseErrors{Title: "This is errs"}
					mockResp = pgdutil.Response{
						Code:    "99",
						Status:  "Error",
						Message: "This is errs",
					}

				})

				It("expect to show response errs ", func() {
					Expect(stHandler.Response).To(Equal(mockResp))
				})

				It("expect to have nil error ", func() {
					Expect(err).To(BeNil())
				})
			})

		})

		Context("error input is not nil", func() {
			BeforeEach(func() {
				errInput = errors.New("cacing")
				mockResp = pgdutil.Response{
					Code:    "99",
					Status:  "Error",
					Message: errInput.Error(),
				}
			})

			It("expect to have response error", func() {
				Expect(stHandler.Response).To(Equal(mockResp))
			})

			It("expect to return nil error", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("SetTotalCount", func() {
		JustBeforeEach(func() {
			iHandler.SetTotalCount("100")
		})

		It("expect to have totalCount on response total count", func() {
			Expect(stHandler.Response.TotalCount).To(Equal("100"))
		})
	})
})
