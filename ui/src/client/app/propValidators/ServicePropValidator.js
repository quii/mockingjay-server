import _ from 'lodash';

const requiredFunctions = [
  'init',
  'getEndpoint',
  'getEndpoints',
  'addNewEndpoint',
  'updateEndpoint',
  'selectEndpoint',
  'deleteEndpoint',
];

function ServiceProp(props, propName, componentName) {

  const functionsMissing = _.filter(requiredFunctions, (f) => props[propName][f] === undefined);

  if (functionsMissing.length > 0) {
    return new Error(
      `Invalid prop ${propName} for ${componentName} - missing functions [${functionsMissing}]`
    );
  }
}

export default ServiceProp;
