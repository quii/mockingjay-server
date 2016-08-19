import React from 'react';
import dialogPolyfill from 'dialog-polyfill';
import { rand, isValidURL } from './util';

class CDC extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      remoteUrl: location.origin,
    };

    this.setDialog = this.setDialog.bind(this);
    this.handleUrlChange = this.handleUrlChange.bind(this);
    this.indicatorClick = this.indicatorClick.bind(this);
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
      $.ajax({
        url: `${this.props.url}?url=${this.state.remoteUrl}`,
        dataType: 'json',
        cache: false,
        success: function (data) {
          this.setState({data});
        }.bind(this),
        error: function (xhr, status, err) {
          console.error(this.props.url, status, err.toString());
        }.bind(this),
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

CDC.label = 'Auto-checking endpoints are equivalent to';

const TestIndicator = React.createClass({
  propTypes: {
    badge: React.PropTypes.string.isRequired,
    indicatorClick: React.PropTypes.func.isRequired,
  },
  render() {
    return (
      <i
        onClick={this.props.indicatorClick} style={{cursor: 'hand'}}
        className="material-icons md-48"
      >{this.props.badge}</i>);
  },
});


const Dialog = React.createClass({
  propTypes: {
    messages: React.PropTypes.array.isRequired,
    title: React.PropTypes.string.isRequired,
  },
  componentDidMount() {
    const dialog = document.querySelector('dialog');

    if (!dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }
  },
  showModal() {
    if (this.props.messages && this.props.messages.length > 0) {
      document.querySelector('dialog').showModal();
    }
  },
  close() {
    document.querySelector('dialog').close();
  },
  render() {
    let errs;
    if (this.props.messages) {
      errs = this.props.messages.map(m => <li key={rand()}>{m}</li>);
    }
    return (
      <dialog className="mdl-dialog" style={{ width: '700px' }}>
        <h4 className="mdl-dialog__title">{this.props.title}</h4>
        <div className="mdl-dialog__content" style={{fontFamily: 'Courier'}}>
          <ul>{errs}</ul>
        </div>
        <div className="mdl-dialog__actions">
          <button type="button" onClick={this.close} className="mdl-button">Close</button>
        </div>
      </dialog>
    );
  },
});


export default CDC;
