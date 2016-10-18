import React from 'react';

class Toaster extends React.Component {

  alert(msgs, timeout = 2750) {
    const message = [].concat(msgs || []).join(', ');
    const notification = document.querySelector('.mdl-js-snackbar');
    notification.MaterialSnackbar.showSnackbar(
      {
        message,
        timeout,
      }
    );
  }

  render() {
    return (
      <div
        aria-live="assertive"
        aria-atomic="true"
        aria-relevant="text"
        className="mdl-snackbar mdl-js-snackbar"
      >
        <div className="mdl-snackbar__text" />
        <button type="button" className="mdl-snackbar__action" />
      </div>
    );
  }
}

Toaster.ErrorDisplayTime = 5000;

export default Toaster;
