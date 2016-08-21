import React from 'react';
import { guid } from './util';

function makeAddButton(onClick) {
  return (
    <a
      key={guid()}
      style={{ paddingLeft: '30px' }}
      onClick={onClick}
      className="mdl-navigation__link mdl-color-text--primary-dark"
    ><i className="material-icons md-32">add</i>Add new endoint</a>
  );
}

function Navigation({ addEndpoint, endpoints, activeEndpoint, openEditor }) {
  const endpointLinks = endpoints.map(endpoint => {
    let cssClass = 'mdl-navigation__link';
    let name = endpoint.Name;

    if (endpoint.Name === activeEndpoint) {
      cssClass += ' mdl-color--primary-contrast mdl-color-text--accent';
      name = <div><i className="material-icons md-32">fingerprint</i>{endpoint.Name}</div>;
    }

    return (
      <a
        key={guid()}
        className={cssClass}
        onClick={(event) => openEditor(endpoint.Name, event)}
      >{name}</a>
    );
  });

  endpointLinks.push(makeAddButton(addEndpoint));

  return <nav className="mdl-navigation">{endpointLinks}</nav>;
}

Navigation.propTypes = {
  addEndpoint: React.PropTypes.func.isRequired,
  openEditor: React.PropTypes.func.isRequired,
  endpoints: React.PropTypes.arrayOf(React.PropTypes.object).isRequired,
  activeEndpoint: React.PropTypes.string,
};

export default Navigation;
