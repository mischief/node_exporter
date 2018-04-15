// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !noboottime

package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

/*
#include <sys/time.h>
#include <time.h>
*/
import "C"

type bootTimeCollector struct{}

func init() {
	registerCollector("boottime", defaultEnabled, NewBootTimeCollector)
}

// NewTimeCollector returns a new Collector exposing the current system boot time in
// seconds since epoch.
func NewBootTimeCollector() (Collector, error) {
	return &bootTimeCollector{}, nil
}

func (c *bootTimeCollector) Update(ch chan<- prometheus.Metric) error {
	var cboottime C.struct_timespec

	_, err := C.clock_gettime(C.CLOCK_BOOTTIME, &cboottime)
	if err != nil {
		return err
	}

	bts := float64(time.Now().Unix() - int64(cboottime.tv_sec))

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "boot_time_seconds"),
			"Node boot time, in unixtime.",
			nil, nil,
		),
		prometheus.GaugeValue, bts,
	)
	return nil
}
