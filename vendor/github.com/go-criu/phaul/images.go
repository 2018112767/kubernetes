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
	"fmt"
	"os"
	"path/filepath"
)

type images struct {
	cursor int
	dir    string
}

func preparePhaulImages(wdir string) (*images, error) {
	return &images{dir: wdir}, nil
}

func (i *images) getPath(idx int) string {
	return fmt.Sprintf(i.dir+"/%d", idx)
}

func (i *images) openNextDir() (*os.File, error) {
	ipath := i.getPath(i.cursor)
	err := os.Mkdir(ipath, 0700)
	if err != nil {
		return nil, err
	}

	i.cursor++
	return os.Open(ipath)
}

func (i *images) lastImagesDir() string {
	var ret string
	if i.cursor == 0 {
		ret = ""
	} else {
		ret, _ = filepath.Abs(i.getPath(i.cursor - 1))
	}
	return ret
}
