import React from 'react';
import request from 'superagent';
import { isValidURL } from '../util';
import TestIndicator from './testIndicator.jsx';
import Dialog from './dialog.jsx';

class CDC extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      remoteUrl: location.origin,
    };

    this.label = 'Auto-checking endpoints are equivalent to';
    this.setDialog = this.setDialog.bind(this);
    this.handleUrlChange = this.handleUrlChange.bind(this);
    this.indicatorClick = this.indicatorClick.bind(this);
    this.checkCompatability = this.checkCompatability.bind(this);
  }

  componentWillMount() {
    this.checkCompatability();
  }

  setDialog(ref) {
    this.dialog = ref;
  }

  sentiment() {
    let sentiment;
    if (this.state && this.state.data) {
      sentiment = this.state.data.Passed ? 'sentiment_satisfied' : 'sentiment_very_dissatisfied';
    } else {
      sentiment = 'sentiment_neutral';
    }
    return sentiment;
  }

  indicatorClick() {
    this.dialog.showModal();
  }

  inputWidth() {
    const w = this.state.remoteUrl.length * 12;
    return w < 350 ? 350 : w;
  }

  handleUrlChange(e) {
    this.setState({
      data: null,
    });

    if (isValidURL(e.target.value)) {
      this.setState({
        remoteUrl: e.target.value,
      }, this.checkCompatability);
    }
  }

  checkCompatability() {
    if (this.state && this.state.remoteUrl && this.state.remoteUrl !== null) {
      request
        .get(`${this.props.url}?url=${this.state.remoteUrl}`)
        .end((err, res) => {
          if (err) {
            console.error(this.props.url, status, err.toString());
          } else {
            this.setState({ data: JSON.parse(res.text) });
          }
        });
    }
  }

  render() {
    const input = (<div className="mdl-textfield__expandable-holder">
      <input
        style={{ width: this.inputWidth() }}
        className="mdl-textfield__input" type="text"
        name="sample"
        id="fixed-header-drawer-exp" onBlur={this.checkCompatability}
        onKeyPress={this.handleUrlChange} defaultValue={this.state.remoteUrl}
      />
    </div>);

    return (
      <header className="mdl-layout__header">
        <div className="mdl-layout__header-row">
          <div
            className="cdc mdl-textfield mdl-js-textfield mdl-textfield--expandable
              mdl-textfield--floating-label mdl-textfield--align-right"
          >
            <TestIndicator indicatorClick={this.indicatorClick} badge={this.sentiment()} />
            {input}
            <label htmlFor="fixed-header-drawer-exp">{this.label}</label>

          </div>
        </div>
        <Dialog
          title="Messages from CDC check"
          messages={this.state && this.state.data ? this.state.data.Messages : []}
          ref={this.setDialog}
        />
      </header>
    );
  }
}

CDC.propTypes = {
  url: React.PropTypes.string.isRequired,
  indicatorClick: React.PropTypes.func,
};


export default CDC;
