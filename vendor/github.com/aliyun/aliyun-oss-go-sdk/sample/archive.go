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
package sample

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// ArchiveSample archives sample
func ArchiveSample() {
	// Create archive bucket
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	err = client.CreateBucket(bucketName, oss.StorageClass(oss.StorageArchive))
	if err != nil {
		HandleError(err)
	}

	archiveBucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Put archive object
	var val = "花间一壶酒，独酌无相亲。 举杯邀明月，对影成三人。"
	err = archiveBucket.PutObject(objectKey, strings.NewReader(val))
	if err != nil {
		HandleError(err)
	}

	// Check whether the object is archive class
	meta, err := archiveBucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		HandleError(err)
	}

	if meta.Get("X-Oss-Storage-Class") == string(oss.StorageArchive) {
		// Restore object
		err = archiveBucket.RestoreObject(objectKey)
		if err != nil {
			HandleError(err)
		}

		// Wait for restore completed
		meta, err = archiveBucket.GetObjectDetailedMeta(objectKey)
		for meta.Get("X-Oss-Restore") == "ongoing-request=\"true\"" {
			fmt.Println("x-oss-restore:" + meta.Get("X-Oss-Restore"))
			time.Sleep(1000 * time.Second)
			meta, err = archiveBucket.GetObjectDetailedMeta(objectKey)
		}
	}

	// Get restored object
	err = archiveBucket.GetObjectToFile(objectKey, localFile)
	if err != nil {
		HandleError(err)
	}

	// Restore repeatedly
	err = archiveBucket.RestoreObject(objectKey)

	// Delete object and bucket
	err = DeleteTestBucketAndObject(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("ArchiveSample completed")
}
