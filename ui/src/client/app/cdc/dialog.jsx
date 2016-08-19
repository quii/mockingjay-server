import React from 'react';
import dialogPolyfill from 'dialog-polyfill';
import { rand } from '../util';

class Dialog extends React.Component {

  constructor(props) {
    super(props);
    this.showModal = this.showModal.bind(this);
  }

  componentDidMount() {
    const dialog = document.querySelector('dialog');

    if (!dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }
  }

  showModal() {
    if (this.props.messages && this.props.messages.length > 0) {
      document.querySelector('dialog').showModal();
    }
  }

  close() {
    document.querySelector('dialog').close();
  }

  render() {
    let errs;
    if (this.props.messages) {
      errs = this.props.messages.map(m => <li key={rand()}>{m}</li>);
    }
    return (
      <dialog className="mdl-dialog" style={{ width: '700px' }}>
        <h4 className="mdl-dialog__title">{this.props.title}</h4>
        <div className="mdl-dialog__content" style={{ fontFamily: 'Courier' }}>
          <ul>{errs}</ul>
        </div>
        <div className="mdl-dialog__actions">
          <button type="button" onClick={this.close} className="mdl-button">Close</button>
        </div>
      </dialog>
    );
  }
}

Dialog.propTypes = {
  messages: React.PropTypes.array.isRequired,
  title: React.PropTypes.string.isRequired,
};

export default Dialog;
