import React from 'react';

const errSpan = (msg) => (msg ? <span className="mdl-textfield__error">{msg}</span> : null);

function TextField({ value, label, errMsg, pattern, name, onChange }) {
  const defaultedValue = value || '';
  const defaultedLabel = label || name;
  return (
    <div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
      <input
        pattern={pattern}
        className="mdl-textfield__input"
        type="text" name={name}
        value={defaultedValue}
        onChange={onChange}
      />
      <label className="mdl-textfield__label" htmlFor={name}>{defaultedLabel}</label>
      {errSpan(errMsg)}
    </div>
  );
}

TextField.propTypes = {
  value: React.PropTypes.string,
  label: React.PropTypes.string,
  errMsg: React.PropTypes.string,
  pattern: React.PropTypes.string,
  name: React.PropTypes.string.isRequired,
  onChange: React.PropTypes.func.isRequired,
};

export default TextField;
