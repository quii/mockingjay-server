import assert from 'assert'
import sinon from 'sinon';
import EndpointService from '../../../src/client/app/EndpointService'

const api = {
  getEndpoints: () => {}
}

describe('Endpoint serverice', () => {

  let sandbox;

  beforeEach(() => {
    sandbox = sinon.sandbox.create();
  });

  afterEach(() => {
    sandbox.restore();
  });

  it('calls the api to get endpoints', () => {

    const someEndpoints = [
      {
        Name: "123"
      },
      {
        Name: "456"
      },
    ];

    sandbox.stub(api, 'getEndpoints').returns(Promise.resolve(someEndpoints));
    const service = new EndpointService(api);

    return service.init()
      .then(() => assert.equal(true, api.getEndpoints.calledOnce));
  })
})
