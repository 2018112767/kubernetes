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
package criu

//Notify interface
type Notify interface {
	PreDump() error
	PostDump() error
	PreRestore() error
	PostRestore(pid int32) error
	NetworkLock() error
	NetworkUnlock() error
	SetupNamespaces(pid int32) error
	PostSetupNamespaces() error
	PostResume() error
}

// NoNotify struct
type NoNotify struct {
}

// PreDump NoNotify
func (c NoNotify) PreDump() error {
	return nil
}

// PostDump NoNotify
func (c NoNotify) PostDump() error {
	return nil
}

// PreRestore NoNotify
func (c NoNotify) PreRestore() error {
	return nil
}

// PostRestore NoNotify
func (c NoNotify) PostRestore(pid int32) error {
	return nil
}

// NetworkLock NoNotify
func (c NoNotify) NetworkLock() error {
	return nil
}

// NetworkUnlock NoNotify
func (c NoNotify) NetworkUnlock() error {
	return nil
}

// SetupNamespaces NoNotify
func (c NoNotify) SetupNamespaces(pid int32) error {
	return nil
}

// PostSetupNamespaces NoNotify
func (c NoNotify) PostSetupNamespaces() error {
	return nil
}

// PostResume NoNotify
func (c NoNotify) PostResume() error {
	return nil
}
