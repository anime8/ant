import React from 'react'
import { Icon, Button, Header, Image, Modal, Divider, Form, TextArea, Dropdown } from 'semantic-ui-react'
import $ from  'jquery'
import Backend from '../conf';
import Modalinfo from '../home/modal-info';
import _ from 'underscore';
import { CleanToken } from '../home/cookies';
import FormFieldInput from '../common/input';

// redis版本列表
const zookeeperOptions = [
  {
    key: 'zookeeper-3.4.14',
    text: 'zookeeper-3.4.14',
    value: 'zookeeper-3.4.14',
    image: { avatar: true, src: '/zookeeper.svg' },
  },
];


// 部署redis弹出框
class DeployZookeeperModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      checkboxLoading: "",
      modalOpen: false,
      zookeeperInput: {
        ClusterName: "",
        Remark: "",
        ClusterNode01: "",
        ClusterNode02: "",
        ClusterNode03: "",
        ZookeeperVersion: "zookeeper-3.4.14",
        DeployPath: "/opt/app",
        ZookeeperData: "/opt/data/zookeeper/data",
        ZookeeperLog: "/opt/data/zookeeper/log",
      }
    };
  }

  // 控制modal打开和关闭
  handleOpen = () => this.setState({ modalOpen: true })
  handleClose = () => {
    this.setState({ modalOpen: false });
    this.setState({ checkboxLoading: "" });
  };

  // 当输入框发生变化时，给state赋值
  handleInputChange = (event) => {
    this.setState({
      zookeeperInput: {
        ...this.state.zookeeperInput,
        [event.target.name]: event.target.value
      }
    });
  }

  // 获取redis版本信息
  handleZookeeperVersionChange = (event) => {
    this.setState({
      zookeeperInput: {
        ...this.state.zookeeperInput,
        ZookeeperVersion: event.target.innerText
      }
    });
  }

  // 点击提交部署任务
  handleInstall = () => {
    // 部署按钮开始loading
    this.setState({ checkboxLoading: "loading disabled" });
    var handleClose = this.handleClose;
    // 设置变量，后面在post里面调用，这是为了防止和post返回的this起冲突
    var refreshZookeeperList = this.props.refreshZookeeperList;
    var url = Backend + "api/zookeeper/install/";
    // 判断是否缺少必要参数，如果缺少则提示用户
    var RequredFileds = [
      "ClusterName",
      "ClusterNode01",
      "ClusterNode02",
      "ClusterNode03",
       "RedisVersion",
       "DeployPath",
       "ZookeeperData",
       "ZookeeperLog",
     ];
    // zookeeperInput赋值是为了防止this在underscore中冲突
    var zookeeperInput = this.state.zookeeperInput;
    var RequredParameters = _.map(
      RequredFileds, function(filed){
        return zookeeperInput[filed];
      }
    );
    var index = _.indexOf(RequredParameters, "");
    // 判断是否缺少必要参数
    if (index !== -1) {
      Modalinfo.fire(<p>缺少必要参数！</p>);
    }
    else {
      $.post(url,
          JSON.stringify(this.state.zookeeperInput),
          function(data,status){
            console.log(data);
            data = JSON.parse(data);
            if (data.Status === "Success") {
              // 刷新zookeeper列表
              refreshZookeeperList();
              handleClose();
              Modalinfo.fire(<p>{data.Data}</p>);
            }else if (data.Status === "NoLogin") {
              CleanToken();
            }else {
              handleClose();
              Modalinfo.fire(<p>{data.Data}</p>);
            }
          });
    }
  }


  render() {
    const zookeeperInput = this.state.zookeeperInput;
    const FormInputList = [
      {label: "zookeeper安装目录", name: "DeployPath"},
      {label: "zookeeper数据目录", name: "ZookeeperData"},
      {label: "zookeeper日志目录", name: "ZookeeperLog"},
      {label: "zookeeper节点1", name: "ClusterNode01"},
      {label: "zookeeper节点2", name: "ClusterNode02"},
      {label: "zookeeper节点3", name: "ClusterNode03"},
    ];
    return (
      <Modal size="large" trigger={
        <Button
          floated='right'
          icon
          labelPosition='left'
          color='twitter'
          size='small'
          onClick={this.handleOpen}
        >
          <Icon name='add' />添加集群
        </Button>
      }
      open={this.state.modalOpen}
      onClose={this.handleClose}
      >
        <Modal.Header>
          <Image wrapped src='/zookeeper.svg' />
          Zookeeper集群部署
        </Modal.Header>
        <Modal.Content>
          <Modal.Description>
          <React.Fragment>
          <Header as='h3'>基础信息</Header>
          <Form>
            <FormFieldInput
              label="集群名称"
              name="ClusterName"
              value={zookeeperInput.ClusterName}
              onChange={this.handleInputChange}
              placeholder="请输入集群名称，此名称仅仅作为一个标识，如：zk-dubbo-account"
            />
            <Form.Field>
              <label>备注</label>
              <TextArea rows={2} name="Remark" value={zookeeperInput.Remark} onChange={this.handleInputChange} placeholder="备注一下吧" />
            </Form.Field>
          </Form>

          <Divider hidden />

          <Header as='h3'>配置信息</Header>
          <Form>
            <Form.Field>
              <label>zookeeper版本</label>
              <Dropdown
                value={zookeeperInput.ZookeeperVersion}
                placeholder='Select version'
                fluid
                selection
                options={zookeeperOptions}
                onChange={this.handleZookeeperVersionChange}
              />
            </Form.Field>
            {
              FormInputList.map(
                (formInput) =>
                <FormFieldInput
                  key={formInput.name}
                  label={formInput.label}
                  name={formInput.name}
                  value={zookeeperInput[formInput.name]}
                  onChange={this.handleInputChange}
                  placeholder={zookeeperInput[formInput.name]}
                />
              )
            }
          </Form>

          </React.Fragment>
          </Modal.Description>
        </Modal.Content>
        <Modal.Actions>
          <Button negative onClick={this.handleClose}>取消</Button>
          <Button
            className={this.state.checkboxLoading}
            onClick={this.handleInstall}
            positive
            icon='checkmark'
            content='提交任务'
          />
        </Modal.Actions>
      </Modal>


    );
  }
}

export default DeployZookeeperModal
