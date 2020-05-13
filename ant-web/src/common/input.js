import React from 'react'
import { Form, Input } from 'semantic-ui-react'



export default class FormFieldInput extends React.Component {

  constructor(props) {
    super(props);
    this.state={};
  }


  render() {
    return  <Form.Field>
              <label>{this.props.label}</label>
              <Input name={this.props.name} value={this.props.value} onChange={this.props.onChange} placeholder={this.props.placeholder} />
            </Form.Field>
  }
}
