package db

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test cases include some initial sanity checks, followed by the ones from the spec
// I then have some benchmark tests to see if I could prove the O log n complexity
// and after graphing them I believe that I accomplished that goal (I did not check
// inside Transactions, but it should have a matching big O to non-transaction operations)

func TestGet(t *testing.T) {
	d := Init()
	d.Database.ReplaceOrInsert(Node{key: "hi", value: "hello"})
	d.Counts.ReplaceOrInsert(CountNode{key: "hello", value: 1})
	d.Database.ReplaceOrInsert(Node{key: "bye", value: "seeya"})
	d.Counts.ReplaceOrInsert(CountNode{key: "seeya", value: 1})
	result := d.Get("hi")
	assert.Equal(t, "hello", result)
	result = d.Get("boo")
	assert.Equal(t, "NULL", result)
}

func TestSet(t *testing.T) {
	d := Init()
	d.Database.ReplaceOrInsert(Node{key: "hi", value: "hello"})
	d.Counts.ReplaceOrInsert(CountNode{key: "hello", value: 1})
	d.Database.ReplaceOrInsert(Node{key: "bye", value: "seeya"})
	d.Counts.ReplaceOrInsert(CountNode{key: "seeya", value: 1})
	result := d.Get("hi")
	assert.Equal(t, "hello", result)
	d.Set("hi", "hulloo")
	result = d.Get("hi")
	assert.Equal(t, "hulloo", result)
}

func TestCase1(t *testing.T) {
	d := Init()
	assert.Equal(t, d.Get("a"), "NULL")
	d.Set("a", "foo")
	d.Set("b", "foo")
	assert.Equal(t, 2, d.Count("foo"))
	assert.Equal(t, 0, d.Count("bar"))
	d.Delete("a")
	assert.Equal(t, 1, d.Count("foo"), 1)
	d.Set("b", "bar")
	assert.Equal(t, 0, d.Count("foo"), 0)
	assert.Equal(t, "bar", d.Get("b"))
	assert.Equal(t, "NULL", d.Get("B"))
}

func TestCase2(t *testing.T) {
	d := Init()
	d.Set("a", "foo")
	d.Set("a", "foo")
	assert.Equal(t, 1, d.Count("foo"))
	assert.Equal(t, "foo", d.Get("a"))
	d.Delete("a")
	assert.Equal(t, "NULL", d.Get("a"))
	assert.Equal(t, 0, d.Count("foo"))
}

func TestCase3(t *testing.T) {
	d := Init()
	d.Begin()
	d.Set("a", "foo")
	assert.Equal(t, "foo", d.Get("a"))
	d.Begin()
	d.Set("a", "bar")
	assert.Equal(t, "bar", d.Get("a"))
	d.Set("a", "baz")
	assert.Equal(t, "baz", d.Get("a"))
	d.Rollback()
	assert.Equal(t, "foo", d.Get("a"))
	d.Rollback()
	assert.Equal(t, "NULL", d.Get("a"))
}

func TestCase4(t *testing.T) {
	d := Init()
	d.Set("a", "foo")
	d.Set("b", "baz")
	d.Begin()
	assert.Equal(t, "foo", d.Get("a"))
	d.Set("a", "bar")
	assert.Equal(t, 1, d.Count("bar"))
	d.Begin()
	assert.Equal(t, 1, d.Count("bar"))
	d.Delete("a")
	assert.Equal(t, "NULL", d.Get("a"))
	assert.Equal(t, 0, d.Count("bar"))
	d.Rollback()
	assert.Equal(t, "bar", d.Get("a"))
	assert.Equal(t, 1, d.Count("bar"))
	d.Commit()
	assert.Equal(t, "bar", d.Get("a"))
	assert.Equal(t, "baz", d.Get("b"))
}

func benchmarkSet(dbSize int, b *testing.B) {
	d := Init()
	for i := 0; i < dbSize; i++ {
		str := strconv.Itoa(i)
		d.Set(str, str)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d.Set("str", "str")
	}
}

func BenchmarkSet0(b *testing.B)      { benchmarkSet(0, b) }
func BenchmarkSet10(b *testing.B)     { benchmarkSet(10, b) }
func BenchmarkSet100(b *testing.B)    { benchmarkSet(100, b) }
func BenchmarkSet1000(b *testing.B)   { benchmarkSet(1000, b) }
func BenchmarkSet10000(b *testing.B)  { benchmarkSet(10000, b) }
func BenchmarkSet100000(b *testing.B) { benchmarkSet(100000, b) }

// func BenchmarkSet1000000(b *testing.B) { benchmarkSet(1000000, b) }

func benchmarkGet(dbSize int, b *testing.B) {
	d := Init()
	var str string
	for i := 0; i < dbSize; i++ {
		str = strconv.Itoa(i)
		d.Set(str, str)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d.Get(str)
	}
}

func BenchmarkGet0(b *testing.B)       { benchmarkGet(0, b) }
func BenchmarkGet10(b *testing.B)      { benchmarkGet(10, b) }
func BenchmarkGet100(b *testing.B)     { benchmarkGet(100, b) }
func BenchmarkGet1000(b *testing.B)    { benchmarkGet(1000, b) }
func BenchmarkGet10000(b *testing.B)   { benchmarkGet(10000, b) }
func BenchmarkGet100000(b *testing.B)  { benchmarkGet(100000, b) }
func BenchmarkGet1000000(b *testing.B) { benchmarkGet(1000000, b) }

// func BenchmarkGet10000000(b *testing.B) { benchmarkGet(10000000, b) }

func benchmarkDelete(dbSize int, b *testing.B) {
	d := Init()
	var str string
	for i := 0; i < dbSize; i++ {
		str = strconv.Itoa(i)
		d.Set(str, str)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d.Set(str, str)
		d.Delete(str)
	}
}

func BenchmarkDelete0(b *testing.B)       { benchmarkDelete(0, b) }
func BenchmarkDelete10(b *testing.B)      { benchmarkDelete(10, b) }
func BenchmarkDelete100(b *testing.B)     { benchmarkDelete(100, b) }
func BenchmarkDelete1000(b *testing.B)    { benchmarkDelete(1000, b) }
func BenchmarkDelete10000(b *testing.B)   { benchmarkDelete(10000, b) }
func BenchmarkDelete100000(b *testing.B)  { benchmarkDelete(100000, b) }
func BenchmarkDelete1000000(b *testing.B) { benchmarkDelete(1000000, b) }

// func BenchmarkDelete10000000(b *testing.B) { benchmarkDelete(10000000, b) }

func benchmarkCount(dbSize int, b *testing.B) {
	d := Init()
	var str string
	for i := 0; i < dbSize; i++ {
		str = strconv.Itoa(i)
		d.Set(str, str)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d.Count(str)
	}
}

func BenchmarkCount0(b *testing.B)       { benchmarkCount(0, b) }
func BenchmarkCount10(b *testing.B)      { benchmarkCount(10, b) }
func BenchmarkCount100(b *testing.B)     { benchmarkCount(100, b) }
func BenchmarkCount1000(b *testing.B)    { benchmarkCount(1000, b) }
func BenchmarkCount10000(b *testing.B)   { benchmarkCount(10000, b) }
func BenchmarkCount100000(b *testing.B)  { benchmarkCount(100000, b) }
func BenchmarkCount1000000(b *testing.B) { benchmarkCount(1000000, b) }

// func BenchmarkGet10000000(b *testing.B) { benchmarkGet(10000000, b) }
