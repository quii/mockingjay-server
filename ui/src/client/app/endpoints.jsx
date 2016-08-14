import React from 'react';
import {HttpDataList, HttpDataEditor} from './httpDataList.jsx';

const Endpoint = React.createClass({
    getInitialState: function() {
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
    startEditing: function(){
        this.setState({
            isEditing: true
        })
    },
    finishEditing: function(){
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
    render: function() {
        var view = (<div className="endpoint">
            <h1>{this.state.name} <button onClick={this.startEditing}>Edit</button></h1>
            <div className="request">
                <h4>Request</h4>
                <p>Method <span className="method" onClick={this.startEditing}>{this.state.method}</span></p>
                <p>URI <span className="uri">{this.state.uri}</span></p>
                <p>Regex URI<code className="regex">{this.state.regex}</code></p>
                <p>Body <span className="reqBody">{this.state.reqBody}</span></p>
                <p>Form data <HttpDataList name="Form data" items={this.state.form} /></p>
                <p>Headers <HttpDataList name="Request headers" items={this.state.reqHeaders} /></p>
            </div>
            <div className="response">
                <h4>Response</h4>
                <p>Status code <span className="code">{this.state.code}</span></p>
                <p>Body <code className="body">{this.state.body}</code></p>
                <p>Headers <HttpDataList name="Response headers" items={this.state.resHeaders} /></p>
            </div>
        </div>);

        var form = <EndpointForm
            name={this.state.name}
            finishEditing={this.finishEditing}
            originalValues={this.state}
            onChange={this.updateValue}                                                             ss
        />;

        return <div>{this.state.isEditing ? form : view}</div>;
    }
});

const EndpointForm = React.createClass({
    render: function () {
        return (
            <div class="editor">
                <h1>Editing {this.props.name}</h1>
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
                <label>URI</label><input type="text" name="uri" value={this.props.originalValues.uri} onChange={this.props.onChange} /><br />
                <label>Regex URI</label><input type="text" name="regex" value={this.props.originalValues.regex} onChange={this.props.onChange} /><br />
                <label>Body</label><input type="text" name="reqBody" value={this.props.originalValues.reqBody} onChange={this.props.onChange} /><br />
                <label>Form</label><HttpDataEditor name="form" items={this.props.originalValues.form} />

                <h4>Response</h4>
                <label>Status code</label><input type="text" name="code" value={this.props.originalValues.code} onChange={this.props.onChange} /><br />
                <label>Body</label><input type="text" name="body" value={this.props.originalValues.body} onChange={this.props.onChange} /><br />
                <button onClick={this.props.finishEditing}>Finish editing</button>
            </div>
        )
    }
})

const EndpointList = React.createClass({
    getInitialState: function() {
        return {
            endpointIds: []
        };
    },
    addEndpoint: function (id) {
        this.state.endpointIds.push(id)
    },
    updateServer: function(){
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
                    Body: state.reqBody
                },
                Response: {
                    Code: parseInt(state.code),
                    Body: state.body
                }
            };
        });
        this.props.putUpdate(JSON.stringify(updatedEndpoints));
    },
    render: function () {
        self = this;
        var i = 0;
        var endpointElements = this.props.data.map(function(endpoint) {
            const endpointName = 'endpoint'+i;
            self.addEndpoint(endpointName);
            i++;
            return (
                <Endpoint
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
            <div className="endpointList">
                {endpointElements}
            </div>
        );
    }
});

export default EndpointList