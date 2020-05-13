import React from 'react'
import { Message, Form } from 'semantic-ui-react'

export default class FormFieldMessage extends React.Component {

  constructor(props) {
    super(props);
    this.state={};
  }


  render() {
    return  <Form.Field>
              <Message>
                  <Message.Header>{this.props.header}</Message.Header>
                  <Message.List>
                    <Message.Item>{this.props.content}</Message.Item>
                  </Message.List>
              </Message>
            </Form.Field>
  }
}
