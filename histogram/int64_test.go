package histogram

import (
	"math/rand"
	"testing"
)

func TestInt64Histogram_WithZero(t *testing.T) {
	h, err := NewSizeHistograms([]Int64HistogramBar{
		{Size: 0, Weight: 1}, {Size: 1, Weight: 1}, {Size: 2, Weight: 1},
	})
	if err != nil {
		t.Fatal("unexpected error at init: ", err)
	}

	t.Log(h.(*int64Histogram).boundary())

	r := rand.New(rand.NewSource(0x1bad1dea))

	for i := 0; i < 100000; i++ {
		v := h.Poll(r)
		if v < 0 || v > 2 {
			t.Fatal("invalid value:", v)
		}
	}
}

func TestInt64Histogram_OnlyZero(t *testing.T) {
	h, err := NewSizeHistograms([]Int64HistogramBar{
		{Size: 0, Weight: 1}, {Size: 0, Weight: 1}, {Size: 0, Weight: 1},
	})
	if err != nil {
		t.Fatal("unexpected error at init: ", err)
	}

	t.Log(h.(*int64Histogram).boundary())

	r := rand.New(rand.NewSource(0x1bad1dea))

	for i := 0; i < 100000; i++ {
		v := h.Poll(r)
		if v != 0 {
			t.Fatal("invalid value:", v)
		}
	}
}

func TestInt64Histogram_Empty(t *testing.T) {
	h, err := NewSizeHistograms([]Int64HistogramBar{})
	if err == nil {
		if h != nil {
			t.Fatal("unexpected success at init (and non-nil histogram returned): ", err)

		} else {
			t.Fatal("unexpected success at init: ", err)
		}
	}
}
