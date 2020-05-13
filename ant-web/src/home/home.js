import React from 'react';
import 'semantic-ui-css/semantic.min.css';
import { Grid } from 'semantic-ui-react';
import MenuList from '../navigation/menu';
import DeployZookeeper from '../deploy/zookeeper';
import DeployRedis from '../deploy/redis';
import DeployKafka from '../deploy/kafka';


// 添加渲染，根据点击不同子菜单来显示内容
function Content (element) {
  if (element.activeItem === "zookeeper") {
    return <DeployZookeeper/>;
  } else if (element.activeItem === "redis") {
    return <DeployRedis/>;
  } else if (element.activeItem === "kafka") {
    return <DeployKafka/>;
  } else {
    return <p>敬请期待</p>;
  }
}

// app主要内容，在这里进行了显示划分，右菜单2格，左边内容14格
class App extends React.Component {
  constructor(props) {
      super(props);
      this.state = {
        activeItem: 'redis'
      };
    }

    // 设置当前选择的子菜单
    handleItemClick = (activeItem) => {
      this.setState({activeItem: activeItem});
    }

    render() {
      return (
          <Grid>
            <Grid.Row>
              <Grid.Column width={14}>
                <Content activeItem={this.state.activeItem}/>
              </Grid.Column>
              <Grid.Column width={2}>
                <MenuList handleItemClick={this.handleItemClick}/>
              </Grid.Column>
            </Grid.Row>
          </Grid>
      );
    }

}

export default App;
