import React from 'react'
// import _ from 'underscore'
import { Table } from 'semantic-ui-react'
import DeployRedisModal from './redis-add-modal';
import ViewRedisModal from './redis-view-modal';
// import $ from  'jquery';
import axios from 'axios';
import Backend from '../conf';
import Modalinfo from '../home/modal-info';
import { CleanToken } from '../home/cookies';

class DeployRedis extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};
  }

  // 获取redis信息，并放入state中
  async componentDidMount() {
    var url = Backend + "api/redis/getall/";
    let RedisList = await axios.get(url);
    RedisList = RedisList.data;
    if (RedisList.Status === "Success"){
      var Data = JSON.parse(RedisList.Data);
      this.setState({
          RedisList: Data
      });
    }else if (RedisList.Status === "NoLogin") {
      CleanToken();
    }else{
      Modalinfo.fire(<p>请求数据失败</p>);
    }
  }

  // 刷新redis列表
  refreshRedisList = async () => {
    var url = Backend + "api/redis/getall/";
    let RedisList = await axios.get(url)
    RedisList = RedisList.data;
    if (RedisList.Status === "Success"){
      var Data = JSON.parse(RedisList.Data);
      this.setState({
          RedisList: Data
      });
    }else if (RedisList.Status === "NoLogin") {
      CleanToken();
    }else{
      Modalinfo.fire(<p>刷新redis列表失败</p>);
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
              {this.state.RedisList && this.state.RedisList.map((redis) => (
                <Table.Row key={redis.Id}>
                  <Table.Cell>
                    {redis.ClusterName}
                  </Table.Cell>
                  <Table.Cell>
                    {redis.DeployStatus}
                  </Table.Cell>
                  <Table.Cell>
                    {redis.Remark}
                  </Table.Cell>
                  <Table.Cell>
                    <ViewRedisModal redis={JSON.stringify(redis)}  refreshRedisList={this.refreshRedisList} DeployStatus={redis.DeployStatus}/>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
            <Table.Footer fullWidth>
              <Table.Row>
                <Table.HeaderCell />
                <Table.HeaderCell colSpan='3'>
                  <DeployRedisModal refreshRedisList={this.refreshRedisList}/>
                </Table.HeaderCell>
              </Table.Row>
            </Table.Footer>
          </Table>
      </div>

    );
  }
}

export default DeployRedis
