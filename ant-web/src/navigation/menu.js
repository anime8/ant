import React from 'react';
import 'semantic-ui-css/semantic.min.css';
import { Menu } from 'semantic-ui-react';

// 定义菜单栏，分为菜单和子菜单
const menus = [
  {
    menu: "部署",
    submenus: ["redis", "zookeeper", "kafka"]
  },
  {
    menu: "其他",
    submenus: ["other"]
  },
];

class MenuList extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      activeItem: 'redis'
    };
  }

  // 设置当前子菜单
  handleClick = (e, element) => {
    // 调用上层函数，设置子菜单
    this.props.handleItemClick(element.name);
    // 设置当前组件子菜单
    this.setState({activeItem: element.name});
  }

  render() {
    const listMenus = menus.map(
      (menu) =>
      <Menu.Item key={menu.menu}>
        <Menu.Header>{menu.menu}</Menu.Header>
        <Menu.Menu>
          { menu.submenus.map((submenu) =>
          <Menu.Item
            key={submenu}
            name={submenu}
            active={this.state.activeItem === submenu}
            onClick={this.handleClick}
          />
          )}
        </Menu.Menu>
      </Menu.Item>
    );
    return (
      <Menu fluid vertical secondary>
        {listMenus}
      </Menu>

    );
  }
}

export default MenuList;
