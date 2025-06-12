package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/handler"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/handler/handlerfakes"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
)

var (
	ErrShortenRequestValidationFailed = errors.New("failed to validate request")
	ErrShortenURLFailed               = errors.New("failed to shorten URL")
	ErrShortCodeValidationFailed      = errors.New("failed to validate short code")
	ErrUrlResolveFailed               = errors.New("failed to resolve URL")
)

var _ = Describe("Handler", func() {
	var (
		ctx       context.Context
		e         *echo.Echo
		recorder  *httptest.ResponseRecorder
		service   *handlerfakes.FakeService
		validator *handlerfakes.FakeValidator
	)

	BeforeEach(func() {
		e = echo.New()
		ctx = context.Background()
		service = &handlerfakes.FakeService{}
		validator = &handlerfakes.FakeValidator{}
		recorder = httptest.NewRecorder()

		handler.RegisterRoutes(ctx, e, service, validator)
	})

	Describe("/api/shorten", func() {
		var (
			req       *http.Request
			url       string
			shortCode string
		)

		BeforeEach(func() {
			url = "https://www.example.com"
			shortCode = uuid.New().String()[:8]

			body, err := json.Marshal(model.ShortenRequest{URL: url})
			Expect(err).NotTo(HaveOccurred())

			req = httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			validator.ValidateShortenRequestReturns(nil)
			service.ShortenURLReturns(shortCode, nil)
		})

		JustBeforeEach(func() {
			e.ServeHTTP(recorder, req)
		})

		It("succeeds", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))

			Expect(validator.ValidateShortenRequestCallCount()).To(Equal(1))
			Expect(service.ShortenURLCallCount()).To(Equal(1))
			actualCtx, actualURL := service.ShortenURLArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualURL).To(Equal(url))
		})

		Context("invalid json is posted", func() {
			BeforeEach(func() {
				req = httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString("{invalid json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			})

			It("returns 400", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("url validation fails", func() {
			BeforeEach(func() {
				validator.ValidateShortenRequestReturns(ErrShortenRequestValidationFailed)
			})

			It("returns 400", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).To(ContainSubstring(ErrShortenRequestValidationFailed.Error()))
			})
		})

		Context("url shorten fails", func() {
			BeforeEach(func() {
				service.ShortenURLReturns("", ErrShortenURLFailed)
			})

			It("returns 500", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring("could not shorten URL"))
			})
		})
	})

	Describe("/:code", func() {
		var (
			req       *http.Request
			url       string
			shortCode string
		)

		BeforeEach(func() {
			url = "https://www.example.com"
			shortCode = uuid.New().String()[:8]

			req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", shortCode), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			validator.ValidateShortCodeReturns(nil)
			service.ResolveURLReturns(url, nil)
		})

		JustBeforeEach(func() {
			e.ServeHTTP(recorder, req)
		})

		It("succeeds", func() {
			Expect(recorder.Code).To(Equal(http.StatusFound))

			Expect(validator.ValidateShortCodeCallCount()).To(Equal(1))
			Expect(service.ResolveURLCallCount()).To(Equal(1))
			actualCtx, actualShortUrl := service.ResolveURLArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualShortUrl).To(Equal(shortCode))
		})

		Context("shortCode validation fails", func() {
			BeforeEach(func() {
				validator.ValidateShortCodeReturns(ErrShortCodeValidationFailed)
			})

			It("returns 400", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(recorder.Body.String()).To(ContainSubstring(ErrShortCodeValidationFailed.Error()))
			})
		})

		Context("resolve url fails", func() {
			BeforeEach(func() {
				service.ResolveURLReturns("", ErrUrlResolveFailed)
			})

			It("returns 500", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(recorder.Body.String()).To(ContainSubstring(ErrUrlResolveFailed.Error()))
			})
		})
	})
})
