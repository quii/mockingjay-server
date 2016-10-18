import React from 'react';
import _ from 'lodash';
import requestEditor from './requestEditor.jsx';
import responeEditor from './responseEditor.jsx';
import miscEditor from './miscInfoEditor.jsx';


class Form extends React.Component {
  componentDidMount() {
    componentHandler.upgradeDom();
  }

  render() {
    const ultraProps = _.cloneDeep(this.props.originalValues);
    ultraProps.onChange = this.props.onChange;
    ultraProps.onCheckboxChange = this.props.onCheckboxChange;
    ultraProps.name = this.props.name;

    return (
      <div>
        {requestEditor(ultraProps)}
        {responeEditor(ultraProps)}
        {miscEditor(ultraProps)}

        <div style={{ margin: '2% 2% 2% 3%' }}>
          <button
            onClick={this.props.finishEditing}
            className="mdl-button mdl-js-button mdl-button--raised mdl-button--accent"
          >
            Save
          </button>
        </div>

      </div>
    );
  }
}

Form.propTypes = {
  onChange: React.PropTypes.func.isRequired,
  onCheckboxChange: React.PropTypes.func.isRequired,
  finishEditing: React.PropTypes.func.isRequired,
  originalValues: React.PropTypes.object.isRequired,
  name: React.PropTypes.string.isRequired,
};

export default Form;
