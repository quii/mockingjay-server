import React from 'react';
import ReactDOM from 'react-dom';

var Endpoint = React.createClass({
    render: function() {
        return (
            <div className="endpoint">
                <h1>{this.props.name}</h1>
                <span className="method">{this.props.method}</span>
                <span className="uri">{this.props.uri}</span>
            </div>
        );
    }
});

var EndpointList = React.createClass({
    render: function () {
        var endpoints = this.props.data.map(function(endpoint) {
            return (
                <Endpoint name={endpoint.Name} method={endpoint.Request.Method} uri={endpoint.Request.URI}/>
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
            <div className="endpointForm">
                <label htmlFor="name">Name</label><input type="text" name="name"/>
                <input type="submit" value="Save"/>
            </div>
        );
    }
});

var UI = React.createClass({
    render: function () {
        return (
            <div className="ui">
                <EndpointList data={this.props.data}/>
                <EndpointForm/>
            </div>
        )
    }
});

var data = [{
    "Name": "Test endpoint",
    "CDCDisabled": false,
    "Request": {
        "URI": "/hello",
        "RegexURI": null,
        "Method": "GET",
        "Headers": null,
        "Body": "",
        "Form": null
    },
    "Response": {
        "Code": 200,
        "Body": "{\"message\": \"hello, world\"}",
        "Headers": {
            "content-type": "text/json"
        }
    }
}, {
    "Name": "Test endpoint 2",
    "CDCDisabled": false,
    "Request": {
        "URI": "/world",
        "RegexURI": null,
        "Method": "DELETE",
        "Headers": null,
        "Body": "",
        "Form": null
    },
    "Response": {
        "Code": 200,
        "Body": "hello, world",
        "Headers": null
    }
}, {
    "Name": "Failing endpoint",
    "CDCDisabled": false,
    "Request": {
        "URI": "/card",
        "RegexURI": null,
        "Method": "POST",
        "Headers": null,
        "Body": "Greetings",
        "Form": null
    },
    "Response": {
        "Code": 500,
        "Body": "Oh bugger",
        "Headers": null
    }
}, {
    "Name": "Endpoint not used for CDC",
    "CDCDisabled": true,
    "Request": {
        "URI": "/burp",
        "RegexURI": null,
        "Method": "POST",
        "Headers": null,
        "Body": "Belch",
        "Form": null
    },
    "Response": {
        "Code": 500,
        "Body": "Oh no",
        "Headers": null
    }
}, {
    "Name": "Posting forms",
    "CDCDisabled": false,
    "Request": {
        "URI": "/cats",
        "RegexURI": null,
        "Method": "POST",
        "Headers": null,
        "Body": "",
        "Form": {
            "age": "10",
            "name": "Hudson"
        }
    },
    "Response": {
        "Code": 201,
        "Body": "Created",
        "Headers": null
    }
}]

ReactDOM.render(
    <UI data={data}/>,
    document.getElementById('app')
);