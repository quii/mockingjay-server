import React from 'react';
import HttpDataList from './httpDataList.jsx';

const Endpoint = React.createClass({
    render: function() {
        return (
            <div className="endpoint">
                <h1>{this.props.name}</h1>
                <div className="request">
                    <span className="method">{this.props.method}</span>
                    <span className="uri">{this.props.uri}</span>
                    <code className="regex">{this.props.regex}</code>
                    <HttpDataList name="Form data" items={this.props.form} />
                    <HttpDataList name="Request headers" items={this.props.reqHeaders} />
                </div>
                <div className="response">
                    <span className="code">{this.props.status}</span>
                    <code className="code">{this.props.body}</code>
                    <HttpDataList name="Response headers" items={this.props.resHeaders} />
                </div>
            </div>
        );
    }
});

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