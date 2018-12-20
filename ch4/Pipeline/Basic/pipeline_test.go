package main

import "testing"

func BenchmarkSimplePipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		simplePipeline()
	}
}

func BenchmarkStreamPipeline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		streamHandlingWithPipeline()
	}
}