import React from 'react';
import {HttpDataList, HttpDataEditor} from './httpDataList.jsx';

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
                    <Chip icon="cloud" value={this.state.method + " " + this.state.uri}/>
                    <Chip icon="face" value={this.state.regex}/>
                    <Body label="Body" value={this.state.reqBody}/>
                    <HttpDataList name="Headers" items={this.state.reqHeaders}/>
                    <HttpDataList name="Form data" items={this.state.form}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <Chip icon="face" value={this.state.code}/>
                    <Body label="Body" value={this.state.body}/>
                    <HttpDataList name="Headers" items={this.state.resHeaders}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <button onClick={this.startEditing} className="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
                        Edit
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
    render: function () {
        return (
            <div className="">

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Request</h3>
                    </div>

                    <label>CDC Disabled?</label><input type="checkbox"
                                                       defaultChecked={this.props.originalValues.cdcDisabled}
                                                       name="cdcDisabled" onClick={this.props.onCheckboxChange}/><br />
                    <label>Method</label>
                    <select name="method" value={this.props.originalValues.method} onChange={this.props.onChange}>
                        <option value="GET">GET</option>
                        <option value="POST">POST</option>
                        <option value="DELETE">DELETE</option>
                        <option value="PUT">PUT</option>
                        <option value="PATCH">PATCH</option>
                        <option value="OPTIONS">OPTIONS</option>
                    </select>
                    <br />
                    <label>URI</label><input type="text" name="uri" value={this.props.originalValues.uri}
                                             onChange={this.props.onChange}/><br />
                    <label>Regex URI</label><input type="text" name="regex" value={this.props.originalValues.regex}
                                                   onChange={this.props.onChange}/><br />
                    <label>Body</label><textarea name="reqBody"
                                                 onChange={this.props.onChange}>{this.props.originalValues.reqBody}</textarea><br />
                    <label>Form</label><HttpDataEditor name="form" items={this.props.originalValues.form}
                                                       onChange={this.props.onChange}/>
                    <label>Headers</label><HttpDataEditor name="reqHeaders" items={this.props.originalValues.reqHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <label>Status code</label><input type="text" name="code" value={this.props.originalValues.code}
                                                     onChange={this.props.onChange}/><br />
                    <label>Body</label><textarea name="body"
                                                 onChange={this.props.onChange}>{this.props.originalValues.body}</textarea><br />
                    <label>Headers</label><HttpDataEditor name="resHeaders" items={this.props.originalValues.resHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <button onClick={this.props.finishEditing} className="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
                    Save
                </button>
            </div>
        )
    }
})


const Chip = React.createClass({
    render: function () {
        if (this.props.value) {
            return <div className="mdl-card__supporting-text"><code className="mdl-color-text--accent">{this.props.value}</code></div>
        } else {
            return null;
        }
    }
});

const Body = React.createClass({
    isJSON: function () {
        try {
            const json = JSON.parse(this.props.value);
            return true;
        }
        catch (e) {
            return false
        }
    },
    renderText: function () {
        if (this.isJSON()) {
            return JSON.stringify(JSON.parse(this.props.value), null, 2);
        } else {
            return this.props.value;
        }
    },
    render: function () {
        if (this.props.value) {
            return (
                <div>
                    <div className="mdl-card__title mdl-card--expand">
                        <h6 className="mdl-card__title-text">{this.props.label}</h6>
                    </div>
                    <div className="mdl-card__supporting-text">
                        <pre className="mdl-color-text--primary">{this.renderText()}</pre>
                    </div>
                </div>
            )
        } else {
            return null;
        }
    }
})

export default Endpoint