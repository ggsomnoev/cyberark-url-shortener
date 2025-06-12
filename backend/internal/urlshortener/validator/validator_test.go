package validator_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/validator"
)

var _ = Describe("Validator", func() {
	var v validator.Validator

	BeforeEach(func() {
		v = *validator.NewValidator()
	})

	Describe("ValidateShortenRequest", func() {
		It("should succeeds", func() {
			req := model.ShortenRequest{
				URL:       "https://example.com",
				ShortCode: "abcDEF12",
			}
			err := v.ValidateShortenRequest(req)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error for an invalid URL", func() {
			req := model.ShortenRequest{
				URL:       "not_a_url",
				ShortCode: "abcDEF12",
			}
			err := v.ValidateShortenRequest(req)
			Expect(err).To(MatchError(validator.ErrInvalidURL))
		})

		It("should return error for invalid short code", func() {
			req := model.ShortenRequest{
				URL:       "https://example.com",
				ShortCode: "invalid_short_code",
			}
			err := v.ValidateShortenRequest(req)
			Expect(err).To(MatchError(validator.ErrInvalidShortCode))
		})

		It("should return nil if short code is empty (optional)", func() {
			req := model.ShortenRequest{
				URL:       "https://example.com",
				ShortCode: "",
			}
			err := v.ValidateShortenRequest(req)
			Expect(err).To(BeNil())
		})

		It("should return error if URL is empty", func() {
			req := model.ShortenRequest{
				URL:       "",
				ShortCode: "abcdefgh",
			}
			err := v.ValidateShortenRequest(req)
			Expect(err).To(MatchError(validator.ErrInvalidURL))
		})
	})

	Describe("ValidateURL", func() {
		It("should succeed", func() {
			err := v.ValidateURL("https://example.com/path")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error for a URL missing scheme", func() {
			err := v.ValidateURL("example.com")
			Expect(err).To(MatchError(validator.ErrInvalidURL))
		})

		It("should return error for a URL missing host", func() {
			err := v.ValidateURL("https:///path")
			Expect(err).To(MatchError(validator.ErrInvalidURL))
		})

		It("should return error for completely invalid URL", func() {
			err := v.ValidateURL("not a url")
			Expect(err).To(MatchError(validator.ErrInvalidURL))
		})
	})

	Describe("ValidateShortCode", func() {
		It("should succeed", func() {
			err := v.ValidateShortCode("AbcD123_")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return error for short code that is too short", func() {
			err := v.ValidateShortCode("abc")
			Expect(err).To(MatchError(validator.ErrInvalidShortCode))
		})

		It("should return error for short code with invalid characters", func() {
			err := v.ValidateShortCode("abc@123$")
			Expect(err).To(MatchError(validator.ErrInvalidShortCode))
		})

		It("should return error for short code that is too long", func() {
			err := v.ValidateShortCode("abcd12345X")
			Expect(err).To(MatchError(validator.ErrInvalidShortCode))
		})
	})

})
