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
// +build go1.7

package oss

import (
	"net"
	"net/http"
)

func newTransport(conn *Conn, config *Config) *http.Transport {
	httpTimeOut := conn.config.HTTPTimeout
	// New Transport
	transport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, httpTimeOut.ConnectTimeout)
			if err != nil {
				return nil, err
			}
			return newTimeoutConn(conn, httpTimeOut.ReadWriteTimeout, httpTimeOut.LongTimeout), nil
		},
		IdleConnTimeout:       httpTimeOut.IdleConnTimeout,
		ResponseHeaderTimeout: httpTimeOut.HeaderTimeout,
	}
	return transport
}
