package rmqprom

import (
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	labelQueue = "queue"
	waitTime   = 5 * time.Second
)

func RecordRmqMetrics(connection rmq.Connection, namespace string) {
	readyCount, rejectedCount, connectionCount, consumerCount, unackedCount := registerCounters(namespace)

	go func() {
		for {
			queues, err := connection.GetOpenQueues()
			if err != nil {
				logrus.Warnf("error fetching open queues: %s", err.Error())
				time.Sleep(waitTime)
				continue
			}

			stats, err := connection.CollectStats(queues)
			if err != nil {
				logrus.Warnf("error collecting stats from open queues: %s", err.Error())
				time.Sleep(waitTime)
				continue
			}

			for queue, queueStats := range stats.QueueStats {
				set(readyCount, queue, float64(queueStats.ReadyCount))
				set(rejectedCount, queue, float64(queueStats.RejectedCount))
				set(connectionCount, queue, float64(queueStats.ConnectionCount()))
				set(consumerCount, queue, float64(queueStats.ConsumerCount()))
				set(unackedCount, queue, float64(queueStats.UnackedCount()))
			}

			time.Sleep(waitTime)
		}
	}()
}

func set(gaugeVec *prometheus.GaugeVec, queue string, value float64) {
	gauge, err := gaugeVec.GetMetricWith(prometheus.Labels{labelQueue: queue})
	if err != nil {
		logrus.Warnf("error sending metric: %s, label: %s", err.Error(), queue)
		return
	}
	gauge.Set(value)
}

func registerCounters(namespace string) (readyCount, rejectedCount, connectionCount, consumerCount, unackedCount *prometheus.GaugeVec) {
	namespace = namespace + "_rmq"
	labels := []string{labelQueue}

	readyCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ready",
			Help:      "Number of ready messages on queue",
		},
		labels,
	)

	rejectedCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "rejected",
			Help:      "Number of rejected messages on queue",
		},
		labels,
	)

	connectionCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "connection",
			Help:      "Number of connections consuming a queue",
		},
		labels,
	)

	consumerCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "consumer",
			Help:      "Number of consumers consuming messages for a queue",
		},
		labels,
	)

	unackedCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "unacked",
			Help:      "Number of unacked messages on a consumer",
		},
		labels,
	)

	prometheus.MustRegister(readyCount, rejectedCount, connectionCount, consumerCount, unackedCount)

	return readyCount, rejectedCount, connectionCount, consumerCount, unackedCount
}
