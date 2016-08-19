import React from 'react';
import TextField from '../form-controllers/textfield.jsx';

function MiscInfoEditor({ onChange, onCheckboxChange, name, cdcDisabled }) {
  return (
    <div className="mdl-card mdl-shadow--2dp">

      <div className="mdl-card__title" style={{ width: '90%' }}>
        <h3 className="mdl-card__title-text">Misc.</h3>
      </div>

      <TextField
        name="name"
        label="Endpoint name"
        value={name}
        onChange={onChange}
      />

      <label className="mdl-checkbox mdl-js-checkbox" htmlFor="cdcDisabled">
        <input
          type="checkbox"
          onClick={onCheckboxChange}
          name="cdcDisabled"
          className="mdl-checkbox__input"
          defaultChecked={cdcDisabled}
        />
        <span className="mdl-checkbox__label">
          <abbr title="Consumer driven contract">CDC</abbr> Disabled?
        </span>
      </label>
    </div>
  );
}

MiscInfoEditor.propTypes = {
  onChange: React.PropTypes.func.isRequired,
  onCheckboxChange: React.PropTypes.func.isRequired,
  name: React.PropTypes.string.isRequired,
  cdcDisabled: React.PropTypes.bool.isRequired,
};

export default MiscInfoEditor;
