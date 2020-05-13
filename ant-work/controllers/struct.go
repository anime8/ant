package controllers

type ResResult struct {
	Status string               // Success Failure
  Data string
}

type Redis struct {
	Id                              int
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNodeRedisChecked01       bool
	ClusterNodeSentinelChecked01    bool
	ClusterNode02                   string
	ClusterNodeRedisChecked02       bool
	ClusterNodeSentinelChecked02    bool
	ClusterNode03                   string
	ClusterNodeRedisChecked03       bool
	ClusterNodeSentinelChecked03    bool
	RedisVersion                    string
	SentinelName                    string
	RedisData                       string
	SentinelData                    string
	RedisLog                        string
	RedisConf                       string
	RedisAuthentication             bool
	RedisPassword                   string
}


type Zookeeper struct {
	Id                              int
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNode02                   string
	ClusterNode03                   string
	ZookeeperVersion                string
	DeployPath                      string
	ZookeeperData                   string
	ZookeeperLog                    string
	TaskId                          string
	DeployStatus                    string
	DeployResult                    string
}

type Kafka struct {
	Id                              int
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNode02                   string
	ClusterNode03                   string
	KafkaVersion                    string
	KafkaPath                       string
	KafkaData                       string
	KafkaZookeeper                  string
	TaskId                          string
	DeployStatus                    string
	DeployResult                    string
}
