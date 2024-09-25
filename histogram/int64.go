// Copyright (c) 2018-2024 Jean-Francois SMIGIELSKI
// Copyright (c) 2024 OVHCloud SAS
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for
// the specific language governing permissions and limitations under the License.
package histogram

import (
	"errors"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

var ErrBadCSV = errors.New("bad CSV line")
var ErrBadWeight = errors.New("negative weight configured")
var ErrBadSize = errors.New("negative size configured")
var ErrEmpty = errors.New("no size configured")

type Int64HistogramBar struct {
	Size   int64 `yaml:"size"`
	Weight int64 `yaml:"weight"`
}

type Int64Distribution interface {
	// Produces a new int64 respecting the distribution and based on the given uniform PRNG
	Poll(r *rand.Rand) int64
}

type histogramBars []Int64HistogramBar

type int64Histogram struct {
	bars histogramBars
}

// Implements sort.Interface
func (s histogramBars) Len() int {
	return len(s)
}

// Implements sort.Interface
func (s histogramBars) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Implements sort.Interface
func (s histogramBars) Less(i, j int) bool {
	return s[i].Size < s[j].Size
}

// NewSizeHistograms returns a Int64Distribution implementing histogram-like buckets, ready to use and based on a
// collection of size bars
func NewSizeHistograms(sizes []Int64HistogramBar) (Int64Distribution, error) {
	sizeHistograms := &int64Histogram{
		bars: make(histogramBars, 0),
	}
	return sizeHistograms.init(sizes)
}

// ParseCSV returns a Int64Distribution implementing histogram-like buckets, ready to use and based on a collection
// of size bars described as a coma-separated sequence of "size:weight" values
func ParseCSV(csv string) (Int64Distribution, error) {
	tokens := strings.Split(csv, ",")
	if len(tokens) < 1 {
		return nil, ErrBadCSV
	}
	return ParseTokens(tokens, ":")
}

// ParseTokens returns a Int64Histogram ready to use, based on a collection of size bars stored in an array of
// "size <separator> weight" tokens
func ParseTokens(pairs []string, separator string) (Int64Distribution, error) {
	histograms := make([]Int64HistogramBar, 0)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.Split(pair, separator)
		bar := Int64HistogramBar{}
		var err error
		if bar.Size, err = strconv.ParseInt(kv[0], 10, 64); err != nil {
			return nil, err
		}
		if bar.Weight, err = strconv.ParseInt(kv[1], 10, 64); err != nil {
			return nil, err
		}
		histograms = append(histograms, bar)
	}
	return NewSizeHistograms(histograms)
}

// `sizes` must be non-empty and doesn't need to be sorted
func (s *int64Histogram) init(sizes []Int64HistogramBar) (*int64Histogram, error) {
	if len(sizes) == 0 {
		return nil, ErrEmpty
	}
	for _, bar := range sizes {
		if bar.Weight < 0 {
			return nil, ErrBadWeight
		}
		if bar.Size < 0 {
			return nil, ErrBadSize
		}
		s.bars = append(s.bars, bar)
	}
	sort.Sort(s.bars)

	// normalize the bars
	total := int64(0)
	for i, _ := range s.bars {
		total += s.bars[i].Weight
		s.bars[i].Weight = total
	}

	return s, nil
}

func (s *int64Histogram) boundary() int64 {
	if len(s.bars) == 0 {
		panic("argl")
	}
	return s.bars[len(s.bars)-1].Weight
}

func (s *int64Histogram) Poll(r *rand.Rand) int64 {
	needle := r.Int63n(s.boundary())
	for i, x := range s.bars { // TODO(jfs): a binary search would maybe perform better
		if x.Weight > needle { // we have the right slot
			prev := int64(0)
			if i > 0 {
				prev = s.bars[i-1].Size
			}
			if x.Size <= 0 {
				return prev
			}
			return prev + r.Int63n(x.Size-prev)
		}
	}
	panic("plop")

}
