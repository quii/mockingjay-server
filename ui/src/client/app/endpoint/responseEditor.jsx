import React from 'react';
import { HttpDataEditor } from '../form-controllers/httpDataList.jsx';
import TextField from '../form-controllers/textfield.jsx';
import TextArea from '../form-controllers/textarea.jsx';

function ResponseEditor({ code, body, onChange, resHeaders }) {
  return (
    <div className="mdl-card mdl-shadow--2dp">

      <div className="mdl-card__title" style={{ width: '90%' }}>
        <h3 className="mdl-card__title-text">Response</h3>
      </div>

      <TextField
        label="Status code"
        pattern="[0-9][0-9][0-9]"
        errMsg="Not valid HTTP status"
        name="code" value={code.toString()}
        onChange={onChange}
      />

      <TextArea
        label="Body"
        name="body"
        value={body}
        onChange={onChange}
      />

      <HttpDataEditor
        label="Headers"
        name="resHeaders"
        items={resHeaders}
        onChange={onChange}
      />

    </div>
  );
}

ResponseEditor.propTypes = {
  onChange: React.PropTypes.func.isRequired,
  code: React.PropTypes.string.isRequired,
  body: React.PropTypes.string,
  resHeaders: React.PropTypes.object.isRequired,
};

export default ResponseEditor;
