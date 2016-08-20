import React from 'react';
import Endpoint from './endpoint/endpoint.jsx';
import CDC from './cdc/CDC.jsx';
import Navigation from './navigation.jsx';
import Toaster from './Toaster.jsx';

class UI extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      endpointService: props.service,
    };

    this.componentDidMount = this.componentDidMount.bind(this);
    this.openEditor = this.openEditor.bind(this);
    this.add = this.add.bind(this);
    this.deleteEndpoint = this.deleteEndpoint.bind(this);
    this.updateEndpoint = this.updateEndpoint.bind(this);
    this.renderCurrentEndpoint = this.renderCurrentEndpoint.bind(this);
    this.currentEndpointName = this.currentEndpointName.bind(this);
    this.endpoints = this.endpoints.bind(this);
    this.setCDCRef = this.setCDCRef.bind(this);
    this.setToasterRef = this.setToasterRef.bind(this);
  }

  componentDidMount() {
    return this.state.endpointService.init()
      .tap(endpointService => this.setState({ endpointService }))
      .catch(err => this.toaster.alert('Problem getting MJ endpoints', err.toString()));
  }

  currentEndpointName() {
    if (this.state) {
      const currentEndpoint = this.state.endpointService.getEndpoint();
      return currentEndpoint ? currentEndpoint.Name : '';
    }

    return '';
  }

  add() {
    return this.state.endpointService.addNewEndpoint()
      .tap(endpointService => this.setState({ endpointService }))
      .tap((endpointService) => {
        this.toaster.alert(`Endpoint "${endpointService.getEndpoint().Name}" added`);
      })
      .catch(err => this.toaster.alert('Problem saving new endpoint', err.toString()));
  }

  deleteEndpoint() {
    const endpointName = this.state.endpointService.getEndpoint().Name;
    return this.state.endpointService.deleteEndpoint()
      .tap(endpointService => this.setState({ endpointService }))
      .tap(() => this.toaster.alert(`Endpoint "${endpointName}" deleted`))
      .catch(err => this.toaster.alert('Problem deleting endpoint', err.toString()));
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

    return this.state.endpointService.updateEndpoint(update)
      .tap(endpointService => this.setState({ endpointService }))
      .tap(() => this.toaster.alert(`Endpoint "${update.Name}" updated`))
      .catch(err => this.toaster.alert('Problem updating endpoints', err.toString()));
  }

  openEditor(endpointName) {
    this.setState({
      endpointService: this.state.endpointService.selectEndpoint(endpointName)
    });
  }

  endpoints() {
    if (this.state.endpointService) {
      return this.state.endpointService.getEndpoints();
    }
    return [];
  }

  setToasterRef(ref) {
    this.toaster = ref;
  }

  setCDCRef(ref) {
    this.cdc = ref;
  }

  renderCurrentEndpoint() {
    const endpoint = this.state.endpointService.getEndpoint();
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

        <CDC ref={this.setCDCRef} url="/mj-check-compatability"/>

        <div className="mdl-layout__drawer">
          <h1 className="mdl-layout-title mdl-color-text--primary">mockingjay server</h1>
          <Navigation
            openEditor={this.openEditor}
            addEndpoint={this.add}
            endpoints={this.endpoints()}
            activeEndpoint={this.currentEndpointName()}
          />
        </div>

        <main className="mdl-layout__content mdl-color--grey-100">
          <div className="page-content">
            {this.renderCurrentEndpoint()}
          </div>
        </main>

        <Toaster ref={this.setToasterRef}/>

      </div>

    );
  }
}

UI.propTypes = {
  service: React.PropTypes.object,
};

export default UI;
