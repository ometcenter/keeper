package tracing

// Код основан на трех примерх-пакетах
// https://logz.io/blog/go-instrumentation-distributed-tracing-jaeger/
// https://github.com/emailtovamos/JaegerQuickExample
// https://github.com/opentracing-contrib/go-amqp

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

// Init returns an instance of Jaeger Tracer.
func Init(serviceName string, addressPostJaeger string) (io.Closer, error) {

	// Вариант с загрузкой из глобальных переменных
	// cfg, err := config.FromEnv()
	// if err != nil {
	// 	panic(fmt.Sprintf("ERROR: failed to read config from env vars: %v\n", err))
	// }
	// tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	// if err != nil {
	// 	panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	// }

	// Tracer = tracer
	// return tracer, closer

	// Вариант попроще без фиксированных переменных среды
	//Tracer = opentracing.GlobalTracer()
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: addressPostJaeger, // Указывать порт с данными, а не с веб-мордой
			LogSpans:           true,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	//defer closer.Close()
	if err != nil {
		return nil, err
		//panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	Tracer = tracer
	opentracing.SetGlobalTracer(Tracer)

	return closer, nil
}
