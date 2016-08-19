import React from 'react';
import { rand } from '../util';

const selectedCSS = 'mdl-button mdl-js-button mdl-button--raised mdl-button--accent';
const notSelectedCSS = 'mdl-button mdl-js-button mdl-button--raised mdl-button--colored';
const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'];

class MethodSwitcher extends React.Component {

  constructor(props) {
    super(props);

    this.handleClick = this.handleClick.bind(this);
    this.createButton = this.createButton.bind(this);
  }

  handleClick(e) {
    this.props.onChange({
      target: {
        name: 'method',
        value: e.target.innerText,
      },
    });
  }

  createButton(methodName, selectedMethod) {
    const clz = methodName === selectedMethod ? selectedCSS : notSelectedCSS;
    return (
      <button
        key={rand()}
        style={{ marginRight: '10px' }}
        className={clz}
        onClick={this.handleClick}
      >{methodName}</button>);
  }

  render() {
    const buttons = methods.map(m => this.createButton(m, this.props.selected));
    return <div className="method-switcher">{buttons}</div>;
  }
}

MethodSwitcher.propTypes = {
  onChange: React.PropTypes.func.isRequired,
  selected: React.PropTypes.string.isRequired,
};

export default MethodSwitcher;
