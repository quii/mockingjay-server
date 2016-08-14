import React from 'react';

export const HttpDataList = React.createClass({
    render: function () {
        const items = mapKeyVals(this.props.items, (key, val) => <li>{key} -> {val}</li>);
        return <HttpDataView name={this.props.name} items={items}/>
    }
});

export const HttpDataEditor = React.createClass({
    updateMap: function(ref){
        const newState = {};
        for(let i=0; i < Object.keys(this.refs).length; i+=2){
            var keyName = Object.keys(this.refs)[i];
            var valueName = Object.keys(this.refs)[i+1];

            const k = this.refs[keyName].value;
            const v = this.refs[valueName].value;
            newState[k] = v;
        }
        this.props.onChange({
            target: {
                name: this.props.name,
                value: newState
            }
        })
    },
    render: function () {
        const items = mapKeyVals(this.props.items, (key, val, i) => {
            return (<li>
                <input onChange={this.updateMap} ref={i+"key"} type="text" value={key}/> ->
                <input onChange={this.updateMap} ref={i+"value"} type="text" value={val}/>
            </li>);
        });
        return <HttpDataView name={this.props.name} items={items}/>
    }
});

const HttpDataView = React.createClass({
    render: function () {
        if (this.props.items && this.props.items.length > 0) {
            return (
                <div className={this.props.name}>
                    <h3>{this.props.name}</h3>
                    <ul>{this.props.items}</ul>
                </div>
            )
        } else {
            return null;
        }
    }
})

function mapKeyVals(items, f) {
    if (items) {
        let i = -1;
        return Object.keys(items).map(function (key) {
            i++;
            let value = items[key];
            return (
                f(key, value, i)
            )
        });
    }
    return [];
}