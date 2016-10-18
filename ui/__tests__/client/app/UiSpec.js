import Promise from 'bluebird';
import { shallow } from 'enzyme';
import React from 'react';
import UI from './../../../src/client/app/ui.jsx';
import Endpoint from './../../../src/client/app/endpoint/endpoint.jsx';
import Navigation from './../../../src/client/app/navigation.jsx';
import EndpointService from '../../../src/client/app/EndpointService';

jest.unmock('./../../../src/client/app/ui.jsx');
jest.unmock('./../../../src/client/app/navigation.jsx');

describe('Endpoint service', () => {
  const service = new EndpointService();

  const endpoints = [
    { Name: 'endpoint 1' },
    { Name: 'endpoint 2' },
  ];

  it('gets endpoints and renders a navigation', () => {
    service.init.mockReturnValue(Promise.resolve(service));
    service.getEndpoints.mockReturnValue(endpoints);

    const ui = shallow(<UI service={service} />);

    expect(ui.find(Navigation).length).toEqual(1);
    expect(ui.find(Navigation).prop('endpoints')).toEqual(endpoints);
  });

  it('renders an endpoint when service returns a currently selected endpoint', () => {
    const selectedEndpoint = endpoints[0];
    service.getEndpoint.mockReturnValue(selectedEndpoint);

    const ui = shallow(<UI service={service} />);

    expect(ui.find(Endpoint).length).toEqual(1);
    expect(ui.find(Endpoint).prop('endpoint')).toEqual(selectedEndpoint);
  });
});
