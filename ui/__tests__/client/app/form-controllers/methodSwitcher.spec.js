jest.unmock('../../../../src/client/app/form-controllers/methodSwitcher.jsx');
import MethodSwitcher from '../../../../src/client/app/form-controllers/methodSwitcher.jsx';

import React from 'react';
import TestUtils from 'react-addons-test-utils';

describe('Method switcher', () => {

  it('changes the highlight as users click', () => {

    const onChange = jest.fn();
    const component = TestUtils.renderIntoDocument(<MethodSwitcher selected="POST" onChange={onChange} />);

    const renderedButtons = TestUtils.scryRenderedDOMComponentsWithTag(component, "button");
    const buttonNames = renderedButtons.map(b => b.textContent);

    expect(buttonNames).toEqual(MethodSwitcher.methods);

    const selectedButton = TestUtils.findRenderedDOMComponentWithClass(component, MethodSwitcher.selectedCSS);

    expect(selectedButton.textContent).toEqual("POST")
  });
});
