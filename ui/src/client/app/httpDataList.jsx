import React from 'react';

export const HttpDataList = React.createClass({
    render: function () {
        const items = mapKeyVals(this.props.items, (key, val) => <li>{key} -> {val}</li>)
        return <HttpDataView name={this.props.name} items={items}/>
    }
});

export const HttpDataEditor = React.createClass({
    render: function () {
        const items = mapKeyVals(this.props.items, (key, val) => <li><input type="text" value={key}/> -> <input
            type="text" value={val}/></li>)
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
    if (items && items.length > 0) {
        return Object.keys(items).map(function (key) {
            let value = items[key];
            return (
                f(key, value)
            )
        });
    }
    return [];
}