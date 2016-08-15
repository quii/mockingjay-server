import React from 'react';
import {HttpDataList, HttpDataEditor} from './httpDataList.jsx';

const Endpoint = React.createClass({
    getInitialState: function () {
        return {
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
            <div className="mdl-card mdl-shadow--2dp">
                <div className="mdl-card__title" style={{width: "90%"}}>
                    <h3 className="mdl-card__title-text">{this.state.name}</h3>
                </div>

                <Chip icon="cloud" value={this.state.method + " " + this.state.uri}/>
                <Chip icon="face" value={this.state.regex}/>
                <Body label="Request body" value={this.state.reqBody}/>

                <HttpDataList name="Form data" items={this.state.form}/>

                <Chip icon="face" value={this.state.code}/>
                <HttpDataList name="Request headers" items={this.state.reqHeaders}/>
                <Body label="Response body" value={this.state.body}/>
                <HttpDataList name="Response headers" items={this.state.resHeaders}/>

                <div className="mdl-card__menu">
                    <button onClick={this.startEditing}
                            className="mdl-button mdl-button--icon mdl-js-button mdl-js-ripple-effect">
                        <i className="material-icons">edit</i>
                    </button>
                </div>
            </div>);

        var form = <EndpointForm
            name={this.state.name}
            finishEditing={this.finishEditing}
            originalValues={this.state}
            onChange={this.updateValue}
            onCheckboxChange={this.updateCheckbox}
        />;

        return <div className="mdl-cell mdl-cell--4-col">{this.state.isEditing ? form : view}</div>;
    }
});

const EndpointForm = React.createClass({
    render: function () {
        return (
            <div className="mdl-card mdl-shadow--2dp">

                <div className="mdl-card__title">
                    <h3 className="mdl-card__title-text">Editing {this.props.name}</h3>
                </div>

                <label>CDC Disabled?</label><input type="checkbox"
                                                   defaultChecked={this.props.originalValues.cdcDisabled}
                                                   name="cdcDisabled" onClick={this.props.onCheckboxChange}/><br />
                <h4>Request</h4>
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

                <h4>Response</h4>
                <label>Status code</label><input type="text" name="code" value={this.props.originalValues.code}
                                                 onChange={this.props.onChange}/><br />
                <label>Body</label><textarea name="body"
                                             onChange={this.props.onChange}>{this.props.originalValues.body}</textarea><br />
                <label>Headers</label><HttpDataEditor name="resHeaders" items={this.props.originalValues.resHeaders}
                                                      onChange={this.props.onChange}/>
                <button onClick={this.props.finishEditing}>Finish editing</button>
            </div>
        )
    }
})

const EndpointList = React.createClass({
    getInitialState: function () {
        return {
            endpointIds: []
        };
    },
    addEndpoint: function (id) {
        this.state.endpointIds.push(id)
    },
    updateServer: function () {
        self = this;
        const updatedEndpoints = this.state.endpointIds.map(function (ref) {
            const state = self.refs[ref].state;
            return {
                Name: state.name,
                CDCDisabled: state.cdcDisabled,
                Request: {
                    URI: state.uri,
                    RegexURI: state.regex,
                    Method: state.method,
                    Body: state.reqBody,
                    Form: state.form,
                    Headers: state.reqHeaders
                },
                Response: {
                    Code: parseInt(state.code),
                    Body: state.body,
                    Headers: state.resHeaders
                }
            };
        });
        this.props.putUpdate(JSON.stringify(updatedEndpoints));
    },
    render: function () {
        self = this;
        var i = 0;
        var endpointElements = this.props.data.map(function (endpoint) {
            const endpointName = 'endpoint' + i;
            self.addEndpoint(endpointName);
            i++;
            return (
                <Endpoint
                    key={endpointName}
                    ref={endpointName}
                    cdcDisabled={endpoint.CDCDisabled}
                    updateServer={self.updateServer}
                    name={endpoint.Name}
                    method={endpoint.Request.Method}
                    reqBody={endpoint.Request.Body}
                    uri={endpoint.Request.URI}
                    regex={endpoint.Request.RegexURI}
                    reqHeaders={endpoint.Request.Headers}
                    form={endpoint.Request.Form}
                    code={endpoint.Response.Code}
                    body={endpoint.Response.Body}
                    resHeaders={endpoint.Response.Headers}
                />
            );
        });
        return (
            <div className="mdl-grid">
                {endpointElements}
            </div>
        );
    }
});

const Chip = React.createClass({
    render: function () {
        if (this.props.value) {
            return (
                <span className="mdl-chip mdl-chip--contact">
                        <span className="mdl-chip__contact mdl-color--teal mdl-color-text--white"><i
                            className="material-icons">{this.props.icon}</i></span>
                        <span className="mdl-chip__text">{this.props.value}</span>
                    </span>
            );
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
                    <pre>{this.renderText()}</pre>
                </div>
            )
        } else {
            return null;
        }
    }
})

export default EndpointList