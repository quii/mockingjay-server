jest.unmock('lodash');
jest.unmock('bluebird');
jest.unmock('../../../src/client/app/EndpointService');

import EndpointService from '../../../src/client/app/EndpointService';
import _ from 'lodash';
import Promise from 'bluebird';
import API from '../../../src/client/app/API'

describe('Endpoint serverice', () => {


  it('gets endpoints and then you can select them', () => {
    const api = new API();
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

    api.getEndpoints.mockReturnValueOnce(Promise.resolve(someEndpoints));

    return service.init()
      .then(() => expect(api.getEndpoints).toBeCalled())
      .then(() => service.selectEndpoint('123'))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint.Value).toBe('cat'));
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

    const api = new API();
    const newEndpointCreator = () => newEndpoint;
    const service = new EndpointService(api, newEndpointCreator);

    const mergedEndpoints = _.concat([newEndpoint], someEndpoints);

    api.getEndpoints.mockReturnValueOnce(Promise.resolve(someEndpoints));
    api.updateEndpoints.mockReturnValueOnce(Promise.resolve(mergedEndpoints));


    return service.init()
      .then(() => service.selectEndpoint('123'))
      .then(() => service.addNewEndpoint())
      .then(() => expect(api.updateEndpoints).toBeCalledWith(JSON.stringify(mergedEndpoints)))
      .then(() => expect(service.endpoints.length).toBe(3))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint).toBe(newEndpoint));
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

    const api = new API();
    const service = new EndpointService(api);

    const endpointsWithItemDeleted = [someEndpoints[0]];

    api.getEndpoints.mockReturnValueOnce(Promise.resolve(someEndpoints));
    api.updateEndpoints.mockReturnValueOnce(Promise.resolve(endpointsWithItemDeleted));


    return service.init()
      .then(() => service.selectEndpoint("456"))
      .then(() => service.deleteEndpoint())
      .then(() => expect(api.updateEndpoints).toBeCalledWith(JSON.stringify(endpointsWithItemDeleted)))
      .then(() => expect(service.endpoints.length).toBe(1))
      .then(() => service.getEndpoint())
      .then(endpoint => expect(endpoint).toBe(undefined));
  });
});
