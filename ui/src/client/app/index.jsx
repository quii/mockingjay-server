import React from 'react';
import ReactDOM from 'react-dom';
import EndpointList from './endpoints.jsx';

const UI = React.createClass({
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
    putUpdate: function(update) {
        console.log('I will PUT', update);
    },
    render: function () {
        return (
            <div className="ui">
                <EndpointList putUpdate={this.putUpdate} data={this.state.data}/>
            </div>
        )
    }
});

ReactDOM.render(
    <UI url="/mj-endpoints"/>,
    document.getElementById('app')
);