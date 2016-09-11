import React from 'react';
import { isJSON } from '../util';

function Body({ value, label }) {
  let textClassName = 'mdl-color-text--primary';

  if (value && value !== '') {
    let text = value;
    if (isJSON(value)) {
      textClassName = textClassName + ' json';
      text = JSON.stringify(JSON.parse(value), null, 2);
    }

    return (
      <div>
        <div className="mdl-card__title mdl-card--expand">
          <h6 className="mdl-card__title-text">{label}</h6>
        </div>
        <div className="mdl-card__supporting-text">
          <pre className={textClassName}>{text}</pre>
        </div>
      </div>
    );
  }

  return null;
}

Body.propTypes = {
  value: React.PropTypes.string,
  label: React.PropTypes.string,
};

export default Body;
