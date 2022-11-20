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
package phaul

import (
	"os"

	"github.com/checkpoint-restore/go-criu/stats"
	"github.com/golang/protobuf/proto"
)

/* FIXME: report stats from CriuResp */
func criuGetDumpStats(imgDir *os.File) (*stats.DumpStatsEntry, error) {
	stf, err := os.Open(imgDir.Name() + "/stats-dump")
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
	// Skip 2 magic values and entry size
	err = proto.Unmarshal(buf[12:sz], st)
	if err != nil {
		return nil, err
	}

	return st.GetDump(), nil
}
