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
package oss

import (
	"math"
	"net/http"

	. "gopkg.in/check.v1"
)

type OssErrorSuite struct{}

var _ = Suite(&OssErrorSuite{})

func (s *OssErrorSuite) TestCheckCRCHasCRCInResp(c *C) {
	headers := http.Header{
		"Expires":              {"-1"},
		"Content-Length":       {"0"},
		"Content-Encoding":     {"gzip"},
		"X-Oss-Hash-Crc64ecma": {"0"},
	}

	resp := &Response{
		StatusCode: 200,
		Headers:    headers,
		Body:       nil,
		ClientCRC:  math.MaxUint64,
		ServerCRC:  math.MaxUint64,
	}

	err := checkCRC(resp, "test")
	c.Assert(err, IsNil)
}

func (s *OssErrorSuite) TestCheckCRCNotHasCRCInResp(c *C) {
	headers := http.Header{
		"Expires":          {"-1"},
		"Content-Length":   {"0"},
		"Content-Encoding": {"gzip"},
	}

	resp := &Response{
		StatusCode: 200,
		Headers:    headers,
		Body:       nil,
		ClientCRC:  math.MaxUint64,
		ServerCRC:  math.MaxUint32,
	}

	err := checkCRC(resp, "test")
	c.Assert(err, IsNil)
}

func (s *OssErrorSuite) TestCheckCRCCNegative(c *C) {
	headers := http.Header{
		"Expires":              {"-1"},
		"Content-Length":       {"0"},
		"Content-Encoding":     {"gzip"},
		"X-Oss-Hash-Crc64ecma": {"0"},
	}

	resp := &Response{
		StatusCode: 200,
		Headers:    headers,
		Body:       nil,
		ClientCRC:  0,
		ServerCRC:  math.MaxUint64,
	}

	err := checkCRC(resp, "test")
	c.Assert(err, NotNil)
	testLogger.Println("error:", err)
}

func (s *OssErrorSuite) TestCheckDownloadCRC(c *C) {
	err := checkDownloadCRC(0xFBF9D9603A6FA020, 0xFBF9D9603A6FA020)
	c.Assert(err, IsNil)

	err = checkDownloadCRC(0, 0)
	c.Assert(err, IsNil)

	err = checkDownloadCRC(0xDB6EFFF26AA94946, 0)
	c.Assert(err, NotNil)
	testLogger.Println("error:", err)
}
