package kafka

import (
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func ConnectKafkaGoKafkaPackage() error {

	// mechanism, err := scram.Mechanism(scram.SHA256, "gokeeper", "1LTRDlfP")
	// if err != nil {
	// 	panic(err)
	// }

	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_USERNAME"),
		Password: os.Getenv("KAFKA_PASSWORD"),
		//Password: "gokeeper",
	}

	// // rootCAs, _ := x509.SystemCertPool()
	// // if rootCAs == nil {
	// // 	rootCAs = x509.NewCertPool()
	// // }

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		//ClientID:      "gokeeper",
		SASLMechanism: mechanism,
		// TLS: &tls.Config{
		// 	InsecureSkipVerify: false,
		// 	RootCAs:            rootCAs},
		//TLS: &tls.Config{InsecureSkipVerify: true},
	}

	_ = dialer

	return nil
}
