jest.unmock('./../../../src/client/app/ui.jsx');
jest.unmock('./../../../src/client/app/navigation.jsx');

import UI from './../../../src/client/app/ui.jsx';
import Promise from 'bluebird';
import EndpointService from '../../../src/client/app/EndpointService'
import React from 'react';
import TestUtils from 'react-addons-test-utils';

describe('Endpoint serverice', () => {

  it('gets endpoints and then you can select them', () => {
    const service = new EndpointService();
    service.init.mockReturnValue(Promise.resolve(service));

    const endpoints = [
      { Name: 'endpoint 1'},
      { Name: 'endpoint 2'},
    ];

    service.getEndpoints.mockReturnValue(endpoints);

    const component = TestUtils.renderIntoDocument(<UI service={service}/>);

    expect(service.init).toBeCalled();
  });


});
