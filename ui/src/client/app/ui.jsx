import React from 'react';
import Endpoint from './endpoint/endpoint.jsx';
import CDC from './cdc/CDC.jsx';
import Navigation from './navigation.jsx';
import Toaster from './Toaster.jsx';
import ServiceProp from './propValidators/ServicePropValidator';

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
    this.setActiveEndpointRef = this.setActiveEndpointRef.bind(this);
    this.restore = this.restore.bind(this);
  }

  componentDidMount() {
    return this.restore();
  }

  setToasterRef(ref) {
    this.toaster = ref;
  }

  setCDCRef(ref) {
    this.cdc = ref;
  }

  setActiveEndpointRef(ref) {
    this.activeEndpoint = ref;
  }

  openEditor(endpointName) {
    this.setState({
      endpointService: this.state.endpointService.selectEndpoint(endpointName),
    });
  }

  endpoints() {
    if (this.state.endpointService) {
      const wee = this.state.endpointService.getEndpoints();
      return wee;
    }
    return [];
  }

  deleteEndpoint() {
    const endpointName = this.state.endpointService.getEndpoint().Name;
    return this.state.endpointService.deleteEndpoint()
      .tap(endpointService => this.setState({ endpointService }))
      .tap(() => this.toaster.alert(`Endpoint "${endpointName}" deleted`))
      .catch(err => {
        this.toaster.alert(['Problem deleting endpoint', err.toString()], Toaster.ErrorDisplayTime);
        return this.restore();
      });
  }

  updateEndpoint() {
    const newEndpointState = this.activeEndpoint.state;

    const update = {
      name: newEndpointState.name,
      cdcdisabled: newEndpointState.cdcDisabled,
      request: {
        uri: newEndpointState.uri,
        regexuri: newEndpointState.regex,
        method: newEndpointState.method,
        body: newEndpointState.reqBody,
        form: newEndpointState.form,
        headers: newEndpointState.reqHeaders,
      },
      response: {
        code: parseInt(newEndpointState.code, 10),
        body: newEndpointState.body,
        headers: newEndpointState.resHeaders,
      },
    };

    return this.state.endpointService.updateEndpoint(update)
      .tap(endpointService => this.setState({ endpointService }))
      .tap(() => this.toaster.alert(`Endpoint "${update.name}" updated`))
      .catch(err => {
        this.toaster.alert([
          'Problem updating endpoints, restoring state (you may lose some changes)',
          err.toString(),
        ], Toaster.ErrorDisplayTime);
        return this.restore();
      });
  }

  add() {
    return this.state.endpointService.addNewEndpoint()
      .tap(endpointService => this.setState({ endpointService }))
      .tap((endpointService) => {
        this.toaster.alert(`Endpoint "${endpointService.getEndpoint().Name}" added`);
      })
      .catch(err => {
        this.toaster.alert(
          ['Problem creating new endpoint', err.toString()],
          Toaster.ErrorDisplayTime);
        return this.restore();
      });
  }

  currentEndpointName() {
    if (this.state) {
      const currentEndpoint = this.state.endpointService.getEndpoint();
      return currentEndpoint ? currentEndpoint.Name : '';
    }

    return '';
  }

  restore() {
    return this.state.endpointService.init()
      .tap(endpointService => this.setState({ endpointService }))
      .catch(err => {
        this.toaster.alert(
          ['Problem getting MJ endpoints', err.toString()],
          Toaster.ErrorDisplayTime);
      });
  }

  renderCurrentEndpoint() {
    const endpoint = this.state.endpointService.getEndpoint();
    if (endpoint) {
      return (
        <Endpoint
          delete={this.deleteEndpoint}
          key={endpoint.Name}
          ref={this.setActiveEndpointRef}
          endpoint={endpoint}
          updateServer={this.updateEndpoint}
        />);
    }
    return null;
  }

  render() {
    return (
      <div className="mdl-layout mdl-js-layout mdl-layout--fixed-drawer">

        <CDC ref={this.setCDCRef} url="/mj-check-compatability" />

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

        <Toaster ref={this.setToasterRef} />

      </div>

    );
  }
}

UI.propTypes = {
  service: ServiceProp,
};

export default UI;
