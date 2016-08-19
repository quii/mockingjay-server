import React from 'react';
import { isJSON } from '../util';

function Body({ value, label }) {
  let text = value || '';
  if (isJSON(value)) {
    text = JSON.stringify(JSON.parse(value), null, 2);
  }

  return (
    <div>
      <div className="mdl-card__title mdl-card--expand">
        <h6 className="mdl-card__title-text">{label}</h6>
      </div>
      <div className="mdl-card__supporting-text">
        <pre className="mdl-color-text--primary">{text}</pre>
      </div>
    </div>
  );
}

Body.propTypes = {
  value: React.PropTypes.string,
  label: React.PropTypes.string,
};

export default Body;
