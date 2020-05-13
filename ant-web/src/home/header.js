import React from 'react'
import { Header, Segment, Image, Dropdown, Icon } from 'semantic-ui-react'
import { CleanToken } from './cookies';
import Modalinfo from '../home/modal-info';
import Cookies from 'universal-cookie';


const cookies = new Cookies();

const LogoImage = () => (
  <Image src="/logo.jpg" size='small' />
)

const trigger = (
  <span>
    <Icon name='user' /> Hello, {cookies.get("username")}
  </span>
)

const options = [
  {
    key: 'user',
    text: (
      <span>
        Signed in as <strong>{cookies.get("username")}</strong>
      </span>
    ),
    disabled: true,
  },
  { key: 'profile', text: '用户中心' },
  { key: 'change-password', text: '更改密码' },
  { key: 'sign-out', text: '退出登录' },
]

class DropdownTrigger extends React.Component {
  constructor(props) {
      super(props);
      this.state = {
        activeItem: '用户中心'
      };
    }

    // 设置当前子菜单
    handleChange = (event) => {
      this.setState({activeItem: event.target.innerText});
      if (event.target.innerText === "用户中心") {
      } else if (event.target.innerText === "更改密码") {
      } else if (event.target.innerText === "退出登录") {
        CleanToken();
      } else {
        Modalinfo.fire(<p>无效参数</p>);
      }
    }

    render() {
       if (cookies.get("username") !== undefined) {
         return (
           <Dropdown trigger={trigger} value={this.state.activeItem} options={options} onChange={this.handleChange} />
         );
       } else {
         return (
           <div/>
         );
       }
    }
}

const HeaderFloating = () => (
  <Segment style={{backgroundColor: "#F5F5F5"}} clearing>
    <Header as='h2' floated='right'>
      <DropdownTrigger />
    </Header>
    <Header as='h2' floated='left'>
      <LogoImage /> Big Ant
    </Header>
  </Segment>
)

export default HeaderFloating
