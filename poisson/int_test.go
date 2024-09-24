// Copyright (c) 2018-2024 Jean-Francois SMIGIELSKI
// Copyright (c) 2024 OVHCloud SAS
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for
// the specific language governing permissions and limitations under the License.

package poisson

import "testing"

func TestPoissonSpeed(t *testing.T) {
	for i := 1; i < 300; i++ {
		t.Logf("lambda=%d", i)
		NewPoissonSlots(i)
	}
}

func TestPoisson(t *testing.T) {
	p := NewPoissonSlots(10)
	for i, slot := range p.probabilities {
		t.Logf("slot %d: %+v", i, slot)
	}
	for i := 0; i < 50; i++ {
		t.Logf("iter %v value %v", i, p.Poll())
	}
}
