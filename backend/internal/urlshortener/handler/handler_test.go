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
	ErrShortenURLFailed = errors.New("failed to shorten URL")
	ErrUrlResolveFailed = errors.New("failed to resolve URL")
)

var _ = Describe("Handler", func() {
	var (
		ctx      context.Context
		e        *echo.Echo
		recorder *httptest.ResponseRecorder
		service  *handlerfakes.FakeService
	)

	BeforeEach(func() {
		e = echo.New()
		ctx = context.Background()
		service = &handlerfakes.FakeService{}
		recorder = httptest.NewRecorder()

		handler.RegisterRoutes(ctx, e, service)
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

			service.ShortenURLReturns(shortCode, nil)
		})

		JustBeforeEach(func() {
			e.ServeHTTP(recorder, req)
		})

		It("succeeds", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))

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

			service.ResolveURLReturns(url, nil)
		})

		JustBeforeEach(func() {
			e.ServeHTTP(recorder, req)
		})

		It("succeeds", func() {
			Expect(recorder.Code).To(Equal(http.StatusFound))

			Expect(service.ResolveURLCallCount()).To(Equal(1))
			actualCtx, actualShortUrl := service.ResolveURLArgsForCall(0)
			Expect(actualCtx).To(Equal(ctx))
			Expect(actualShortUrl).To(Equal(shortCode))
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
