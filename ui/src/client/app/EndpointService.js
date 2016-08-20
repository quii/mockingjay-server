import _ from 'lodash';
import { guid } from './util';

class EndpointService {
  constructor(api) {
    this.endpoints = [];
    this.selectedEndpointIndex = null;
    this.api = api;

    this.selectEndpoint = this.selectEndpoint.bind(this);
    this.getEndpoint = this.getEndpoint.bind(this);
    this.getEndpoints = this.getEndpoints.bind(this);
    this.updateEndpoint = this.updateEndpoint.bind(this);
    this.addNewEndpoint = this.addNewEndpoint.bind(this);
    this.deleteEndpoint = this.deleteEndpoint.bind(this);
    this.init = this.init.bind(this);
    this.sendUpdateToServer = this.sendUpdateToServer.bind(this);
  }

  init() {
    return this.api.getEndpoints()
      .then(endpoints => {
        this.selectedEndpointIndex = null;
        this.endpoints = endpoints;
      })
      .then(() => this);
  }

  sendUpdateToServer() {
    return this.api.updateEndpoints(JSON.stringify(this.getEndpoints()))
      .then(endpoints => {
        this.endpoints = endpoints;
      })
      .then(() => this);
  }

  selectEndpoint(name) {
    this.selectedEndpointIndex = _.findIndex(this.endpoints, ep => ep.Name === name);
    return this;
  }

  getEndpoint() {
    return this.endpoints[this.selectedEndpointIndex];
  }

  getEndpoints() {
    return this.endpoints;
  }

  addNewEndpoint() {
    const newEndpoint = {
      Name: guid(),
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

    this.endpoints.unshift(newEndpoint);
    this.selectedEndpointIndex = 0;

    return this.sendUpdateToServer();
  }

  deleteEndpoint() {
    this.endpoints.splice(this.selectedEndpointIndex, 1);
    this.selectedEndpointIndex = null;
    return this.sendUpdateToServer();
  }

  updateEndpoint(updatedEndpoint) {
    this.endpoints[this.selectedEndpointIndex] = updatedEndpoint;
    return this.sendUpdateToServer();
  }

}

export default EndpointService;
