import React from 'react';

function TestIndicator({ indicatorClick, badge }) {
  return (
    <i
      onClick={indicatorClick} style={{ cursor: 'hand' }}
      className="material-icons md-48"
    >{badge}</i>);
}

TestIndicator.propTypes = {
  badge: React.PropTypes.string.isRequired,
  indicatorClick: React.PropTypes.func.isRequired,
};


export default TestIndicator;
