package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/gomega"
)

var testDBDSN = "postgres://user:pass@localhost:5432/urlshortenerdb"

func MustInitDBConnectio(ctx context.Context) *pgxpool.Pool {
	poolCfg, err := pgxpool.ParseConfig(testDBDSN)
	Expect(err).NotTo(HaveOccurred())

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	Expect(err).NotTo(HaveOccurred())

	Expect(pool.Ping(ctx)).To(Succeed())

	return pool
}
