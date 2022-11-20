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
Copyright 2014 The Kubernetes Authors.

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

package migrate

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"k8s.io/client-go/dynamic"
)

const (
	PodSucceeded  = "Succeeded"
	PodHalfFailed = "HalfFailed"
	PodFailed     = "Failed"
)

type MigrateOptions struct {
	DynamicClient dynamic.Interface
	Mapper        meta.RESTMapper
	Result        *resource.Result

	FilenameOptions resource.FilenameOptions
	Node            string

	genericclioptions.IOStreams
}

var (
	migrateLong = templates.LongDesc(i18n.T(`
		Migrate a pod to a new node.`))

	migrateExample = templates.Examples(i18n.T(`
		# Migrate a pod using the data in pod.json.
		kubectl migrate -f pod.yaml --node=nodelables`))
)

func NewMigrateOptions(ioStreams genericclioptions.IOStreams) *MigrateOptions {
	return &MigrateOptions{
		IOStreams: ioStreams,
	}
}

func NewCmdMigrate(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := NewMigrateOptions(ioStreams)

	cmd := &cobra.Command{
		Use:                   "migrate -f FILENAME --storage filePath --node nodeLable",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("migrate a pod to a new node."),
		Long:                  migrateLong,
		Example:               migrateExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.ValidateArgs(cmd, args))
			cmdutil.CheckErr(o.RunMigrate(f, cmd, args))
		},
	}

	// bind flag structs
	usage := "to migrate a pod"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	//cmd.MarkFlagRequired("filename")
	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddApplyAnnotationFlags(cmd)
	cmd.Flags().StringVarP(&o.Node, "node", "l", o.Node, "The node to which the pod is to be migrated, uses label for node(e.g. -n NODENAME)")
	return cmd
}

func (o *MigrateOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
	if len(o.FilenameOptions.Filenames) < 1 {
		return cmdutil.UsageErrorf(cmd, "File for podcheckpoint is null!!")
	}
	return nil
}

func (o *MigrateOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	var err error

	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}
	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}
	return nil
}

func (o *MigrateOptions) RunMigrate(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var (
		err  error
		conf map[string]interface{}
		flag bool
	)

	home := os.Getenv("HOME")
	filenames := o.FilenameOptions.Filenames
	logrus.WithFields(logrus.Fields{
		"filenames": filenames,
	}).Info("podcheckpoint's yaml is: ")

	for _, s := range filenames {
		getConf(&conf, s)
	}
	logrus.Info("success analyse podcheckpoint's yaml")

	metadata := conf["metadata"].(map[interface{}]interface{})
	podcheckpointName := metadata["name"].(string)

	spec := conf["spec"].(map[interface{}]interface{})
	podName := spec["podName"].(string) //带迁移pod名称

	// 创建podcheckpoint CR
	logrus.Info("invoke  o.RunCreateCheckpoint")
	err = o.RunCreateCheckpoint(f, cmd)
	if err != nil {
		logrus.Info("RunCreateCheckpoint Error: ", err.Error())
		return err
	}

	filename := home + "/" + "podcheckpoint-" + podcheckpointName + ".yaml"
	logrus.Info("full filename is ", filename)

	var iter int
	iter = 0
	for {
		status, err := getPodCheckpointStatus(filename, podcheckpointName)
		iter++
		if err == nil {
			flag = false
			switch status {
			case PodSucceeded:
				logrus.Info("Create Checkpoint Succeeded!!!")
				flag = true
				break
			case PodFailed:
				logrus.Info("Create Checkpoint Failed!!!")
				return nil
			case PodHalfFailed:
				logrus.Info("Create Checkpoint HalfFailed!!!")
				return nil
			default: //podcheckpoint处于中间状态，sleep等待迁移完成
				if iter == 1 {
					logrus.Info("Waitting for Create Checkpoint....")
				}
				flag = false
				time.Sleep(time.Millisecond * 2)
			}
			if flag == true {
				break
			}
		} else {
			logrus.Info("getPodCheckpointStatus Failed!!!")
		}
	}

	// 成功创建podcheckpoint， 在目标节点创建pod

	logrus.Info("invoke o.RunCreateNewPod")
	err = clear(filename)
	err = o.RunCreateNewPod(home, podcheckpointName, podName)
	if err != nil {
		fmt.Println("o.RunCreateNewPod ERROR: ", err.Error())
		return err
	}
	return nil
}

