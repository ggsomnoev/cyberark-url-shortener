package store_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testdb "github.com/ggsomnoev/cyberark-url-shortener/test/pg"
)

var (
	ctx  context.Context
	pool *pgxpool.Pool
)

func TestStore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store Suite")
}

var _ = BeforeEach(func() {
	ctx = context.Background()
	pool = testdb.MustInitDBConnectio(ctx)
})

var _ = AfterEach(func() {
	pool.Close()
})
