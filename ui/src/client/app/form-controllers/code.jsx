import React from 'react';

function Code({ value }) {
  if (value) {
    return (
      <div className="mdl-card__supporting-text">
        <code className="mdl-color-text--accent">{value}</code>
      </div>
    );
  }
  return null;
}

Code.propTypes = {
  value: React.PropTypes.string,
};

export default Code;
