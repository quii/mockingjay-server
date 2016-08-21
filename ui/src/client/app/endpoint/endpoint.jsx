import React from 'react';

import View from './view.jsx';
import Form from './form.jsx';

class Endpoint extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      cdcDisabled: props.endpoint.CDCDisabled,
      isEditing: false,
      name: props.endpoint.Name,
      method: props.endpoint.Request.Method,
      uri: props.endpoint.Request.URI,
      regex: props.endpoint.Request.RegexURI,
      reqBody: props.endpoint.Request.Body,
      form: props.endpoint.Request.Form,
      reqHeaders: props.endpoint.Request.Headers,
      code: props.endpoint.Response.Code,
      body: props.endpoint.Response.Body,
      resHeaders: props.endpoint.Response.Headers,
    };

    this.updateServer = props.updateServer;
    this.delete = props.delete;

    this.startEditing = this.startEditing.bind(this);
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

  finishEditing() {
    this.setState({
      isEditing: false,
    });
    this.updateServer();
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

Endpoint.propTypes = {
  endpoint: React.PropTypes.object,
  delete: React.PropTypes.func,
  updateServer: React.PropTypes.func,
};

export default Endpoint;
