package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	netRequestTime = 20
	produceDelay   = 30
)

// Планировщик задач.
type scheduler struct {
	mux   sync.Mutex
	tasks []task
}

// Задание на отправку уведомления.
type task struct {
	t   int64
	msg message // Уведомление UserNotificationV2
}

// Сообщение - 10 КБайт.
type message [10_000]byte

var (
	queueLen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "q_len",
		Help: "Queue length.",
	})
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":80", nil)

	var s scheduler
	go produce(&s)
	go consume(&s)

	for {
		time.Sleep(time.Second * 1)
		s.mux.Lock()
		fmt.Printf("Размер/емкость очереди: %d/%d\n", len(s.tasks), cap(s.tasks))
		queueLen.Set(float64(len(s.tasks)))
		s.mux.Unlock()
	}
}

// Получение заданий на уведомление из внешнего источника.
func produce(s *scheduler) {
	for {
		s.mux.Lock()
		add(s)
		s.mux.Unlock()
		time.Sleep(time.Millisecond * produceDelay)
	}
}

// Выполнение заданий.
func consume(s *scheduler) {
	for {
		now := time.Now().Unix()
		s.mux.Lock()
		for i := len(s.tasks) - 1; i >= 0; i-- {
			if s.tasks[i].t <= now {
				exec(s.tasks[i])
				remove(s, i)
			}
		}
		s.mux.Unlock()
	}
}

func exec(task) {
	time.Sleep(time.Millisecond * netRequestTime)
}

func add(s *scheduler) {
	t := task{
		t:   time.Now().Unix() + rand.Int63n(120),
		msg: message{},
	}
	s.tasks = append(s.tasks, t)
}

func remove(s *scheduler, i int) {
	s.tasks[i] = s.tasks[len(s.tasks)-1]
	s.tasks = s.tasks[:len(s.tasks)-1]
}
