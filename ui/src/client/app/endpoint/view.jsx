import React from 'react';

// import Curl from './curl.jsx';
import Body from '../form-controllers/body.jsx';
import Code from '../form-controllers/code.jsx';
import { HttpDataList } from '../form-controllers/httpDataList.jsx';

// function couldBeDodgyCurlFormStuff(reqHeaders, reqBody) {
//   const noHeaders = !reqHeaders || Object.keys(reqHeaders).length === 0;
//   return reqBody !== '' && noHeaders;
// }

function View({
  method,
  uri,
  regex,
  reqBody,
  reqHeaders,
  form,
  code,
  body,
  resHeaders,
  // name,
  startEditing,
  deleteEndpoint,
}) {
  return (
    <div>

      <div className="mdl-card mdl-shadow--2dp">
        <div className="mdl-card__title" style={{ width: '90%' }}>
          <h3 className="mdl-card__title-text">Request</h3>
        </div>
        <Code icon="cloud" value={`${method} ${uri}`} />
        <Code icon="face" value={regex} />
        <Body label="Body" value={reqBody} />
        <HttpDataList name="Headers" items={reqHeaders} />
        <HttpDataList name="Form data" items={form} />
      </div>

      <div className="mdl-card mdl-shadow--2dp">
        <div className="mdl-card__title" style={{ width: '90%' }}>
          <h3 className="mdl-card__title-text">Response</h3>
        </div>
        <Code icon="face" value={code.toString()} />
        <Body label="Body" value={body} />
        <HttpDataList name="Headers" items={resHeaders} />
      </div>

      {/* <Curl*/}
      {/* baseURL={location.origin}*/}
      {/* name={name}*/}
      {/* showPostHint={couldBeDodgyCurlFormStuff(reqHeaders, reqBody)}*/}
      {/* />*/}

      <div style={{ margin: '2% 2% 2% 3%' }}>
        <button
          style={{ margin: '0% 1% 0% 0%' }}
          onClick={startEditing}
          className="mdl-button mdl-button--raised mdl-button--accent"
        >Edit
        </button>

        <button
          onClick={deleteEndpoint}
          className="mdl-button mdl-button--raised mdl-button--primary"
        >Delete
        </button>
      </div>
    </div>
  );
}

View.propTypes = {
  method: React.PropTypes.string.isRequired,
  uri: React.PropTypes.string.isRequired,
  regex: React.PropTypes.string,
  reqBody: React.PropTypes.string,
  reqHeaders: React.PropTypes.object,
  form: React.PropTypes.object,
  code: React.PropTypes.string.isRequired,
  body: React.PropTypes.string,
  resHeaders: React.PropTypes.object,
  name: React.PropTypes.string.isRequired,
  startEditing: React.PropTypes.func.isRequired,
  deleteEndpoint: React.PropTypes.func.isRequired,
};

export default View;
