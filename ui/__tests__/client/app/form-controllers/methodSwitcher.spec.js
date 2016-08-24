jest.unmock('../../../../src/client/app/form-controllers/methodSwitcher.jsx');
jest.unmock('react');
jest.unmock('react-dom');
jest.unmock('react-addons-test-utils');
import MethodSwitcher from '../../../../src/client/app/form-controllers/methodSwitcher.jsx';

import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

describe('Method switcher', () => {
  it('changes the highlight as users click', () => {

    const onChange = jest.fn();

    const methodSwitcher = TestUtils.renderIntoDocument(
      <MethodSwitcher selected="POST" onChange={onChange} />
    );

    const methodSwitcherNode = ReactDOM.findDOMNode(methodSwitcher);

    console.log('method switcher node', methodSwitcherNode);

    expect(1).toEqual(1);

  });
});
