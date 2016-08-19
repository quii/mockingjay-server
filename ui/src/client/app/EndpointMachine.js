import _ from 'lodash';
import { guid } from './util';

class EndpointMachine {
  constructor(endpoints) {
    this.endpoints = endpoints;
    this.selectedEndpointIndex = null;

    this.selectEndpoint = this.selectEndpoint.bind(this);
    this.getEndpoint = this.getEndpoint.bind(this);
    this.getEndpoints = this.getEndpoints.bind(this);
    this.updateEndpoint = this.updateEndpoint.bind(this);
    this.updateEndpoints = this.updateEndpoints.bind(this);
    this.addNewEndpoint = this.addNewEndpoint.bind(this);
    this.deleteEndpoint = this.deleteEndpoint.bind(this);
    this.asJSON = this.asJSON.bind(this);
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

    this.endpoints.unshift(newEndpoint);
    this.selectedEndpointIndex = 0;

    return this;
  }

  deleteEndpoint() {
    this.endpoints.splice(this.selectedEndpointIndex, 1);
    this.selectedEndpointIndex = null;
    return this;
  }

  updateEndpoint(updatedEndpoint) {
    this.endpoints[this.selectedEndpointIndex] = updatedEndpoint;
    return this;
  }

  updateEndpoints(update) {
    this.endpoints = update;
    return this;
  }

  asJSON() {
    return JSON.stringify(this.getEndpoints());
  }

}

export default EndpointMachine;
