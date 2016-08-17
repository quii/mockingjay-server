import React from 'react';
import {rand} from './util';

export const MethodSwitcher = React.createClass({

    selectedCSS: "mdl-button mdl-js-button mdl-button--raised mdl-button--accent",
    notSelectedCSS:  "mdl-button mdl-js-button mdl-button--raised mdl-button--colored",
    methods: ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"],

    handleClick: function (e) {
        this.props.onChange({
            target: {
                name: 'method',
                value:e.target.innerText
            }
        })
    },
    createButton: function (methodName, selectedMethod) {
        const clz = methodName===selectedMethod ? this.selectedCSS : this.notSelectedCSS;
        return <button key={rand()} style={{"marginRight": "10px"}} className={clz} onClick={this.handleClick}>{methodName}</button>
    },
    render: function () {
        const buttons = this.methods.map(m => this.createButton(m, this.props.selected));
        return <div>{buttons}</div>
    }
});

export const TextField = React.createClass({
    render: function () {
        const value = this.props.value || "";
        const label = this.props.label || this.props.name;
        const errorSpan = this.props.errMsg ? <span className="mdl-textfield__error">{this.props.errMsg}</span> : null
        return (
            <div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
                <input pattern={this.props.pattern} ref="user" className="mdl-textfield__input" type="text" name={this.props.name} value={value} onChange={this.props.onChange} />
                <label className="mdl-textfield__label" htmlFor={this.props.name}>{label}</label>
                {errorSpan}
            </div>
        )
    }
});

export const TextArea = React.createClass({
    render: function () {
        const label = this.props.label || this.props.name;
        return (
            <div style={{width: "100%"}} className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
                <textarea ref="user" className="mdl-textfield__input" type="text" rows="5" name={this.props.name} value={this.props.value} onChange={this.props.onChange} />
                <label className="mdl-textfield__label" htmlFor={this.props.name}>{label}</label>
            </div>
        )
    }
});


export const Code = React.createClass({
    render: function () {
        if (this.props.value) {
            return <div className="mdl-card__supporting-text"><code className="mdl-color-text--accent">{this.props.value}</code></div>
        } else {
            return null;
        }
    }
});

export const Body = React.createClass({
    isJSON: function () {
        try {
            JSON.parse(this.props.value);
            return true;
        }
        catch (e) {
            return false
        }
    },
    renderText: function () {
        if (this.isJSON()) {
            return JSON.stringify(JSON.parse(this.props.value), null, 2);
        } else {
            return this.props.value;
        }
    },
    render: function () {
        if (this.props.value) {
            return (
                <div>
                    <div className="mdl-card__title mdl-card--expand">
                        <h6 className="mdl-card__title-text">{this.props.label}</h6>
                    </div>
                    <div className="mdl-card__supporting-text">
                        <pre className="mdl-color-text--primary">{this.renderText()}</pre>
                    </div>
                </div>
            )
        } else {
            return null;
        }
    }
});