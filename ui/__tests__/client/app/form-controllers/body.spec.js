import React from 'react';
import { shallow } from 'enzyme';
import Body from '../../../../src/client/app/form-controllers/body.jsx';

jest.unmock('../../../../src/client/app/form-controllers/body.jsx');
jest.unmock('../../../../src/client/app/util.js');

describe('Body renderer', () => {
  it('renders json bodies', () => {
    const body = '{"foo": "bar"}';
    const component = shallow(<Body label={"whatever"} value={body} />);
    expect(component.find('.json').length).toEqual(1);
  });

  it('doesnt render json bodies when it isnt json', () => {
    const body = 'how now, brown cow';
    const component = shallow(<Body label={"whatever"} value={body} />);
    expect(component.find('.json').length).toEqual(0);
  });
});

