import {expect} from 'chai';
import _ from 'lodash';
import Promise from 'bluebird';
import sinon from 'sinon';

import EndpointService from '../../../src/client/app/EndpointService';

const api = {
  getEndpoints: () => {
  },
  updateEndpoints: () => {
  },
};

describe('Endpoint serverice', () => {
  let sandbox;

  beforeEach(() => {
    sandbox = sinon.sandbox.create();
  });

  afterEach(() => {
    sandbox.restore();
  });


  it('gets endpoints and then you can select them', () => {
    const service = new EndpointService(api);

    const someEndpoints = [
      {
        Name: '123',
        Value: 'cat',
      },
      {
        Name: '456',
        Value: 'dog',
      },
    ];

    sandbox.stub(api, 'getEndpoints').returns(Promise.resolve(someEndpoints));

    return service.init()
      .then(() => expect(api.getEndpoints.calledOnce).to.be.true)
      .then(() => service.selectEndpoint('123'))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint.Value).to.equal('cat'));
  });

  it('adding a new endpoint stores it and sets it as the current endpoint', () => {
    const newEndpoint = { Name: '789', Value: 'sheep' };

    const someEndpoints = [
      {
        Name: '123',
        Value: 'cat',
      },
      {
        Name: '456',
        Value: 'dog',
      },
    ];

    const newEndpointCreator = () => newEndpoint;
    const service = new EndpointService(api, newEndpointCreator);

    const mergedEndpoints = _.concat([newEndpoint], someEndpoints);

    sandbox.stub(api, 'getEndpoints').returns(Promise.resolve(someEndpoints));
    sandbox.stub(api, 'updateEndpoints').returns(Promise.resolve(mergedEndpoints));


    return service.init()
      .then(() => service.selectEndpoint('123'))
      .then(() => service.addNewEndpoint())
      .then(() => expect(api.updateEndpoints.calledWith(JSON.stringify(mergedEndpoints))).to.be.true)
      .then(() => expect(service.endpoints).to.have.lengthOf(3))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint).to.eql(newEndpoint));
  });

  it('sets currently selected to null when deleting', () => {

    const someEndpoints = [
      {
        Name: '123',
        Value: 'cat',
      },
      {
        Name: '456',
        Value: 'dog',
      },
    ];

    const service = new EndpointService(api);

    const endpointsWithItemDeleted = [someEndpoints[0]];

    sandbox.stub(api, 'getEndpoints').returns(Promise.resolve(someEndpoints));
    sandbox.stub(api, 'updateEndpoints').returns(Promise.resolve(endpointsWithItemDeleted));


    return service.init()
      .then(() => service.selectEndpoint("456"))
      .then(() => service.deleteEndpoint())
      .then(() => expect(api.updateEndpoints.calledWith(JSON.stringify(endpointsWithItemDeleted))).to.be.true)
      .then(() => expect(service.endpoints).to.have.lengthOf(1))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint).to.be.an('undefined'));
  });
});
