jest.unmock('./../../../src/client/app/ui.jsx');
jest.unmock('./../../../src/client/app/navigation.jsx');

import UI from './../../../src/client/app/ui.jsx';
import { shallow } from 'enzyme';
import Navigation from './../../../src/client/app/navigation.jsx';
import Promise from 'bluebird';
import EndpointService from '../../../src/client/app/EndpointService'
import React from 'react';

describe('Endpoint serverice', () => {

  it('gets endpoints and then you can select them', () => {
    const service = new EndpointService();

    const endpoints = [
      { Name: 'endpoint 1'},
      { Name: 'endpoint 2'},
    ];

    service.init.mockReturnValue(Promise.resolve(service));
    service.getEndpoints.mockReturnValue(endpoints);

    const ui = shallow(<UI service={service} />);

    expect(service.getEndpoints).toBeCalled();
    expect(ui.find(Navigation).length).toEqual(1);

  });


});
