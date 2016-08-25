jest.unmock('../../../../src/client/app/form-controllers/methodSwitcher.jsx');
import MethodSwitcher from '../../../../src/client/app/form-controllers/methodSwitcher.jsx';

import React from 'react';
import ReactDOM from 'react-dom';
import TestUtils from 'react-addons-test-utils';

describe('Method switcher', () => {

  let onChange, component, renderedDOM;

  beforeEach(() => {
    onChange = jest.fn();
    component = TestUtils.renderIntoDocument(<MethodSwitcher selected="POST" onChange={onChange} />);
  });

  it('changes the highlight as users click', () => {

    let renderedButtons = TestUtils.scryRenderedDOMComponentsWithTag(component, "button");
    expect(renderedButtons.length).toEqual(MethodSwitcher.methods.length);
    
  });
});