// RunCreateCheckpoint -- 使用kubectl create -f pod.yaml 命令创建podcheckpoint
func (o *MigrateOptions) RunCreateCheckpoint(f cmdutil.Factory, cmd *cobra.Command) error {
	createCmd := "kubectl create -f " + o.FilenameOptions.Filenames[0]
	logrus.Info("createCmd = ", createCmd)

	c := exec.Command("bash", "-c", createCmd)
	if err := c.Run(); err != nil {
		logrus.Info("RunCreateCheckpoint Error: ", err)
		return err
	}
	return nil
}

func (o *MigrateOptions) RunCreateNewPod(home string, podcheckpointName string, podName string) error {
	var err error
	inFilepath := home + "/" + podName + "-" + "podcheckpoint" + ".yaml"
	outFilepath := home + "/" + podName + "-" + "podcheckpoint-new" + ".yaml"

	//获取原节点pod的相关信息
	getConfCmd := "kubectl get pod " + podName + " -o yaml > " + inFilepath

	logrus.Info("getConfCmd is ", getConfCmd)
	c := exec.Command("bash", "-c", getConfCmd)
	if err = c.Run(); err != nil {
		return err
	}

	//删除原节点pod
	deleteCmd := "kubectl delete pod " + podName

	logrus.Debug("deleteCmd is ", deleteCmd)
	c = exec.Command("bash", "-c", deleteCmd)
	if err = c.Run(); err != nil {
		return err
	}

	//根据原pod的信息，生成新pod的yaml文件
	HandleYamlFile(inFilepath, outFilepath, o.Node, podcheckpointName)

	//创建新pod
	createCmd := "kubectl create -f " + outFilepath
	logrus.Debug("createCmd is ", createCmd)
	c = exec.Command("bash", "-c", createCmd)
	if err = c.Run(); err != nil {
		return err
	}

	err = clear(inFilepath)
	err = clear(outFilepath)
	return nil
}

// 解析yaml文件  获得podcheckpoint配置
func getConf(conf *map[string]interface{}, filepath string) *map[string]interface{} {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		logrus.Info(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		logrus.Info(err.Error())
	}
	return conf
}

func HandleYamlFile(inFile string, outFile string, nodeName string, podcheckpointName string) error {
	// Read buffer from jsonFile
	byteValue, err := ioutil.ReadFile(inFile)
	if err != nil {
		return err
	}
	// We have known the outer json object is a map, so we define  result as map.
	// otherwise, result could be defined as slice if outer is an array
	var result map[string]interface{}
	err = yaml.Unmarshal(byteValue, &result)
	if err != nil {
		return err
	}
	// handle peers
	spec := result["spec"].(map[interface{}]interface{})

	// 指定要调度的目标节点
	spec["nodeName"] = nodeName

	//添加一条annotation
	// podCheckpoint： podcheckpointName
	metadata := result["metadata"].(map[interface{}]interface{})
	annots := metadata["annotations"]
	if annots == nil {
		var m map[string]string
		m = make(map[string]string)
		m["podCheckpoint"] = podcheckpointName
		metadata["annotations"] = m
	} else {
		annot := annots.(map[interface{}]interface{})
		annot["podCheckpoint"] = podcheckpointName
	}

	// Convert golang object back to byte
	byteValue, err = yaml.Marshal(result)
	if err != nil {
		return err
	}
	// Write back to file
	err = ioutil.WriteFile(outFile, byteValue, 0644)
	if err != nil {
		return err
	}
	return err
}

func getPodCheckpointStatus(filename string, podcheckpointName string) (string, error) {
	logrus.Info("invoke getPodCheckpointStatus")

	getConfCmd := "kubectl get podcheckpoint " + podcheckpointName + " -o yaml > " + filename

	logrus.Info("getConfCmd is ", getConfCmd)

	c := exec.Command("bash", "-c", getConfCmd)
	if err := c.Run(); err != nil {
		return "", err
	}

	var conf map[string]interface{}
	getConf(&conf, filename)
	status := conf["status"].(map[interface{}]interface{})
	
	if status == nil {
		return "", errors.New("podcheckpoint has not status")
	}
	podcheckpointStatus := status["phase"].(string)

	return podcheckpointStatus, nil
}

func clear(filepath string) error {
	c := exec.Command("rm", "-rf", filepath)
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
