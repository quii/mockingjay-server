import React from 'react';
import ReactDOM from 'react-dom';
import Endpoint from './endpoints.jsx';
import _ from 'lodash';

const UI = React.createClass({
    getInitialState: function() {
        return {
            data: [],
            activeEndpoint: null,
            endpointIds: []
        };
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
    getMenuLinks: function () {
        return this.state.data.map(endpoint => {
            let cssClass = "mdl-navigation__link";
            if(endpoint.Name===this.state.activeEndpoint){
                cssClass += " --accent";
            }

            return (
                <a
                    ref={'menu-'+endpoint.Name}
                    className={cssClass}
                    onClick={(event)=>this.openEditor(endpoint.Name, event)}>
                    {endpoint.Name}
                </a>
            )
        });
    },
    openEditor: function (endpointName) {
        this.setState({
            activeEndpoint: endpointName
        })
    },
    updateServer: function () {
        const newEndpointState = this.refs[this.state.activeEndpoint].state;

        let data = _.cloneDeep( this.state.data);

        data[newEndpointState.index] = {
            Name: newEndpointState.name,
            CDCDisabled: newEndpointState.cdcDisabled,
            Request: {
                URI: newEndpointState.uri,
                RegexURI: newEndpointState.regex,
                Method: newEndpointState.method,
                Body: newEndpointState.reqBody,
                Form: newEndpointState.form,
                Headers: newEndpointState.reqHeaders
            },
            Response: {
                Code: parseInt(newEndpointState.code),
                Body: newEndpointState.body,
                Headers: newEndpointState.resHeaders
            }
        };

        const json = JSON.stringify(data);

        this.putUpdate(json);
    },
    renderCurrentEndpoint: function(){
        if(this.state.activeEndpoint) {
            const index = _.findIndex(this.state.data, ep => ep.Name==this.state.activeEndpoint)
            const endpoint = this.state.data.find(ep => ep.Name===this.state.activeEndpoint);
            return (
                <Endpoint
                    index={index}
                    key={endpoint.Name}
                    ref={endpoint.Name}
                    cdcDisabled={endpoint.CDCDisabled}
                    updateServer={this.updateServer}
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
                />);

        }else {
            return null;
        }
    },
    render: function () {
        return (
        <div className="mdl-layout mdl-js-layout mdl-layout--fixed-drawer">
            <div className="mdl-layout__drawer">
                <span className="mdl-layout-title">mockingjay server</span>
                <nav className="mdl-navigation">
                    {this.getMenuLinks()}
                </nav>
            </div>
            <main className="mdl-layout__content">
                <div className="page-content">{this.renderCurrentEndpoint()}</div>
            </main>
        </div>

        )
    }
});
ReactDOM.render(
    <UI url="/mj-endpoints"/>,
    document.getElementById('app')
);