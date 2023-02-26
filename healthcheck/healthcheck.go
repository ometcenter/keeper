package healthcheck

import (
	"fmt"
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
	//health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))

	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-dep-dns",
		healthcheck.DNSResolveCheck("mos.ru", 50*time.Millisecond))

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

func HealthCheckTestGoroutines() {

	// vars := mux.Vars(r)
	// JobId := vars["id"]

	// err := store.CancelJob(JobId)

	// if err != nil {
	// 	log.Impl.Error(err)
	// 	fmt.Fprintf(w, err.Error())
	// 	return
	// }

	jobs := make(chan int, 5)
	done := make(chan bool)

	go func() {
		for {
			j, closed := <-jobs
			if closed {
				fmt.Println("received job", j)
			} else {
				fmt.Println("received all jobs")
				done <- true
				return
			}
		}
	}()

	for i := 0; i < 320; i++ {
		go func(counter int) {
			for j := 0; j < 100; j++ {

				// select {
				// case msg1 := <-done:
				// 	fmt.Println("received done", msg1)
				// 	if msg1 {
				// 		return
				// 	}
				// default:
				jobs <- counter
				time.Sleep(1 * time.Second)
				// }

			}
		}(i)
	}

	go func() {
		time.Sleep(120 * time.Second)
		fmt.Println("sent all jobs")
		done <- true
		close(jobs)
	}()

}
