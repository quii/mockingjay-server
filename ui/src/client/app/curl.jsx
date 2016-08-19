import React from 'react';
import request from 'superagent';

export const Curl = React.createClass({
  componentWillMount() {
    let baseURL = location.origin;
    if (this.props.baseURL) {
      baseURL = this.props.baseURL;
    }

    const url = `/mj-curl?name=${this.props.name}&baseURL=${baseURL}`;
    request
      .get(url)
      .end((err, res) => {
        if (err) {
          console.error(this.props.url, status, err.toString());
        } else {
          this.setState({ curl: res.text });
        }
      });
  },
  hint() {
    if (this.props.showPostHint) {
      return (
                <div className="mdl-card__supporting-text hint">
                    <code>curl</code>'s default behaviour with -d is to send <code>application/x-www-form-urlencoded</code> header if no other headers are specified. MJ will not match this request unless you specify form values <strong>or</strong> the http header in the request explicitly.
                </div>
                );
    } else {
      return null;
    }
  },
  render() {
    if (this.state) {
      return (
                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__supporting-text">
                        <code>{this.state.curl}</code>
                    </div>
                    {this.hint()}
                </div>
            );
    } else {
      return null;
    }
  },
});
