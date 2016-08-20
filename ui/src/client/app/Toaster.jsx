import React from 'react';

class Toaster extends React.Component {

  alert(msgs) {
    const message = [].concat(msgs || []).join(', ');
    const notification = document.querySelector('.mdl-js-snackbar');
    notification.MaterialSnackbar.showSnackbar(
      {
        message: message,
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

export default Toaster;
