import React from 'react';
import ReactDOM from 'react-dom';
import UI from './ui.jsx';
import API from './API.js';
import EndpointService from './EndpointService';

const service = new EndpointService(new API('/mj-endpoints'));

ReactDOM.render(
  <UI service={service} />,
  document.getElementById('app')
);
