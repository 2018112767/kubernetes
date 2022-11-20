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

package util

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"k8s.io/api/core/v1"
)

const DefaultCheckpointDir = "/tmp"

type OSSOption struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	ObjectName      string
}

const (
	UPLOAD     = "upload"
	DOWNLOAD   = "download"
	CHECKPOINT = "checkpoint"
	CLEAR      = "clear"
)

func ProcessStorage(storage string, checkpointID string, op string, secret *v1.Secret) (string, error) {
	fmt.Println("Start Process Storage! storage = ", storage)
	var (
		checkpointDir string
		err           error
	)

	start := strings.Index(storage, "://")
	if start == 0 {
		fmt.Println("Error: Storage format is invalid")
		return "", nil
	}

	prefix := storage[0:start]
	path := storage[start+3:]
	fmt.Println("prefix = ", prefix)
	fmt.Println("path = ", path)

	if prefix == "oss" {
		checkpointDir = DefaultCheckpointDir
		ossOption := getOssOptionFromStorage(path, secret)
		ossOption.ObjectName = checkpointID
		fmt.Printf("ossOption: %v", ossOption)
		if op == DOWNLOAD {
			fmt.Printf("invoke DOWNLOAD")
			if err = downloadCheckpoint(ossOption); err != nil {
				handleError(err)
				return "", err
			}
			c := exec.Command("tar", "-xvjf", checkpointDir+"/"+ossOption.ObjectName+".tar.gz", "-C", "/")
			if err = c.Run(); err != nil {
				fmt.Println("Error: ", err)
			}
		} else if op == UPLOAD {
			fmt.Printf("invoke UPLOAD")
			c := exec.Command("tar", "-cvjf", checkpointDir+"/"+ossOption.ObjectName+".tar.gz", checkpointDir+"/"+checkpointID)
			if err = c.Run(); err != nil {
				fmt.Println("Error: ", err)
			}
			if err = uploadCheckpoint(ossOption); err != nil {
				handleError(err)
				return "", err
			}
		} else if op == CHECKPOINT {
			fmt.Printf("invoke CHECKPOINT")
			return checkpointDir, nil
		} else if op == CLEAR {
			fmt.Println("invoke CLEAR")
			fmt.Println("entering clear, remove checkpoint files, file path is: ", checkpointDir+"/"+checkpointID+"*")
			Cmd := "rm -rf " + checkpointDir + "/" + checkpointID + "*"
			c := exec.Command("bash", "-c", Cmd)
			if err = c.Run(); err != nil {
				fmt.Println("Error: ", err)
			}
		}
	} else if prefix == "file" {
		checkpointDir = path
		fmt.Println("checkpointDir is %s\n", path)
	}
	return checkpointDir, nil
}

func handleError(err error) {
	fmt.Println("Error:", err)
}

func downloadCheckpoint(option *OSSOption) error {
	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(option.Bucket)
	if err != nil {
		return err
	}
	path := DefaultCheckpointDir + "/" + option.ObjectName + ".tar.gz"
	err = bucket.GetObjectToFile(option.ObjectName, path)
	return err
}

func uploadCheckpoint(option *OSSOption) error {
	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(option.Bucket)
	if err != nil {
		return err
	}
	path := DefaultCheckpointDir + "/" + option.ObjectName + ".tar.gz"
	err = bucket.PutObjectFromFile(option.ObjectName, path)
	if err != nil {
		return err
	}

	return nil
}

func getOssOptionFromStorage(storage string, secret *v1.Secret) *OSSOption {
	ossOption := &OSSOption{}
	var pre string
	var back string
	pre, back = getNextSplitStr(storage, "/")
	fmt.Println("pre = ", pre)
	fmt.Println("back = ", back)
	//fmt.Printf("secret: %v", secret)

	if back == "" {
		fmt.Println("PodCheckpoint.Spec.Storage is invalid")
		return nil
	}
	ossOption.Endpoint = pre
	//ossOption.AccessKeyId = string(secret.Data["accessKeyId"])
	//ossOption.AccessKeySecret = string(secret.Data["accessKeySecret"])
	ossOption.AccessKeyId = "LTAI5tGJ5u6nCqaY4dn98Tnw"
	ossOption.AccessKeySecret = "qhT13h5UbM568FDcIykpsCuNnFrZHr"
	ossOption.Bucket = back
	return ossOption
}

func getNextSplitStr(s string, div string) (string, string) {
	i := strings.Index(s, div)
	pre := s[0:i]
	back := s[i+len(div):]
	return pre, back
}
