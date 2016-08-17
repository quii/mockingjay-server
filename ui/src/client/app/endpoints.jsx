import React from 'react';
import {HttpDataList, HttpDataEditor} from './httpDataList.jsx';
import {Code, MethodSwitcher, Body, TextArea, TextField} from './formbits.jsx';

const Endpoint = React.createClass({
    getInitialState: function () {
        return {
            index: this.props.index,
            cdcDisabled: this.props.cdcDisabled,
            isEditing: false,
            name: this.props.name,
            method: this.props.method,
            uri: this.props.uri,
            regex: this.props.regex,
            reqBody: this.props.reqBody,
            form: this.props.form,
            reqHeaders: this.props.reqHeaders,
            code: this.props.code,
            body: this.props.body,
            resHeaders: this.props.resHeaders,
        };
    },
    startEditing: function () {
        this.setState({
            isEditing: true
        })
    },
    delete: function () {
        this.props.delete();
    },
    finishEditing: function () {
        this.setState({
            isEditing: false
        });
        this.props.updateServer();
    },
    updateValue: function (e) {
        this.setState({
            [e.target.name]: e.target.value
        })
    },
    updateCheckbox: function (e) {
        this.setState({
            [e.target.name]: e.target.value === 'on'
        })
    },
    render: function () {

        const view = (
            <div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Request</h3>
                    </div>
                    <Code icon="cloud" value={this.state.method + " " + this.state.uri}/>
                    <Code icon="face" value={this.state.regex}/>
                    <Body label="Body" value={this.state.reqBody}/>
                    <HttpDataList name="Headers" items={this.state.reqHeaders}/>
                    <HttpDataList name="Form data" items={this.state.form}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <Code icon="face" value={this.state.code}/>
                    <Body label="Body" value={this.state.body}/>
                    <HttpDataList name="Headers" items={this.state.resHeaders}/>
                </div>

                <div style={{margin:"2% 2% 2% 3%"}}>
                    <button style={{margin:"0% 1% 0% 0%"}} onClick={this.startEditing} className="mdl-button mdl-button--raised mdl-button--accent">
                        Edit
                    </button>
                    <button onClick={this.delete} className="mdl-button mdl-button--raised mdl-button--primary">
                        Delete
                    </button>
                </div>
            </div>);

        const form = <EndpointForm
            name={this.state.name}
            finishEditing={this.finishEditing}
            originalValues={this.state}
            onChange={this.updateValue}
            onCheckboxChange={this.updateCheckbox}
        />;

        return this.state.isEditing ? form : view;
    }
});

const EndpointForm = React.createClass({
    componentDidMount: function () {
        componentHandler.upgradeDom();
    },
    render: function () {
        return (
            <div>
                <div className="mdl-card mdl-shadow--2dp">
                    <TextField name="name" label="Endpoint name" value={this.props.originalValues.name} onChange={this.props.onChange} />
                    <label className="mdl-checkbox mdl-js-checkbox" htmlFor="cdcDisabled">
                        <input type="checkbox" onClick={this.props.onCheckboxChange} name="cdcDisabled" className="mdl-checkbox__input" defaultChecked={this.props.originalValues.cdcDisabled} />
                            <span className="mdl-checkbox__label">CDC Disabled?</span>
                    </label>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Request</h3>
                    </div>

                    <TextField label="URI" name="uri" value={this.props.originalValues.uri} onChange={this.props.onChange} />
                    <TextField label="Regex URI (optional)" name="regex" value={this.props.originalValues.regex} onChange={this.props.onChange} />
                    <MethodSwitcher selected={this.props.originalValues.method} onChange={this.props.onChange} />
                    <TextArea label="Body" name="reqBody" value={this.props.originalValues.reqBody} onChange={this.props.onChange} />


                    <HttpDataEditor label="Form" name="form" items={this.props.originalValues.form}
                                                       onChange={this.props.onChange}/>
                    <HttpDataEditor label="Headers" name="reqHeaders" items={this.props.originalValues.reqHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <TextField label="Status code" name="code" value={this.props.originalValues.code} onChange={this.props.onChange} />

                    <TextArea label="Body" name="body" value={this.props.originalValues.body} onChange={this.props.onChange} />
                    <HttpDataEditor label="Headers" name="resHeaders" items={this.props.originalValues.resHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <div style={{margin:"2% 2% 2% 3%"}}>
                    <button onClick={this.props.finishEditing} className="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
                        Save
                    </button>
                </div>


            </div>
        )
    }
});

export default Endpoint