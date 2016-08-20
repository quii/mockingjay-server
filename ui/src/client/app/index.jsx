import React from 'react';

import ReactDOM from 'react-dom';
import Endpoint from './endpoint/endpoint.jsx';
import CDC from './cdc/CDC.jsx';
import Navigation from './navigation.jsx';
import API from './API.js';
import EndpointMachine from './EndpointMachine';
import Toaster from './Toaster.jsx';

class UI extends React.Component {

  constructor(props) {
    super(props);

    this.api = new API(props.url);

    this.state = {
      endpointMachine: new EndpointMachine([]),
    };

    this.componentDidMount = this.componentDidMount.bind(this);
    this.putUpdate = this.putUpdate.bind(this);
    this.openEditor = this.openEditor.bind(this);
    this.add = this.add.bind(this);
    this.deleteEndpoint = this.deleteEndpoint.bind(this);
    this.updateEndpoint = this.updateEndpoint.bind(this);
    this.renderCurrentEndpoint = this.renderCurrentEndpoint.bind(this);
    this.currentEndpointName = this.currentEndpointName.bind(this);
  }

  componentDidMount() {
    return this.api.getEndpoints()
      .tap(data => this.setState({ endpointMachine: new EndpointMachine(data) }))
      .catch(err => console.error('Problem getting MJ endpoints', status, err.toString()));
  }

  putUpdate() {
    return this.api.updateEndpoints(this.state.endpointMachine.asJSON())
      .tap(data => {
        this.setState({ endpointMachine: this.state.endpointMachine.updateEndpoints(data) });
      })
      .then(() => this.cdc.checkCompatability())
      .catch(err => this.toaster.alert(`Problem with PUT to ${this.props.url} ${err.toString()}`));
  }

  currentEndpointName() {
    if (this.state) {
      const currentEndpoint = this.state.endpointMachine.getEndpoint();
      return currentEndpoint ? currentEndpoint.Name : '';
    }
    return '';
  }

  add() {
    this.setState({
      endpointMachine: this.state.endpointMachine.addNewEndpoint(),
    });

    return this.putUpdate();
  }

  deleteEndpoint() {
    this.setState({
      endpointMachine: this.state.endpointMachine.deleteEndpoint(),
    });
    this.toaster.alert('Endpoint deleted');
    return this.putUpdate();
  }

  updateEndpoint() {
    const newEndpointState = this.activeEndpoint.state;

    const update = {
      Name: newEndpointState.name,
      CDCDisabled: newEndpointState.cdcDisabled,
      Request: {
        URI: newEndpointState.uri,
        RegexURI: newEndpointState.regex,
        Method: newEndpointState.method,
        Body: newEndpointState.reqBody,
        Form: newEndpointState.form,
        Headers: newEndpointState.reqHeaders,
      },
      Response: {
        Code: parseInt(newEndpointState.code, 10),
        Body: newEndpointState.body,
        Headers: newEndpointState.resHeaders,
      },
    };

    this.setState({
      endpointMachine: this.state.endpointMachine.updateEndpoint(update),
    });

    this.putUpdate();
  }

  openEditor(endpointName) {
    this.setState({
      endpointMachine: this.state.endpointMachine.selectEndpoint(endpointName),
    });
  }

  renderCurrentEndpoint() {
    const endpoint = this.state.endpointMachine.getEndpoint();
    if (endpoint) {
      return (
        <Endpoint
          delete={this.deleteEndpoint}
          key={endpoint.Name}
          ref={r => {
            this.activeEndpoint = r;
          }}
          cdcDisabled={endpoint.CDCDisabled}
          updateServer={this.updateEndpoint}
          name={endpoint.Name}
          method={endpoint.Request.Method}
          reqBody={endpoint.Request.Body}
          uri={endpoint.Request.URI}
          regex={endpoint.Request.RegexURI}
          reqHeaders={endpoint.Request.Headers}
          form={endpoint.Request.Form}
          code={endpoint.Response.Code}
          body={endpoint.Response.Body}
          resHeaders={endpoint.Response.Headers}
        />);
    }
    return null;
  }

  render() {
    return (
      <div className="mdl-layout mdl-js-layout mdl-layout--fixed-drawer">

        <CDC
          ref={r => {
            this.cdc = r;
          }}
          url="/mj-check-compatability"
        />

        <div className="mdl-layout__drawer">
          <h1 className="mdl-layout-title mdl-color-text--primary">mockingjay server</h1>
          <Navigation
            openEditor={this.openEditor}
            addEndpoint={this.add}
            endpoints={this.state.endpointMachine.getEndpoints()}
            activeEndpoint={this.currentEndpointName()}
          />
        </div>

        <main className="mdl-layout__content mdl-color--grey-100">
          <div className="page-content">
            {this.renderCurrentEndpoint()}
          </div>
        </main>

        <Toaster
          ref={r => {
            this.toaster = r;
          }}
        />

      </div>

    );
  }
}

ReactDOM.render(
  <UI url="/mj-endpoints"/>,
  document.getElementById('app')
);
