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

type Int64HistogramBar struct {
	Size   int64 `yaml:"size"`
	Weight int64 `yaml:"weight"`
}

type Int64Histogram []Int64HistogramBar

// Implements sort.Interface
func (s Int64Histogram) Len() int {
	return len(s)
}

// Implements sort.Interface
func (s Int64Histogram) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Implements sort.Interface
func (s Int64Histogram) Less(i, j int) bool {
	return s[i].Size < s[j].Size
}

// Returns a Int64Histogram ready to use, based on a collection of size slots
func NewSizeHistograms(sizes Int64Histogram) Int64Histogram {
	sizeHistograms := make(Int64Histogram, len(sizes))
	sizeHistograms.Init(sizes)
	return sizeHistograms
}

// Returns a Int64Histogram ready to use, based on a collection of size slots described as a coma-separated sequence
// of "size:weight" values
func ParseCSV(csv string) (Int64Histogram, error) {
	tokens := strings.Split(csv, ",")
	if len(tokens) < 1 {
		return nil, ErrBadCSV
	}
	return ParseTokens(tokens, ":")
}

// Returns a Int64Histogram ready to use, based on a collection of size slots stored in an array of "size <separator> weight"
// tokens
func ParseTokens(pairs []string, separator string) (Int64Histogram, error) {
	histograms := make(Int64Histogram, 0)
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
	return NewSizeHistograms(histograms), nil
}

// `sizes` must be non-empty and doesn't need to be sorted
func (s Int64Histogram) Init(sizes []Int64HistogramBar) {
	copy(s, sizes)
	sort.Sort(s)
	total := int64(0)
	for i, _ := range s {
		total += (s)[i].Weight
		(s)[i].Weight = total
	}
}

func (s Int64Histogram) locate(needle int64) int64 {
	for i, x := range s {
		if x.Weight > needle { // we have the right slot
			prev := int64(0)
			if i > 0 {
				prev = s[i-1].Size
			}
			return prev + rand.Int63n(x.Size-prev)
		}
	}
	panic("plop")
}

func (s Int64Histogram) boundary() int64            { return s[len(s)-1].Weight }
func (s Int64Histogram) PollRand(r rand.Rand) int64 { return s.locate(r.Int63n(s.boundary())) }
func (s Int64Histogram) Poll() int64                { return s.locate(rand.Int63n(s.boundary())) }
