import React from 'react';
import HttpDataList from './httpDataList.jsx';

const Endpoint = React.createClass({
    getInitialState: function() {
        return {
            isEditing: false,
            name: this.props.name,
            method: this.props.method,
            uri: this.props.uri,
            regex: this.props.regex,
            form: this.props.form,
            reqHeaders: this.props.reqHeaders,
            status: this.props.status,
            body: this.props.body,
            form: this.props.resHeaders,
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
        })
    },
    updateValue: function (e) {
        this.setState({
            [e.target.name]: e.target.value
        })
    },
    render: function() {
        var view = (<div className="endpoint">
            <h1>{this.state.name}</h1>
            <button onClick={this.startEditing}>Edit</button>
            <div className="request">
                <h4>Request</h4>
                <span className="method" onClick={this.startEditing}>{this.state.method}</span>
                <span className="uri">{this.state.uri}</span>
                <code className="regex">{this.state.regex}</code>
                <HttpDataList name="Form data" items={this.state.form} />
                <HttpDataList name="Request headers" items={this.state.reqHeaders} />
            </div>
            <div className="response">
                <h4>Response</h4>
                <span className="code">{this.state.status}</span>
                <code className="code">{this.state.body}</code>
                <HttpDataList name="Response headers" items={this.state.resHeaders} />
            </div>
        </div>);

        var form = <EndpointForm
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
                <h4>Request</h4>
                <label>Method</label><input type="text" name="method" value={this.props.originalValues.method} onChange={this.props.onChange} /><br />
                <label>URI</label><input type="text" name="uri" value={this.props.originalValues.uri} onChange={this.props.onChange} /><br />
                <label>Regex URI</label><input type="text" name="regex" value={this.props.originalValues.regex} onChange={this.props.onChange} /><br />
                <h4>Response</h4>
                <label>Status</label><input type="text" name="regex" value={this.props.originalValues.status} onChange={this.props.onChange} /><br />
                <label>Body</label><input type="text" name="regex" value={this.props.originalValues.body} onChange={this.props.onChange} /><br />
                <button onClick={this.props.finishEditing}>Save</button>
            </div>
        )
    }
})

const EndpointList = React.createClass({
    render: function () {
        var endpoints = this.props.data.map(function(endpoint) {
            return (
                <Endpoint
                    name={endpoint.Name}
                    method={endpoint.Request.Method}
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
                {endpoints}
            </div>
        );
    }
});

export default EndpointList