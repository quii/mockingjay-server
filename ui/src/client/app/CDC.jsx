import React from 'react';
import dialogPolyfill from 'dialog-polyfill';
import {rand} from './util';

const CDC = React.createClass({
    getInitialState: function () {
        return {
            remoteUrl: location.origin
        }
    },
    componentWillMount: function () {
        this.checkCompatability()
    },
    checkCompatability: function () {
        if(this.state && this.state.remoteUrl && this.state.remoteUrl!==null) {
            $.ajax({
                url: this.props.url + "?url=" + this.state.remoteUrl,
                dataType: 'json',
                cache: false,
                success: function (data) {
                    this.setState({data: data});
                }.bind(this),
                error: function (xhr, status, err) {
                    console.error(this.props.url, status, err.toString());
                }.bind(this)
            });
        }
    },
    handleUrlChange: function(e) {
        this.setState({
            data: null
        });

        if ((!e.key || e.key === 'Enter') && isValidURL(e.target.value)) {
            this.setState({
                remoteUrl: e.target.value
            }, this.checkCompatability);
        }
    },
    label: function () {
            return "Automatically checking your endpoints are equivalent to (click to change)"
    },
    indicatorClick: function () {
        this.refs['dialog'].showModal();
    },
    render: function () {
        let checkDetails, messages;
        if(this.state && this.state.data){
            checkDetails = this.state.data.Passed ? <TestIndicator indicatorClick={this.indicatorClick} badge="sentiment_satisfied"/> : <TestIndicator indicatorClick={this.indicatorClick} badge="sentiment_very_dissatisfied" />;
            messages = this.state.data.Messages;
        }else{
            checkDetails = <TestIndicator indicatorClick={this.indicatorClick} badge="sentiment_neutral" />;
            messages = [];
        }
        return (
            <header className="mdl-layout__header">
                <div className="mdl-layout__header-row">
                    <div className="mdl-layout-spacer"></div>
                    <div className="mdl-textfield mdl-js-textfield mdl-textfield--expandable
                  mdl-textfield--floating-label mdl-textfield--align-right">
                        <label htmlFor="fixed-header-drawer-exp">{this.label()}</label>
                        <div className="mdl-textfield__expandable-holder">
                            <input className="mdl-textfield__input" type="text" name="sample"
                                   id="fixed-header-drawer-exp" onBlur={this.checkCompatability} onKeyPress={this.handleUrlChange} defaultValue={this.state.remoteUrl} />
                        </div>
                        {checkDetails}
                    </div>
                </div>
                <Dialog title="Messages from CDC check" messages={messages} ref="dialog" />
            </header>

        )
    }
});

const TestIndicator = React.createClass({
    render: function () {
        return <i onClick={this.props.indicatorClick} style={{cursor: "hand"}} className="material-icons md-48">{this.props.badge}</i>
    }
});


const Dialog = React.createClass({
    componentDidMount: function () {
        const dialog = document.querySelector('dialog');

        if(!dialog.showModal) {
            dialogPolyfill.registerDialog(dialog);
        }
    },
    showModal: function () {
        if(this.props.messages && this.props.messages.length > 0) {
            document.querySelector('dialog').showModal();
        }
    },
    close: function () {
        document.querySelector('dialog').close();
    },
    render: function () {
        let errs;
        if(this.props.messages) {
            errs = this.props.messages.map(m => <li key={rand()}>{m}</li>)
        }
        return (
            <dialog className="mdl-dialog" style={{width:"700px"}}>
                <h4 className="mdl-dialog__title">{this.props.title}</h4>
                <div className="mdl-dialog__content" style={{"fontFamily": "Courier"}}>
                    <ul>{errs}</ul>
                </div>
                <div className="mdl-dialog__actions">
                    <button type="button" onClick={this.close} className="mdl-button">Close</button>
                </div>
            </dialog>
        )
    }
})

function isValidURL(str) {
    var a  = document.createElement('a');
    a.href = str;
    return a.host;
}


export default CDC;