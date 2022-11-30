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
	"k8s.io/klog"
	"os/exec"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"k8s.io/api/core/v1"
)

// const DefaultCheckpointDir = "/var/lib/docker/containers"
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

func ProcessStorage(storage string, containerID string, checkpointID string, op string, secret *v1.Secret) (string, error) {
	klog.Infoln("Start Process Storage! storage = ", storage)
	var (
		checkpointDir string
		err           error
	)

	start := strings.Index(storage, "://")
	if start == 0 {
		klog.Infoln("Error: Storage format is invalid")
		return "", nil
	}

	prefix := storage[0:start]
	path := storage[start+3:]

	// /var/lib/docker/containers/{$dockerid}/checkpoints/{$checkpoint_name}
	if prefix == "oss" {
		//checkpointDir = DefaultCheckpointDir + "/" + containerID + "/checkpoints"
		checkpointDir = DefaultCheckpointDir
		ossOption := getOssOptionFromStorage(path, secret)
		//ossOption.ObjectName = checkpointID + "/" + checkpointName
		ossOption.ObjectName = checkpointID

		fmt.Printf("ossOption: %v", ossOption)
		if op == DOWNLOAD {
			klog.Infoln("invoke DOWNLOAD")
			downloadPath := checkpointDir + "/" + checkpointID + ".tar.gz"
			klog.Infoln("downloadPath = ", downloadPath)
			if err = downloadCheckpoint(ossOption, downloadPath); err != nil {
				handleError(err)
				return "", err
			}
			c := exec.Command("tar", "-xvjf", downloadPath, "-C", "/")
			if err = c.Run(); err != nil {
				klog.Infoln("Error: ", err)
			}
		} else if op == UPLOAD {
			klog.Infoln("invoke UPLOAD")

			checkpointFilePath := checkpointDir + "/" + checkpointID
			uploadPath := checkpointDir + "/" + ossOption.ObjectName + ".tar.gz"

			klog.Infoln("uploadPath = ", uploadPath)
			klog.Infoln("checkpointFilePath = ", checkpointFilePath)
			c := exec.Command("tar", "-cvjf", uploadPath, checkpointFilePath)
			if err = c.Run(); err != nil {
				klog.Infoln("Error: ", err)
			}
			if err = uploadCheckpoint(ossOption, uploadPath); err != nil {
				handleError(err)
				return "", err
			}
		} else if op == CHECKPOINT {
			klog.Infoln("invoke CHECKPOINT")
			return checkpointDir, nil
		} else if op == CLEAR {
			klog.Infoln("invoke CLEAR")
			klog.Infoln("entering clear, remove checkpoint files, file path is: ", checkpointDir+"/"+checkpointID+"*")
			Cmd := "rm -rf " + checkpointDir + "/" + checkpointID + "*"
			c := exec.Command("bash", "-c", Cmd)
			if err = c.Run(); err != nil {
				klog.Errorln("Error: ", err)
			}

			Cmd = "rm -rf " + "/var/lib/docker/containers/" + containerID + "/checkpoints/" + checkpointID
			c = exec.Command("bash", "-c", Cmd)
			if err = c.Run(); err != nil {
				klog.Errorln("Error: ", err)
			}
		}
	} else if prefix == "file" {
		checkpointDir = path
		klog.Infoln("checkpointDir is %s\n", path)
	}
	return checkpointDir, nil
}

func handleError(err error) {
	klog.Infoln("Error:", err)
}

func downloadCheckpoint(option *OSSOption, path string) error {
	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(option.Bucket)
	if err != nil {
		return err
	}

	err = bucket.GetObjectToFile(option.ObjectName, path)
	return err
}

func uploadCheckpoint(option *OSSOption, path string) error {
	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(option.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(option.ObjectName, path)
	if err != nil {
		return err
	}

	return nil
}

func MoveCheckpointFileToContainerDirectory(srcPath string, dstPath string) error {
	Cmd := "mv " + srcPath + " " + dstPath
	c := exec.Command("bash", "-c", Cmd)
	err := c.Run()
	if err != nil {
		klog.Infoln("Error: ", err)
	}
	return err
}

func getOssOptionFromStorage(storage string, secret *v1.Secret) *OSSOption {
	ossOption := &OSSOption{}
	var pre string
	var back string
	pre, back = getNextSplitStr(storage, "/")
	//fmt.Printf("secret: %v", secret)

	if back == "" {
		klog.Infoln("PodCheckpoint.Spec.Storage is invalid")
		return nil
	}

	for k, v := range secret.Data {
		klog.Infoln("k = %s, v = %s", k, v)
	}

	ossOption.Endpoint = pre
	//ossOption.AccessKeyId = "LTAI5tDncGtnwEdFNYn9MUVC"
	//ossOption.AccessKeySecret = "ge7ktyzN3FLPKfvBVpu6Uk202L43Qe"

	content, ok := secret.Data["accessKeyId"]
	if !ok {
		klog.Errorln("can not read accessKeyId")
		return nil
	}
	ossOption.AccessKeyId = string(content)
	content, ok = secret.Data["accessKeySecret"]
	if !ok {
		klog.Errorln("can not read accessKeySecret")
		return nil
	}
	ossOption.AccessKeySecret = string(content)
	//ossOption.AccessKeyId = string(secret.Data["accessKeyId"])
	//ossOption.AccessKeySecret = string(secret.Data["accessKeySecret"])

	ossOption.Bucket = back
	return ossOption
}

func getNextSplitStr(s string, div string) (string, string) {
	i := strings.Index(s, div)
	pre := s[0:i]
	back := s[i+len(div):]
	return pre, back
}
