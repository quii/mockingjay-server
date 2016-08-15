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
        $.ajax({
            url: this.props.url,
            dataType: 'json',
            type: 'PUT',
            cache: false,
            data: update,
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
            <div className="mdl-layout mdl-js-layout mdl-layout--fixed-header">
                <header className="mdl-layout__header">
                    <div className="mdl-layout__header-row">
                        <span className="mdl-layout-title">mockingjay</span>
                        </div>
                    </header>
            <div className="ui">
                <EndpointList putUpdate={this.putUpdate} data={this.state.data}/>
            </div>
                </div>
        )
    }
});

ReactDOM.render(
    <UI url="/mj-endpoints"/>,
    document.getElementById('app')
);