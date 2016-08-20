import Promise from 'bluebird';
import request from 'superagent';

function getJSONFromURL(url) {
  return new Promise((resolve, reject) => {
    request
      .get(url)
      .end((err, res) => {
        if (err) {
          reject(err);
        } else {
          resolve(res.text);
        }
      });
  }).then(body => JSON.parse(body));
}

function putToURL(url, data) {
  return new Promise((resolve, reject) => {
    request
      .put(url)
      .send(data)
      .end((err, res) => {
        if (err) {
          const error = `Got ${res.statusText} when sending update to ${url}. Response: ${res.text}`;
          reject(new Error(error));
        } else {
          resolve(res.text);
        }
      });
  }).then(body => JSON.parse(body));
}

class API {

  constructor(baseURL) {
    this.baseURL = baseURL;
    this.getEndpoints = this.getEndpoints.bind(this);
    this.updateEndpoints = this.updateEndpoints.bind(this);
  }

  getEndpoints() {
    return getJSONFromURL(this.baseURL);
  }

  checkCompatability(remoteURL) {
    return getJSONFromURL(`${this.baseURL}?url=${remoteURL}`);
  }

  updateEndpoints(data) {
    return putToURL(this.baseURL, data);
  }

}

export default API;
