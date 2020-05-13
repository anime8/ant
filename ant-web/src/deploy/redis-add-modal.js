import React from 'react'
import { Icon, Table, Button, Header, Image, Modal, Divider, Form, Input, TextArea, Grid, Radio, Dropdown } from 'semantic-ui-react'
import $ from  'jquery'
import Backend from '../conf';
import Modalinfo from '../home/modal-info';
import _ from 'underscore';
import { CleanToken } from '../home/cookies';
import FormFieldInput from '../common/input';

// redis版本列表
const redisOptions = [
  {
    key: 'redis-2.8.10',
    text: 'redis-2.8.10',
    value: 'redis-2.8.10',
    image: { avatar: true, src: '/redis-white.png' },
  },
  {
    key: 'redis-3.2.0',
    text: 'redis-3.2.0',
    value: 'redis-3.2.0',
    image: { avatar: true, src: '/redis-white.png' },
  },
  {
    key: 'redis-5.0.5',
    text: 'redis-5.0.5',
    value: 'redis-5.0.5',
    image: { avatar: true, src: '/redis-white.png' },
  },
];


// 部署redis弹出框
class DeployRedisModal extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      checkboxLoading: "",
      modalOpen: false,
      redisInput: {
        ClusterName: "",
        Remark: "",
        ClusterNode01: "",
        ClusterNodeRedisChecked01: true,
        ClusterNodeSentinelChecked01: true,
        ClusterNode02: "",
        ClusterNodeRedisChecked02: true,
        ClusterNodeSentinelChecked02: true,
        ClusterNode03: "",
        ClusterNodeRedisChecked03: true,
        ClusterNodeSentinelChecked03: true,
        RedisVersion: "redis-2.8.10",
        SentinelName: "mymaster",
        RedisData: "/opt/data/redis",
        SentinelData: "/opt/data/redis/sentinel",
        RedisLog: "/logs/redis",
        RedisConf: "/etc/redis",
        RedisAuthentication: false,
        RedisPassword: "",
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
      redisInput: {
        ...this.state.redisInput,
        [event.target.name]: event.target.value
      }
    });
  }

  // 获取redis版本信息
  handleRedisVersionChange = (event) => {
    this.setState({
      redisInput: {
        ...this.state.redisInput,
        RedisVersion: event.target.innerText
      }
    });
  }

  // 设置check
  handleChecked = (e, { name, checked }) => {
    this.setState({
      redisInput: {
        ...this.state.redisInput,
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
    var refreshRedisList = this.props.refreshRedisList;
    var url = Backend + "api/redis/install/";
    // 判断是否缺少必要参数，如果缺少则提示用户
    var RequredFileds = [
      "ClusterName",
      "ClusterNode01",
      "ClusterNode02",
      "ClusterNode03",
       "RedisVersion",
       "SentinelName",
       "RedisData",
       "SentinelData",
       "RedisLog",
       "RedisConf"
     ];
    // redisInput赋值是为了防止this在underscore中冲突
    var redisInput = this.state.redisInput;
    var RequredParameters = _.map(
      RequredFileds, function(filed){
        return redisInput[filed];
      }
    );
    var index = _.indexOf(RequredParameters, "");
    // 判断是否缺少必要参数
    if (index !== -1) {
      Modalinfo.fire(<p>缺少必要参数！</p>);
    }
    else {
      $.post(url,
          JSON.stringify(this.state.redisInput),
          function(data,status){
            console.log(data);
            data = JSON.parse(data);
            if (data.Status === "Success") {
              // 刷新redis列表
              refreshRedisList();
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
    const redisInput = this.state.redisInput;
    const FormInputList = [
      {label: "哨兵名称", name: "SentinelName"},
      {label: "redis数据目录", name: "RedisData"},
      {label: "sentinel数据目录", name: "SentinelData"},
      {label: "redis日志目录", name: "RedisLog"},
      {label: "redis配置目录", name: "RedisConf"},
    ];
    const CheckInputList = [
      {name: "ClusterNode01", redis: "ClusterNodeRedisChecked01", sentinel: "ClusterNodeSentinelChecked01"},
      {name: "ClusterNode02", redis: "ClusterNodeRedisChecked02", sentinel: "ClusterNodeSentinelChecked02"},
      {name: "ClusterNode03", redis: "ClusterNodeRedisChecked03", sentinel: "ClusterNodeSentinelChecked03"},
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
          <Image wrapped src='/redis-white.png' size='mini' />
          Rdis集群部署
        </Modal.Header>
        <Modal.Content>
          <Modal.Description>
          <React.Fragment>
          <Header as='h3'>基础信息</Header>
          <Form>
            <FormFieldInput
              label="集群名称"
              name="ClusterName"
              value={redisInput.ClusterName}
              onChange={this.handleInputChange}
              placeholder="请输入集群名称，此名称仅仅作为一个标识，如：redis-account"
            />
            <Form.Field>
              <label>备注</label>
              <TextArea rows={2} name="Remark" value={redisInput.Remark} onChange={this.handleInputChange} placeholder="备注一下吧" />
            </Form.Field>
          </Form>

          <Divider hidden />

          <Header as='h3'>配置信息</Header>
          <Form>
            <Form.Field>
              <label>redis版本</label>
              <Dropdown
                value={redisInput.RedisVersion}
                placeholder='Select version'
                fluid
                selection
                options={redisOptions}
                onChange={this.handleRedisVersionChange}
              />
            </Form.Field>
            {
              FormInputList.map(
                (formInput) =>
                <FormFieldInput
                  key={formInput.name}
                  label={formInput.label}
                  name={formInput.name}
                  value={redisInput[formInput.name]}
                  onChange={this.handleInputChange}
                  placeholder={redisInput[formInput.name]}
                />
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
                            (checkInput) =>
                            <Table.Row key={checkInput.name}>
                              <Table.Cell collapsing>
                                <Input name={checkInput.name} value={redisInput[checkInput.name]} onChange={this.handleInputChange} placeholder="请输入集群节点IP" />
                              </Table.Cell>
                              <Table.Cell>
                                <Radio
                                  name={checkInput.redis}
                                  checked={redisInput[checkInput.redis]}
                                  label='redis'
                                  onClick={this.handleChecked}
                                />
                              </Table.Cell>
                              <Table.Cell>
                                <Radio
                                  name={checkInput.sentinel}
                                  checked={redisInput[checkInput.sentinel]}
                                  label='sentinel'
                                  onClick={this.handleChecked}
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
                            name="RedisAuthentication"
                            checked={redisInput.RedisAuthentication}
                            label='密码验证'
                            onClick={this.handleChecked}
                          />
                        </Table.Cell>
                        <Table.Cell>
                          <Input disabled={!redisInput.RedisAuthentication} name="RedisPassword" value={redisInput.RedisPassword} onChange={this.handleInputChange} placeholder='请输入redis密码' />
                        </Table.Cell>
                      </Table.Row>
                    </Table.Body>
                  </Table>
                </Grid.Column>
              </Grid>
            </Form.Field>
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

export default DeployRedisModal
