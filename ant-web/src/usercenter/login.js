import React from 'react';
import $ from  'jquery'
import Cookies from 'universal-cookie';
import Backend from '../conf';
import 'semantic-ui-css/semantic.min.css';
import { Container, Button, Form, Grid, Segment } from 'semantic-ui-react';
import Modalinfo from '../home/modal-info';

const cookies = new Cookies();

class Logging extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      username: "",
      password: ""
    };
  }

  // 当用户名和密码发生变化时，给state赋值
  handleInputChange = (event) => {
    if (event.target.name==="username") {
      this.setState({username: event.target.value});
    }else if (event.target.name==="password") {
      this.setState({password: event.target.value});
    }
    else {
      Modalinfo.fire(<p>无效参数</p>);
    }
  }


  // 进行登录操作，验证成功则设置cookies，跳转到首页
  handleClick = () => {
    var url = Backend + "api/user/login/"
    var username = this.state.username;
    var password = this.state.password;
    $.post(url,
        JSON.stringify({
            username: username,
            password: password
        }),
        function(data,status){
          data = JSON.parse(data);
          if (data.Status === "Success") {
            cookies.set("username", username, { path: '/' });
            cookies.set("usertoken", data.Data, { path: '/' });
            window.location.href="/";
          } else {
            Modalinfo.fire(<p>登录失败</p>);
          }
        });
  }

  render() {
    return (
      <div style={{marginTop: 150 + "px"}}>
        <Container>
            <Grid columns='equal'>
            <Grid.Column>
            </Grid.Column>
            <Grid.Column>
              <Segment placeholder>
                <Grid columns={1} relaxed='very' stackable>
                  <Grid.Column>
                    <Form>
                      <Form.Input
                        icon='user'
                        iconPosition='left'
                        label='Username'
                        placeholder='Username' name="username" value={this.state.username} onChange={this.handleInputChange}
                      />
                      <Form.Input
                        icon='lock'
                        iconPosition='left'
                        label='Password'
                        type='password' name="password" value={this.state.password} onChange={this.handleInputChange}
                      />

                      <Button content='Login' primary onClick={this.handleClick} />
                    </Form>
                  </Grid.Column>
                </Grid>
              </Segment>
            </Grid.Column>
            <Grid.Column>
            </Grid.Column>
          </Grid>
        </Container>
      </div>

    );
  }
}

export default Logging;
