import React from 'react';

import View from './view.jsx';
import Form from './form.jsx';

class Endpoint extends React.Component {
  constructor(props) {
    super(props);

    this.setState({
      index: this.props.index,
      cdcDisabled: this.props.cdcDisabled,
      isEditing: false,
      name: this.props.name,
      method: this.props.method,
      uri: this.props.uri,
      regex: this.props.regex,
      reqBody: this.props.reqBody,
      form: this.props.form,
      reqHeaders: this.props.reqHeaders,
      code: this.props.code,
      body: this.props.body,
      resHeaders: this.props.resHeaders,
    });

    this.startEditing = this.startEditing.bind(this);
    this.delete = this.delete.bind(this);
    this.finishEditing = this.finishEditing.bind(this);
    this.updateValue = this.updateValue.bind(this);
    this.updateCheckbox = this.updateCheckbox.bind(this);
    this.render = this.render.bind(this);
  }

  startEditing() {
    this.setState({
      isEditing: true,
    });
  }

  delete() {
    this.props.delete();
  }

  finishEditing() {
    this.setState({
      isEditing: false,
    });
    this.props.updateServer();
  }

  updateValue(e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  updateCheckbox(e) {
    this.setState({
      [e.target.name]: e.target.value === 'on',
    });
  }

  render() {
    if (this.state.isEditing) {
      return (<Form
        name={this.state.name}
        finishEditing={this.finishEditing}
        originalValues={this.state}
        onChange={this.updateValue}
        onCheckboxChange={this.updateCheckbox}
      />);
    }

    return (
      <View
        method={this.state.method}
        uri={this.state.uri}
        regex={this.state.regex}
        reqBody={this.state.reqBody}
        reqHeaders={this.state.reqHeaders}
        form={this.state.form}
        code={this.state.code.toString()}
        body={this.state.body}
        resHeaders={this.state.resHeaders}
        name={this.state.name}
        startEditing={this.startEditing}
        deleteEndpoint={this.delete}
      />
    );
  }
}

export default Endpoint;
