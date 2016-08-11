import React from 'react';
import ReactDOM from 'react-dom';

var Endpoint = React.createClass({
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

var HttpDataList = React.createClass({
    render: function(){
        if(this.props.items) {
            self = this;
            var items = Object.keys(this.props.items).map(function (key) {
                let value = self.props.items[key];
                return (
                    <li>{key} -> {value}</li>
                )
            });
            return (
                <div className={this.props.name}>
                    <h3>{this.props.name}</h3>
                    <ul>{items}</ul>
                </div>
            )
        }else{
            return null;
        }
    }
})

var EndpointList = React.createClass({
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
})

var EndpointForm = React.createClass({
    render: function () {
        return (
            <form className="endpointForm">
                <label htmlFor="name">Name</label><input type="text" name="name"/>
                <input type="submit" value="Save"/>
            </form>
        );
    }
});

var UI = React.createClass({
    getInitialState: function() {
        return {data: []};
    },
    componentDidMount: function() {
        $.ajax({
            url: this.props.url,
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({data: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    render: function () {
        return (
            <div className="ui">
                <EndpointList data={this.state.data}/>
                <EndpointForm/>
            </div>
        )
    }
});

ReactDOM.render(
    <UI url="/mj-endpoints"/>,
    document.getElementById('app')
);