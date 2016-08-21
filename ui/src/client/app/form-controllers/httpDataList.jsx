import React from 'react';
import { rand } from '../util';
import _ from 'lodash';

function mapKeyVals(items, f) {
  if (items) {
    let i = -1;
    return Object.keys(items).map(key => {
      i++;
      const value = items[key];
      return (
        f(key, value, i)
      );
    });
  }
  return [];
}

export function HttpDataList({ items, label, name }) {
  const itemRows = mapKeyVals(items, (key, val) => <tr key={rand()}>
    <td className="mdl-data-table__cell--non-numeric">{key}</td>
    <td className="mdl-data-table__cell--non-numeric">{val}</td>
  </tr>);

  if (itemRows.length > 0) {
    return (<div>
      <div className="mdl-card__title mdl-card--expand">
        <h6 className="mdl-card__title-text">{label || name}</h6>
      </div>
      <table className="mdl-data-table mdl-js-data-table mdl-data-table">
        <thead>
        <tr>
          <th className="mdl-data-table__cell--non-numeric">Name</th>
          <th className="mdl-data-table__cell--non-numeric">Value</th>
        </tr>
        </thead>
        <tbody>
        {itemRows}
        </tbody>
      </table>
    </div>);
  }

  return null;
}

HttpDataList.propTypes = {
  label: React.PropTypes.string,
  name: React.PropTypes.string.isRequired,
  items: React.PropTypes.object,
};

export class HttpDataEditor extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      numberOfItems: this.props.items ? Object.keys(this.props.items).length : 0,
    };

    this.addItem = this.addItem.bind(this);
    this.updateMap = this.updateMap.bind(this);
    this.createInput = this.createInput.bind(this);
    this.render = this.render.bind(this);
  }

  addItem() {
    this.setState({
      numberOfItems: this.state.numberOfItems + 1,
    });
  }

  updateMap() {
    const newState = {};
    for (let i = 0; i < Object.keys(this.refs).length; i += 2) {
      const keyName = Object.keys(this.refs)[i];
      const valueName = Object.keys(this.refs)[i + 1];


      const k = this.refs[keyName].value;
      const v = this.refs[valueName].value;

      if (k !== '' || v !== '') {
        newState[k] = v;
      }
    }

    const change = _.isEmpty(newState) ? null : newState;

    this.props.onChange({
      target: {
        name: this.props.name,
        value: change,
      },
    });
  }

  createInput(ref, key, val) {
    const keyRef = `${this.props.name}${ref}key`;
    const valRef = `${this.props.name}${ref}val`;

    return (
      <div className="mdl-textfield mdl-js-textfield">
        <input
          ref={keyRef}
          pattern={this.props.keyPattern}
          className="mdl-textfield__input"
          type="text"
          value={key}
          onChange={this.updateMap}
        />
        <i className="material-icons">chevron_right</i>
        <input
          ref={valRef}
          pattern={this.props.valPattern}
          className="mdl-textfield__input"
          type="text"
          value={val}
          onChange={this.updateMap}
        />
      </div>
    );
  }

  render() {
    const label = this.props.label || this.props.name;
    const items = mapKeyVals(this.props.items, (key, val, i) => this.createInput(i, key, val));
    items.push(this.createInput(items.length + 1, '', ''));

    const remainingItems = (this.state.numberOfItems + 1) - items.length;

    for (let i = 0; i < remainingItems; i++) {
      const newItem = this.createInput(i + items.length, '', '');
      items.push(newItem);
    }

    return (
      <div className="list-editor">
        <h4>{label}</h4>
        <ul>{items}</ul>
      </div>);
  }
}

HttpDataEditor.propTypes = {
  label: React.PropTypes.string,
  name: React.PropTypes.string.isRequired,
  items: React.PropTypes.object,
  onChange: React.PropTypes.func.isRequired,
  servce: React.PropTypes.object.isRequired,
  keyPattern: React.PropTypes.string,
  valPattern: React.PropTypes.string,
};

