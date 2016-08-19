import React from 'react';
import { HttpDataEditor } from '../form-controllers/httpDataList.jsx';
import TextField from '../form-controllers/textfield.jsx';
import TextArea from '../form-controllers/textarea.jsx';
import MethodSwitcher from '../form-controllers/methodSwitcher.jsx';

function RequestEditor({onChange, uri, regex, method,reqBody, form, reqHeaders}){
  return (
    <div className="mdl-card mdl-shadow--2dp">

      <div className="mdl-card__title" style={{ width: '90%' }}>
        <h3 className="mdl-card__title-text">Request</h3>
      </div>

      <TextField label="URI" name="uri" value={uri} onChange={onChange} />
      <TextField label="Regex URI (optional)" name="regex" value={regex} onChange={onChange} />
      <MethodSwitcher selected={method} onChange={onChange}/>
      <TextArea label="Body" name="reqBody" value={reqBody} onChange={onChange}/>
      <HttpDataEditor label="x-www-form-urlencoded" name="form" items={form} onChange={onChange}/>

      <HttpDataEditor
        label="Headers"
        keyPattern="[A-Za-z0-9\S]{1,25}"
        valPattern="[A-Za-z0-9\S]{1,25}"
        name="reqHeaders"
        items={reqHeaders}
        onChange={onChange}
      />
    </div>
  );
}

RequestEditor.propTypes = {
  onChange: React.PropTypes.func.isRequired,
  uri: React.PropTypes.string.isRequired,
  regex: React.PropTypes.string,
  method: React.PropTypes.string.isRequired,
  reqBody: React.PropTypes.string.isRequired,
  form: React.PropTypes.object.isRequired,
  reqHeaders: React.PropTypes.object.isRequired,
};

export default RequestEditor;
