import React from 'react';
import _ from 'lodash';
import request from 'superagent';

import ReactDOM from 'react-dom';
import Endpoint from './endpoint/endpoint.jsx';
import CDC from './cdc/CDC.jsx';
import { guid } from './util';
import Navigation from './navigation.jsx';

const UI = React.createClass({
  getInitialState() {
    return {
      data: [],
      activeEndpoint: null,
      endpointIds: [],
    };
  },
  componentDidMount() {
    request
      .get(this.props.url)
      .end((err, res) => {
        if (err) {
          console.error(this.props.url, status, err.toString());
        } else {
          this.setState({ data: JSON.parse(res.text) });
        }
      });
  },
  putUpdate(update) {
    request
      .put(this.props.url)
      .send(update)
      .end((err, res) => {
        if (err) {
          this.toasty(`Problem with PUT to ${this.props.url} ${err.toString()}`);
          console.error(this.props.url, status, err.toString());
        } else {
          this.refs.cdc.checkCompatability();
          this.setState({ data: JSON.parse(res.text) });
        }
      });
  },
  add() {
    const data = _.cloneDeep(this.state.data);

    const newEndpointName = guid();

    const newEndpoint = {
      Name: newEndpointName,
      CDCDisabled: false,
      Request: {
        URI: '/hello',
        Method: 'GET',
      },
      Response: {
        Code: 200,
        Body: 'World!',
      },
    };

    data.unshift(newEndpoint);

    this.setState({
      data,
      activeEndpoint: newEndpointName,
      endpointIds: [],
    });

    const json = JSON.stringify(data);

    this.putUpdate(json);
  },
  toasty(msg) {
    const notification = document.querySelector('.mdl-js-snackbar');
    notification.MaterialSnackbar.showSnackbar(
      {
        message: msg,
      }
    );
  },
  openEditor(endpointName) {
    this.setState({
      activeEndpoint: endpointName,
    });
  },
  deleteEndpoint() {
    const indexToDelete = this.refs[this.state.activeEndpoint].state.index;

    const data = _.cloneDeep(this.state.data);
    data.splice(indexToDelete, 1);
    const json = JSON.stringify(data);

    this.toasty('Endpoint deleted');

    this.putUpdate(json);
  },
  updateServer() {
    const newEndpointState = this.refs[this.state.activeEndpoint].state;

    const data = _.cloneDeep(this.state.data);

    data[newEndpointState.index] = {
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
        Code: parseInt(newEndpointState.code),
        Body: newEndpointState.body,
        Headers: newEndpointState.resHeaders,
      },
    };

    const json = JSON.stringify(data);

    this.setState({
      activeEndpoint: newEndpointState.name,
    });

    this.putUpdate(json);
  },
  renderCurrentEndpoint() {
    if (this.state.activeEndpoint) {
      const index = _.findIndex(this.state.data, ep => ep.Name == this.state.activeEndpoint);
      const endpoint = this.state.data.find(ep => ep.Name === this.state.activeEndpoint);
      if (endpoint) {
        return (
          <Endpoint
            index={index}
            delete={this.deleteEndpoint}
            key={endpoint.Name}
            ref={endpoint.Name}
            cdcDisabled={endpoint.CDCDisabled}
            updateServer={this.updateServer}
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
    }
    return null;
  },
  render() {
    return (
      <div className="mdl-layout mdl-js-layout mdl-layout--fixed-drawer">

        <CDC ref="cdc" url="/mj-check-compatability"/>
        <div className="mdl-layout__drawer">
          <h1 className="mdl-layout-title mdl-color-text--primary">mockingjay server</h1>
          <Navigation
            openEditor={this.openEditor}
            addEndpoint={this.add}
            endpoints={this.state.data}
            activeEndpoint={this.state.activeEndpoint}
          />
        </div>
        <main className="mdl-layout__content mdl-color--grey-100">
          <div className="page-content">
            {this.renderCurrentEndpoint()}
          </div>
        </main>
        <div aria-live="assertive" aria-atomic="true" aria-relevant="text"
             className="mdl-snackbar mdl-js-snackbar">
          <div className="mdl-snackbar__text"></div>
          <button type="button" className="mdl-snackbar__action"></button>
        </div>
      </div>

    );
  },
});
ReactDOM.render(
  <UI url="/mj-endpoints"/>,
  document.getElementById('app')
);
