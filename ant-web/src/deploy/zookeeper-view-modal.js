import React from 'react'
import { Progress, Button, Header, Image, Modal, Divider, Form } from 'semantic-ui-react'
// import axios from 'axios';
import Backend from '../conf';
import _ from 'underscore';
import $ from  'jquery';
import Modalinfo from '../home/modal-info';
import { CleanToken } from '../home/cookies';
import FormFieldMessage from '../common/message';

class ViewZookeeperModal extends React.Component {

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


  // 获取zookeeper安装进度
  async getZookeeperInfo() {
    // 重新给react的this换个名字，防止和post返回的this冲突
    var reactThis = this;
    // 当this.state.redis取到值后，函数尝试获取数据
    if (this.state.zookeeper) {
      // 如果状态DeployStatus数组中的状态，则去后端定时获取安装进度
      var DeployStatus = ["PENDING", "RECEIVED", "STARTED", "RETRY"];
      var index = _.indexOf(DeployStatus, this.state.zookeeper.DeployStatus);
      if (index !== -1) {
        var url = Backend + "api/zookeeper/getone/"
        $.post(url,
            JSON.stringify({
              Id: this.state.zookeeper.Id,
            }),
            function(data,status){
              data = JSON.parse(data);
              if (data.Status === "Success") {
                // 刷新改组件redis信息
                var zookeeper = JSON.parse(data.Data);
                reactThis.setState({ zookeeper: zookeeper });
                reactThis.handleProgress(zookeeper.DeployStatus);
              }else if (data.Status === "NoLogin") {
                CleanToken();
              }else {
                Modalinfo.fire(<p>{data.Data}</p>);
              }
            });
      }else{
        // 如果安装状态为SUCCESS，则将progress设置为100
        if (this.state.zookeeper.DeployStatus === "SUCCESS") {
          this.setState({ percent: 100 });
        }
      }
    }
  }

  // 打开modal
  handleOpen = () => {
    this.setState({ modalOpen: true });
    // 将传递到该元素的props中的redis信息赋值到state
    let zookeeper = JSON.parse(this.props.zookeeper);
    this.setState({ zookeeper: zookeeper });
    // 如果初始DeployStatus为SUCCESS，则进度条为100%
    if (zookeeper.DeployStatus === "SUCCESS") {
      this.setState({ percent: 100 });
    }
    // 设置定时获取redis最新信息
    this.timerID = setInterval(
      () => this.getZookeeperInfo(),
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
      this.props.refreshZookeeperList();
    }
  }


  render() {
    const zookeeper = this.state.zookeeper;
    const baseInfo = [
      {header: "集群名称", name: "ClusterName"},
      {header: "备注", name: "Remark"},
    ];
    const confInfo = [
      {header: "zookeeper安装目录", name: "DeployPath"},
      {header: "zookeeper数据目录", name: "ZookeeperData"},
      {headerheader: "zookeeper日志目录", name: "ZookeeperLog"},
      {header: "zookeeper节点1", name: "ClusterNode01"},
      {header: "zookeeper节点2", name: "ClusterNode02"},
      {header: "zookeeper节点3", name: "ClusterNode03"},
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
          <Image wrapped src='/zookeeper.svg' />
          Zookeeper集群信息
        </Modal.Header>
        <Modal.Content>
          <Modal.Description>
          <React.Fragment>
          <Header as='h3'>基础信息</Header>
          <Form>
            {
              zookeeper && baseInfo.map(
                (base) => <FormFieldMessage key={base.name} header={base.header} content={zookeeper[base.name]} />
              )
            }
          </Form>

          <Divider hidden />

          <Header as='h3'>配置信息</Header>
          <Form>
            {
              zookeeper && confInfo.map(
                (conf) => <FormFieldMessage key={conf.name} header={conf.header} content={zookeeper[conf.name]} />
              )
            }
          </Form>
          <Divider hidden />

          <Header as='h3'>安装进度</Header>
          <div>
            <Progress percent={this.state.percent} indicating progress label={zookeeper ? zookeeper.DeployStatus : ''} />
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

export default ViewZookeeperModal
