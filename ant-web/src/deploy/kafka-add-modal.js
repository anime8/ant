import React from 'react'
import { Icon, Button, Header, Image, Modal, Divider, Form, TextArea, Dropdown } from 'semantic-ui-react'
import $ from  'jquery'
import Backend from '../conf';
import Modalinfo from '../home/modal-info';
import _ from 'underscore';
import { CleanToken } from '../home/cookies';
import FormFieldInput from '../common/input';

// kafka版本列表
const kafkaOptions = [
  {
    key: 'kafka_2.11-1.1.0',
    text: 'kafka_2.11-1.1.0',
    value: 'kafka_2.11-1.1.0',
    image: { avatar: true, src: '/kafka.png' },
  },
];


// 部署kafka弹出框
class DeployKafkaModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      checkboxLoading: "",
      modalOpen: false,
      kafkaInput: {
        ClusterName: "",
        Remark: "",
        ClusterNode01: "",
        ClusterNode02: "",
        ClusterNode03: "",
        KafkaVersion: "kafka_2.11-1.1.0",
        KafkaPath: "/opt/app",
        KafkaData: "/opt/data/kafka",
        KafkaZookeeper: "zk1:2181,zk2:2181,zk3:2181",
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
      kafkaInput: {
        ...this.state.kafkaInput,
        [event.target.name]: event.target.value
      }
    });
  }

  // 获取kafka版本信息
  handleKafkaVersionChange = (event) => {
    this.setState({
      kafkaInput: {
        ...this.state.kafkaInput,
        KafkaVersion: event.target.innerText
      }
    });
  }

  // 设置check
  handleChecked = (e, { name, checked }) => {
    this.setState({
      kafkaInput: {
        ...this.state.kafkaInput,
        [name]: checked
        }
    });
  }

  // 点击提交部署任务
  handleInstall = () => {
    // 部署按钮开始loading
    this.setState({ checkboxLoading: "loading disabled" });
    var handleClose = this.handleClose;
    // 设置变量，后面在post里面调用，这是为了防止和post返回的this起冲突
    var refreshKafkaList = this.props.refreshKafkaList;
    var url = Backend + "api/kafka/install/";
    // 判断是否缺少必要参数，如果缺少则提示用户
    var RequredFileds = [
      "ClusterName",
      "ClusterNode01",
      "ClusterNode02",
      "ClusterNode03",
       "KafkaVersion",
       "KafkaPath",
       "KafkaData",
       "KafkaZookeeper",
     ];
    // kafkaInput赋值是为了防止this在underscore中冲突
    var kafkaInput = this.state.kafkaInput;
    var RequredParameters = _.map(
      RequredFileds, function(filed){
        return kafkaInput[filed];
      }
    );
    var index = _.indexOf(RequredParameters, "");
    // 判断是否缺少必要参数
    if (index !== -1) {
      Modalinfo.fire(<p>缺少必要参数！</p>);
    }
    else {
      $.post(url,
          JSON.stringify(this.state.kafkaInput),
          function(data,status){
            console.log(data);
            data = JSON.parse(data);
            if (data.Status === "Success") {
              // 刷新kafka列表
              refreshKafkaList();
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
    const kafkaInput = this.state.kafkaInput;
    const FormInputList = [
      {label: "kafka安装目录", name: "KafkaPath"},
      {label: "kafka数据目录", name: "KafkaData"},
      {label: "Zookeeper配置", name: "KafkaZookeeper"},
      {label: "kafka节点1", name: "ClusterNode01"},
      {label: "kafka节点2", name: "ClusterNode02"},
      {label: "kafka节点3", name: "ClusterNode03"},
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
          <Image wrapped src='/kafka.png' size='tiny' />
          Kafka集群部署
        </Modal.Header>
        <Modal.Content>
          <Modal.Description>
          <React.Fragment>
          <Header as='h3'>基础信息</Header>
          <Form>
            <FormFieldInput
              label="集群名称"
              name="ClusterName"
              value={kafkaInput.ClusterName}
              onChange={this.handleInputChange}
              placeholder="请输入集群名称，此名称仅仅作为一个标识，如：kafka-account"
            />
            <Form.Field>
              <label>备注</label>
              <TextArea rows={2} name="Remark" value={kafkaInput.Remark} onChange={this.handleInputChange} placeholder="备注一下吧" />
            </Form.Field>
          </Form>

          <Divider hidden />

          <Header as='h3'>配置信息</Header>
          <Form>
            <Form.Field>
              <label>kafka版本</label>
              <Dropdown
                value={kafkaInput.KafkaVersion}
                placeholder='Select version'
                fluid
                selection
                options={kafkaOptions}
                onChange={this.handleKafkaVersionChange}
              />
            </Form.Field>
            {
              FormInputList.map(
                (formInput) =>
                <FormFieldInput
                  key={formInput.name}
                  label={formInput.label}
                  name={formInput.name}
                  value={kafkaInput[formInput.name]}
                  onChange={this.handleInputChange}
                  placeholder={kafkaInput[formInput.name]}
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

export default DeployKafkaModal
