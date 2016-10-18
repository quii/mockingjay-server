import React from 'react';

function TextArea({ label, name, value, onChange }) {
  const defaultLabel = label || name;
  return (
    <div
      style={{ width: '100%' }}
      className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label"
    >
      <textarea
        className="mdl-textfield__input"
        type="text"
        rows="5"
        name={name}
        value={value}
        onChange={onChange}
      />
      <label className="mdl-textfield__label" htmlFor={name}>{defaultLabel}</label>
    </div>
  );
}

TextArea.propTypes = {
  label: React.PropTypes.string,
  name: React.PropTypes.string.isRequired,
  value: React.PropTypes.string.isRequired,
  onChange: React.PropTypes.func.isRequired,
};

export default TextArea;
