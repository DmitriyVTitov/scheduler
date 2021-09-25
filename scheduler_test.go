package main

import "testing"

func Benchmark_add(b *testing.B) {
	s := scheduler{}
	for i := 0; i < b.N; i++ {
		add(&s)
	}
}

func Benchmark_remove(b *testing.B) {
	s := scheduler{}
	for i := 0; i < 10_000; i++ {
		add(&s)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		add(&s)
		b.StartTimer()
		remove(&s, 0)
	}
}
