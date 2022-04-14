package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test cases include some initial sanity checks, followed by the ones from the spec
// I was going to go into some benchmark tests to see if I could prove the O log n
// but ran out of time

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
