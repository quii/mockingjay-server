import React from 'react';
import ReactDOM from 'react-dom';

var Endpoint = React.createClass({
    render: function() {
        return (
            <div className="endpoint">
                I am an endpoint
            </div>
        );
    }
});

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
                <Endpoint/>
                <EndpointForm/>
            </div>
        )
    }
})

ReactDOM.render(
    <UI/>,
    document.getElementById('app')
);