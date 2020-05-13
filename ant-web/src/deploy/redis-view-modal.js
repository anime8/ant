import React from 'react'
import { Progress, Table, Button, Header, Image, Modal, Divider, Form, Grid, Radio } from 'semantic-ui-react'
// import axios from 'axios';
import Backend from '../conf';
import _ from 'underscore';
import $ from  'jquery';
import Modalinfo from '../home/modal-info';
import { CleanToken } from '../home/cookies';
import FormFieldMessage from '../common/message';

class ViewRedisModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      modalOpen: false,
      percent: 0
    };
  }

  // 卸载的时候关闭定时获取
  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  // 设置progress进度
  handleProgress = (DeployStatus) => {
    const PercentList = {
      PENDING: 0,
      RECEIVED: 10,
      STARTED: 50,
      RETRY: 0,
      FAILURE: 60,
      SUCCESS: 100
    };
    this.setState({ percent: PercentList[DeployStatus] });
  }


  // 获取redis安装进度
  async getRedisInfo() {
    // 重新给react的this换个名字，防止和post返回的this冲突
    var reactThis = this;
    // 当this.state.redis取到值后，函数尝试获取数据
    if (this.state.redis) {
      // 如果状态DeployStatus数组中的状态，则去后端定时获取安装进度
      var DeployStatus = ["PENDING", "RECEIVED", "STARTED", "RETRY"];
      var index = _.indexOf(DeployStatus, this.state.redis.DeployStatus);
      if (index !== -1) {
        var url = Backend + "api/redis/getone/"
        $.post(url,
            JSON.stringify({
              Id: this.state.redis.Id,
            }),
            function(data,status){
              data = JSON.parse(data);
              if (data.Status === "Success") {
                // 刷新改组件redis信息
                var redis = JSON.parse(data.Data);
                reactThis.setState({ redis: redis });
                reactThis.handleProgress(redis.DeployStatus);
              }else if (data.Status === "NoLogin") {
                CleanToken();
              }else {
                Modalinfo.fire(<p>{data.Data}</p>);
              }
            });
      }else{
        // 如果安装状态为SUCCESS，则将progress设置为100
        if (this.state.redis.DeployStatus === "SUCCESS") {
          this.setState({ percent: 100 });
        }
      }
    }
  }

  // 打开modal
  handleOpen = () => {
    this.setState({ modalOpen: true });
    // 将传递到该元素的props中的redis信息赋值到state
    let redis = JSON.parse(this.props.redis);
    this.setState({ redis: redis });
    // 如果初始DeployStatus为SUCCESS，则进度条为100%
    if (redis.DeployStatus === "SUCCESS") {
      this.setState({ percent: 100 });
    }
    // 设置定时获取redis最新信息
    this.timerID = setInterval(
      () => this.getRedisInfo(),
      1000
    );
  }

  // 关闭modal，同时关闭定时获取数据
  handleClose = () => {
    this.setState({ modalOpen: false });
    clearInterval(this.timerID);

    // 如果redis列表中该条redis信息的部署状态不为FAILURE或者是SUCCESS，则刷新redis列表信息
    var DeployStatus = ["FAILURE", "SUCCESS"];
    var index = _.indexOf(DeployStatus, this.props.DeployStatus);
    if (index === -1) {
      this.props.refreshRedisList();
    }
  }


  render() {
    const redis = this.state.redis;
    const baseInfo = [
      {header: "集群名称", name: "ClusterName"},
      {header: "备注", name: "Remark"},
    ];
    const confInfo = [
      {header: "redis版本", name: "RedisVersion"},
      {header: "哨兵名称", name: "SentinelName"},
      {header: "redis数据目录", name: "RedisData"},
      {header: "sentinel数据目录", name: "SentinelData"},
      {header: "redis日志目录", name: "RedisLog"},
      {header: "redis配置目录", name: "RedisConf"},
    ];
    const CheckInputList = [
      {name: "ClusterNode01", redis: "ClusterNodeRedisChecked01", sentinel: "ClusterNodeSentinelChecked01"},
      {name: "ClusterNode02", redis: "ClusterNodeRedisChecked02", sentinel: "ClusterNodeSentinelChecked02"},
      {name: "ClusterNode03", redis: "ClusterNodeRedisChecked03", sentinel: "ClusterNodeSentinelChecked03"},
    ];
    return (
      <Modal size="large" trigger={
        <Button
          circular
          icon='eye'
          onClick={this.handleOpen}
        >
        </Button>
      }
      open={this.state.modalOpen}
      onClose={this.handleClose}
      >
        <Modal.Header>
          <Image wrapped src='/redis-white.png' size='mini' />
          Rdis集群信息
        </Modal.Header>
        <Modal.Content>
          <Modal.Description>
          <React.Fragment>
          <Header as='h3'>基础信息</Header>
          <Form>
            {
              redis && baseInfo.map(
                (base) => <FormFieldMessage key={base.name} header={base.header} content={redis[base.name]} />
              )
            }
          </Form>

          <Divider hidden />

          <Header as='h3'>配置信息</Header>
          <Form>
            {
              redis && confInfo.map(
                (conf) => <FormFieldMessage key={conf.name} header={conf.header} content={redis[conf.name]} />
              )
            }
            <Form.Field>
              <label>节点信息</label>
                <Grid>
                  <Grid.Column width={16}>
                    <Table>
                      <Table.Body>
                        {
                          CheckInputList.map(
                            (CheckInput) =>
                            <Table.Row key={CheckInput.name}>
                              <Table.Cell collapsing>
                                <p>{redis ? redis[CheckInput.name] : ''}</p>
                              </Table.Cell>
                              <Table.Cell>
                                <Radio
                                  checked={redis ? redis[CheckInput.redis] : ''}
                                  label='redis'
                                />
                              </Table.Cell>
                              <Table.Cell>
                                <Radio
                                  checked={redis ? redis[CheckInput.sentinel] : ''}
                                  label='sentinel'
                                />
                              </Table.Cell>
                            </Table.Row>
                          )
                        }
                      </Table.Body>
                    </Table>
                  </Grid.Column>
                </Grid>
            </Form.Field>
            <Form.Field>
              <label>是否开启验证</label>
              <Grid>
                <Grid.Column width={16}>
                  <Table>
                    <Table.Body>
                      <Table.Row>
                        <Table.Cell collapsing>
                          <Radio
                            checked={redis ? redis.RedisAuthentication : ''}
                            label='密码验证'
                          />
                        </Table.Cell>
                        <Table.Cell>
                          <p>{redis ? redis.RedisPassword : ''}</p>
                        </Table.Cell>
                      </Table.Row>
                    </Table.Body>
                  </Table>
                </Grid.Column>
              </Grid>
            </Form.Field>
          </Form>
          <Divider hidden />

          <Header as='h3'>安装进度</Header>
          <div>
            <Progress percent={this.state.percent} indicating progress label={redis ? redis.DeployStatus : ''} />
          </div>

          </React.Fragment>
          </Modal.Description>
        </Modal.Content>
        <Modal.Actions>
          <Button negative onClick={this.handleClose}>取消</Button>
        </Modal.Actions>
      </Modal>

    );
  }
}

export default ViewRedisModal
