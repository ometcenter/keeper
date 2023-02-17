package healthcheck

import (
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/ometcenter/keeper/config"
	shareStore "github.com/ometcenter/keeper/store"
)

func StartHealthCheck() {

	health := healthcheck.NewHandler()

	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-dep-dns",
		healthcheck.DNSResolveCheck("mos.ru", 50*time.Millisecond))

	DB, err := shareStore.GetDB(config.Conf.DatabaseURL)
	if err != nil {
		//return nil, err
	}

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(DB, 5*time.Second))

	go http.ListenAndServe("0.0.0.0:8086", health)

}

func StartHealthCheckLight() {

	health := healthcheck.NewHandler()

	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(300))

	// // Our app is not ready if we can't resolve our upstream dependency in DNS.
	// health.AddReadinessCheck(
	// 	"upstream-dep-dns",
	// 	healthcheck.DNSResolveCheck("4444google777.com", 50*time.Millisecond)) //mos.ru

	// DB, err := shareStore.GetDB(config.Conf.DatabaseURL)
	// if err != nil {
	// 	//return nil, err
	// }

	// // Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	// health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(DB, 5*time.Second))

	go http.ListenAndServe("0.0.0.0:8086", health)

}

func StartHealthCheckTestError() {

	health := healthcheck.NewHandler()

	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(300))

	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-dep-dns",
		healthcheck.DNSResolveCheck("4444google777.com", 50*time.Millisecond)) //mos.ru

	DB, err := shareStore.GetDB(config.Conf.DatabaseURL)
	if err != nil {
		//return nil, err
	}

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(DB, 5*time.Second))

	go http.ListenAndServe("0.0.0.0:8086", health)

}
