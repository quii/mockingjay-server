jest.unmock('../../../../src/client/app/form-controllers/methodSwitcher.jsx');
import MethodSwitcher from '../../../../src/client/app/form-controllers/methodSwitcher.jsx';

import _ from 'lodash';
import React from 'react';
import TestUtils from 'react-addons-test-utils';

describe('Method switcher', () => {

    it('highlights selected method and sends correct data on clicks', () => {

        const initiallySelectedMethod = "POST";
        const onChange = jest.fn();
        const component = TestUtils.renderIntoDocument(<MethodSwitcher selected={initiallySelectedMethod}
                                                                       onChange={onChange}/>);
        const buttons = TestUtils.scryRenderedDOMComponentsWithTag(component, "button");
        const buttonNames = buttons.map(b => b.textContent);

        expect(buttonNames).toEqual(MethodSwitcher.methods);
        expect(getSelectedMethod(component)).toEqual(initiallySelectedMethod);

        const putButton = findButtonWithMethod(buttons, 'PUT');
        TestUtils.Simulate.click(putButton);

        expect(onChange).toBeCalledWith(
            {
                target: {
                    name: 'method',
                    value: 'PUT',
                }
            });
    });
});

function findButtonWithMethod(buttons, method) {
    return _.find(buttons, b => b.textContent === method)
}

function getSelectedMethod(component) {
    return TestUtils.findRenderedDOMComponentWithClass(component, MethodSwitcher.selectedCSS).textContent;
}
