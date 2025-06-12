package store_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/model"
	"github.com/ggsomnoev/cyberark-url-shortener/internal/urlshortener/store"
)

var _ = Describe("Store", func() {
	When("creating", func() {
		It("should create a store instance", func() {
			Expect(store.NewStore(nil)).NotTo(BeNil())
		})
	})

	Describe("Instance", func() {
		var (
			st *store.Store
		)

		BeforeEach(func() {
			st = store.NewStore(pool)
		})

		Describe("Save", func() {
			var (
				id             uuid.UUID
				url, shortCode string
				urlEntity      model.URL

				errAction error
			)

			BeforeEach(func() {
				id = uuid.New()
				url = "https://www.example.com"
				shortCode = uuid.New().String()[:8]
				urlEntity = model.URL{
					ID:        id,
					Original:  url,
					ShortCode: shortCode,
				}
			})

			JustBeforeEach(func() {
				errAction = st.Save(ctx, urlEntity)
			})

			JustAfterEach(func() {
				Expect(st.DeleteUrlEntityByID(ctx, id)).To(Succeed())
			})

			It("succeeds", func() {
				Expect(errAction).NotTo(HaveOccurred())
			})

			It("inserts the correct url entity", func() {
				storeUrlEntity, err := st.FindByShortCode(ctx, shortCode)
				Expect(err).NotTo(HaveOccurred())
				Expect(storeUrlEntity.ID).To(Equal(id))
				Expect(storeUrlEntity.Original).To(Equal(url))
				Expect(storeUrlEntity.ShortCode).To(Equal(shortCode))
				Expect(storeUrlEntity.CreatedAt).NotTo(BeZero())
			})
		})

		Describe("FindByShortCode", func() {
			var (
				id                        uuid.UUID
				url, shortCode            string
				storeURLEntity, urlEntity model.URL

				errAction error
			)

			BeforeEach(func() {
				id = uuid.New()
				url = "https://www.example.com"
				shortCode = uuid.New().String()[:8]
				urlEntity = model.URL{
					ID:        id,
					Original:  url,
					ShortCode: shortCode,
				}
			})

			JustBeforeEach(func() {
				storeURLEntity, errAction = st.FindByShortCode(ctx, shortCode)
			})

			It("does not find any records", func() {
				Expect(errAction).To(MatchError(store.ErrRecordNotFound))
			})

			When("a record is added", func() {
				BeforeEach(func() {
					Expect(st.Save(ctx, urlEntity)).To(Succeed())
				})

				AfterEach(func() {
					Expect(st.DeleteUrlEntityByID(ctx, id))
				})

				It("succeeds", func() {
					Expect(errAction).NotTo(HaveOccurred())
				})

				It("finds the record", func() {
					Expect(storeURLEntity.ID).To(Equal(id))
					Expect(storeURLEntity.Original).To(Equal(url))
					Expect(storeURLEntity.ShortCode).To(Equal(shortCode))
					Expect(storeURLEntity.CreatedAt).ToNot(BeZero())
				})
			})
		})
	})
})
