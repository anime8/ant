package controllers

import (
	"encoding/json"
)

// 定义返回结果
type ResResult struct {
	Status string      // Success Failure NoLogin  用户没有登录则跳转到登录页面
  Data string
}
// ResResult的interface
type ResResultInterface interface {
	ConvertToString() string          // 将ResResult转换成字符串
}
// 将ResResult转换成字符串
func (res ResResult) ConvertToString() string {
	data,_ := json.Marshal(res)
	return string(data)
}


// TaskResult represents an actual return value of a processed task
type TaskResult struct {
  Type  string      `bson:"type"`
  Value interface{} `bson:"value"`
}


// TaskState represents a state of a task
type TaskState struct {
  TaskUUID  string        `bson:"_id"`
  State     string        `bson:"state"`
  Results   []*TaskResult `bson:"results"`
  Error     string        `bson:"error"`
}
// TaskState的interface
type TaskStateInterface interface {
	GetDeployResult() ResResult           // 获取任务的安装结果
}
// 获取任务的安装结果
func (ts TaskState) GetDeployResult() ResResult {
	var res ResResult
	taskState := ts.Results[0].Value.(string)
	json.Unmarshal([]byte(taskState), &res)
	return res
}
