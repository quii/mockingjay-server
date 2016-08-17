import React from 'react';
import {HttpDataList, HttpDataEditor} from './httpDataList.jsx';

const Endpoint = React.createClass({
    getInitialState: function () {
        return {
            index: this.props.index,
            cdcDisabled: this.props.cdcDisabled,
            isEditing: false,
            name: this.props.name,
            method: this.props.method,
            uri: this.props.uri,
            regex: this.props.regex,
            reqBody: this.props.reqBody,
            form: this.props.form,
            reqHeaders: this.props.reqHeaders,
            code: this.props.code,
            body: this.props.body,
            resHeaders: this.props.resHeaders,
        };
    },
    startEditing: function () {
        this.setState({
            isEditing: true
        })
    },
    delete: function () {
        this.props.delete();
    },
    finishEditing: function () {
        this.setState({
            isEditing: false
        });
        this.props.updateServer();
    },
    updateValue: function (e) {
        this.setState({
            [e.target.name]: e.target.value
        })
    },
    updateCheckbox: function (e) {
        this.setState({
            [e.target.name]: e.target.value === 'on'
        })
    },
    render: function () {

        const view = (
            <div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Request</h3>
                    </div>
                    <Chip icon="cloud" value={this.state.method + " " + this.state.uri}/>
                    <Chip icon="face" value={this.state.regex}/>
                    <Body label="Body" value={this.state.reqBody}/>
                    <HttpDataList name="Headers" items={this.state.reqHeaders}/>
                    <HttpDataList name="Form data" items={this.state.form}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <Chip icon="face" value={this.state.code}/>
                    <Body label="Body" value={this.state.body}/>
                    <HttpDataList name="Headers" items={this.state.resHeaders}/>
                </div>

                <div style={{margin:"2% 2% 2% 3%"}}>
                    <button style={{margin:"0% 1% 0% 0%"}} onClick={this.startEditing} className="mdl-button mdl-button--raised mdl-button--accent">
                        Edit
                    </button>
                    <button onClick={this.delete} className="mdl-button mdl-button--raised mdl-button--primary">
                        Delete
                    </button>
                </div>
            </div>);

        const form = <EndpointForm
            name={this.state.name}
            finishEditing={this.finishEditing}
            originalValues={this.state}
            onChange={this.updateValue}
            onCheckboxChange={this.updateCheckbox}
        />;

        return this.state.isEditing ? form : view;
    }
});

const EndpointForm = React.createClass({
    componentDidMount: function () {
        componentHandler.upgradeDom();
    },
    render: function () {
        return (
            <div>
                <div className="mdl-card mdl-shadow--2dp">
                    <TextField name="name" value={this.props.originalValues.name} onChange={this.props.onChange} />
                    <label class="mdl-checkbox mdl-js-checkbox" for="cdcDisabled">
                        <input type="checkbox" onClick={this.props.onCheckboxChange} name="cdcDisabled" class="mdl-checkbox__input" defaultChecked={this.props.originalValues.cdcDisabled} />
                            <span class="mdl-checkbox__label">CDC Disabled?</span>
                    </label>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Request</h3>
                    </div>

                    <TextField name="uri" value={this.props.originalValues.uri} onChange={this.props.onChange} />
                    <TextField name="regex" value={this.props.originalValues.regex} onChange={this.props.onChange} />
                    <MethodSwitcher selected={this.props.originalValues.method} onChange={this.props.onChange} />
                    <TextArea name="reqBody" value={this.props.originalValues.reqBody} onChange={this.props.onChange} />


                    <label>Form</label><HttpDataEditor name="form" items={this.props.originalValues.form}
                                                       onChange={this.props.onChange}/>
                    <label>Headers</label><HttpDataEditor name="reqHeaders" items={this.props.originalValues.reqHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <div className="mdl-card mdl-shadow--2dp">
                    <div className="mdl-card__title" style={{width: "90%"}}>
                        <h3 className="mdl-card__title-text">Response</h3>
                    </div>
                    <label>Status code</label><input type="text" name="code" value={this.props.originalValues.code}
                                                     onChange={this.props.onChange}/><br />
                    <label>Body</label><textarea name="body"
                                                 onChange={this.props.onChange}>{this.props.originalValues.body}</textarea><br />
                    <label>Headers</label><HttpDataEditor name="resHeaders" items={this.props.originalValues.resHeaders}
                                                          onChange={this.props.onChange}/>
                </div>

                <button onClick={this.props.finishEditing} className="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
                    Save
                </button>
            </div>
        )
    }
});

const MethodSwitcher = React.createClass({

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
        return <button style={{"margin-right": "10px"}} className={clz} onClick={this.handleClick}>{methodName}</button>
    },
    render: function () {
        const buttons = this.methods.map(m => this.createButton(m, this.props.selected));
        return <div>{buttons}</div>
    }
});

const TextField = React.createClass({
    render: function () {
        return (
            <div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
                <input ref="user" className="mdl-textfield__input" type="text" name={this.props.name} value={this.props.value} onChange={this.props.onChange} />
                <label className="mdl-textfield__label" htmlFor={this.props.name}>{this.props.name}</label>
                {/*<span class="mdl-textfield__error">Only alphabet and no spaces, please!</span>  add pattern to enable validation http://webdesign.tutsplus.com/tutorials/learning-material-design-lite-text-fields--cms-24614*/}
            </div>
        )
    }
});

const TextArea = React.createClass({
    render: function () {
        return (
            <div className="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
                <textarea ref="user" className="mdl-textfield__input" type="text" rows="5" name={this.props.name} value={this.props.value} onChange={this.props.onChange} />
                <label className="mdl-textfield__label" htmlFor={this.props.name}>{this.props.name}</label>
            </div>
        )
    }
});


const Chip = React.createClass({
    render: function () {
        if (this.props.value) {
            return <div className="mdl-card__supporting-text"><code className="mdl-color-text--accent">{this.props.value}</code></div>
        } else {
            return null;
        }
    }
});

const Body = React.createClass({
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
})

export default Endpoint