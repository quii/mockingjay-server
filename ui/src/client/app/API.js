import Promise from 'bluebird';
import request from 'superagent';

class API {

  constructor(baseURL) {
    this.baseURL = baseURL;
    this.getEndpoints = this.getEndpoints.bind(this);
    this.updateEndpoints = this.updateEndpoints.bind(this);
  }

  getEndpoints() {
    console.log('api being used');
    return new Promise((resolve, reject) => {
      request
        .get(this.baseURL)
        .end((err, res) => {
          if (err) {
            reject(err);
          } else {
            resolve(res.text);
          }
        });
    }).then(body => JSON.parse(body));
  }

  updateEndpoints(data) {
    console.log('api being for put');

    return new Promise((resolve, reject) => {
      request
        .put(this.baseURL)
        .send(data)
        .end((err, res) => {
          if (err) {
            reject(err);
          } else {
            resolve(res.text);
          }
        });
    }).then(body => JSON.parse(body));
  }

}

export default API;
