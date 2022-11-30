/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dumpstats

import (
	"k8s.io/klog"
	"os"
	"strconv"

	"github.com/go-criu/stats"
	"github.com/golang/protobuf/proto"
)

var (
	MAX_ITER_COUNT   = 7
	MIN_PAGE_WRITTEN = 500
	MAX_GROW_RATE    = 10
)

func CheckPreCopy(iterCount int, checkpointDir string, containerName string) (bool, error) {
	klog.Warningln("Check if need precopy!")

	if iterCount <= 1 {
		dumpStatFile := checkpointDir + "/" + containerName + "_" + "Dump" + strconv.Itoa(iterCount) + "/criu.work" + "/stats-dump"
		dumpstat, err := getDumpStats(dumpStatFile)
		if err != nil {
			klog.Errorln("Get DumpStat Failed!!! error is %v", err)
			return false, err
		}
		var minPage uint64 = uint64(MIN_PAGE_WRITTEN)
		if *dumpstat.PagesWritten <= minPage {
			klog.Errorln("The dump pages is small enough to do final dump!!!")
			return false, nil
		} else {
			return true, nil
		}
	} else if iterCount >= MAX_ITER_COUNT {
		klog.Warningln("Reach the max iter count!!!")
		return false, nil
	} else {
		dumpStatFile := checkpointDir + "/" + containerName + "_" + "Dump" + strconv.Itoa(iterCount) + "/criu.work" + "/stats-dump"
		dumpstat, err := getDumpStats(dumpStatFile)
		if err != nil {
			klog.Errorln("Get DumpStat Failed!!! error is %v", err)
			return false, err
		}

		predumpStatFile := checkpointDir + "/" + containerName + "_" + "Dump" + strconv.Itoa(iterCount-1) + "/criu.work" + "/stats-dump"
		predumpstat, err := getDumpStats(predumpStatFile)
		if err != nil {
			klog.Errorln("Get PreDumpStat Failed!!! error is %v", err)
			return false, err
		}
		growRate := calGrowRate(*dumpstat.PagesWritten, *predumpstat.PagesWritten)
		var Rate uint64 = uint64(MAX_GROW_RATE)
		if growRate >= Rate {
			klog.Errorln("Written pages grows too fast with iteration!!!!")
			return false, nil
		}
	}
	return true, nil
}

func getDumpStats(filePath string) (*stats.DumpStatsEntry, error) {
	stf, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer stf.Close()

	buf := make([]byte, 2*4096)
	sz, err := stf.Read(buf)
	if err != nil {
		return nil, err
	}

	st := &stats.StatsEntry{}
	proto.Unmarshal(buf[12:sz], st)

	return st.GetDump(), nil
}

func calGrowRate(pagesWritten uint64, prepagesWritten uint64) uint64 {
	diff := pagesWritten - prepagesWritten
	growRate := diff * 100 / prepagesWritten
	return growRate
}
