import React from 'react'
// import _ from 'underscore'
import { Table } from 'semantic-ui-react'
import DeployKafkaModal from './kafka-add-modal';
import ViewKafkaModal from './kafka-view-modal';
// import $ from  'jquery';
import axios from 'axios';
import Backend from '../conf';
import Modalinfo from '../home/modal-info';
import { CleanToken } from '../home/cookies';

class DeployKafka extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};
  }

  // 获取kafka信息，并放入state中
  async componentDidMount() {
    var url = Backend + "api/kafka/getall/";
    let KafkaList = await axios.get(url);
    KafkaList = KafkaList.data;
    if (KafkaList.Status === "Success"){
      var Data = JSON.parse(KafkaList.Data);
      this.setState({
          KafkaList: Data
      });
    }else if (KafkaList.Status === "NoLogin") {
      CleanToken();
    }else{
      Modalinfo.fire(<p>请求数据失败</p>);
    }
  }

  // 刷新kafka列表
  refreshKafkaList = async () => {
    var url = Backend + "api/kafka/getall/";
    let KafkaList = await axios.get(url)
    KafkaList = KafkaList.data;
    if (KafkaList.Status === "Success"){
      var Data = JSON.parse(KafkaList.Data);
      this.setState({
          KafkaList: Data
      });
    }else if (KafkaList.Status === "NoLogin") {
      CleanToken();
    }else{
      Modalinfo.fire(<p>刷新kafka列表失败</p>);
    }
  }

  render() {
    return (
      <div>
          <Table fixed color="olive">
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell>集群名称</Table.HeaderCell>
                <Table.HeaderCell>部署进度</Table.HeaderCell>
                <Table.HeaderCell>备注</Table.HeaderCell>
                <Table.HeaderCell>详情</Table.HeaderCell>
              </Table.Row>
            </Table.Header>

            <Table.Body>
              {this.state.KafkaList && this.state.KafkaList.map((kafka) => (
                <Table.Row key={kafka.Id}>
                  <Table.Cell>
                    {kafka.ClusterName}
                  </Table.Cell>
                  <Table.Cell>
                    {kafka.DeployStatus}
                  </Table.Cell>
                  <Table.Cell>
                    {kafka.Remark}
                  </Table.Cell>
                  <Table.Cell>
                    <ViewKafkaModal kafka={JSON.stringify(kafka)}  refreshKafkaList={this.refreshKafkaList} DeployStatus={kafka.DeployStatus}/>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
            <Table.Footer fullWidth>
              <Table.Row>
                <Table.HeaderCell />
                <Table.HeaderCell colSpan='3'>
                  <DeployKafkaModal refreshKafkaList={this.refreshKafkaList}/>
                </Table.HeaderCell>
              </Table.Row>
            </Table.Footer>
          </Table>
      </div>

    );
  }
}

export default DeployKafka
