import React from 'react';
import dialogPolyfill from 'dialog-polyfill';
import {rand, isValidURL} from './util';

const CDC = React.createClass({
  getInitialState: function () {
    return {
      remoteUrl: location.origin
    };
  },
  componentWillMount: function () {
    this.checkCompatability();
  },
  checkCompatability: function () {
    if (this.state && this.state.remoteUrl && this.state.remoteUrl !== null) {
      $.ajax({
        url: this.props.url + '?url=' + this.state.remoteUrl,
        dataType: 'json',
        cache: false,
        success: function (data) {
          this.setState({data: data});
        }.bind(this),
        error: function (xhr, status, err) {
          console.error(this.props.url, status, err.toString());
        }.bind(this)
      });
    }
  },
  handleUrlChange: function (e) {
    this.setState({
      data: null
    });

    if (isValidURL(e.target.value)) {
      this.setState({
        remoteUrl: e.target.value
      }, this.checkCompatability);
    }
  },
  inputWidth: function () {
    let w = this.state.remoteUrl.length * 12;
    return w < 350 ? 350 : w;
  },
  sentiment: function () {
    let sentiment;
    if (this.state && this.state.data) {
      sentiment = this.state.data.Passed ? 'sentiment_satisfied' : 'sentiment_very_dissatisfied';
    } else {
      sentiment = 'sentiment_neutral';
    }
    return sentiment;
  },
  label: 'Auto-checking endpoints are equivalent to',
  indicatorClick: function () {
    this.refs['dialog'].showModal();
  },
  render: function () {
    return (
            <header className="mdl-layout__header">
                <div className="mdl-layout__header-row">
                    <div className="cdc mdl-textfield mdl-js-textfield mdl-textfield--expandable
                  mdl-textfield--floating-label mdl-textfield--align-right">
                        <TestIndicator indicatorClick={this.indicatorClick} badge={this.sentiment()}/>

                        <div className="mdl-textfield__expandable-holder">
                            <input style={{width: this.inputWidth()}} className="mdl-textfield__input" type="text"
                                   name="sample"
                                   id="fixed-header-drawer-exp" onBlur={this.checkCompatability}
                                   onKeyPress={this.handleUrlChange} defaultValue={this.state.remoteUrl}/>
                        </div>

                        <label htmlFor="fixed-header-drawer-exp">{this.label}</label>

                    </div>
                </div>
                <Dialog
                    title="Messages from CDC check"
                    messages={this.state && this.state.data ? this.state.data.Messages : []}
                    ref="dialog"
                />
            </header>
        );
  }
});

const TestIndicator = React.createClass({
  render: function () {
    return (<i onClick={this.props.indicatorClick} style={{cursor: 'hand'}}
                  className="material-icons md-48">{this.props.badge}</i>);
  }
});


const Dialog = React.createClass({
  componentDidMount: function () {
    const dialog = document.querySelector('dialog');

    if (!dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }
  },
  showModal: function () {
    if (this.props.messages && this.props.messages.length > 0) {
      document.querySelector('dialog').showModal();
    }
  },
  close: function () {
    document.querySelector('dialog').close();
  },
  render: function () {
    let errs;
    if (this.props.messages) {
      errs = this.props.messages.map(m => <li key={rand()}>{m}</li>);
    }
    return (
            <dialog className="mdl-dialog" style={{width: '700px'}}>
                <h4 className="mdl-dialog__title">{this.props.title}</h4>
                <div className="mdl-dialog__content" style={{'fontFamily': 'Courier'}}>
                    <ul>{errs}</ul>
                </div>
                <div className="mdl-dialog__actions">
                    <button type="button" onClick={this.close} className="mdl-button">Close</button>
                </div>
            </dialog>
        );
  }
});


export default CDC;
