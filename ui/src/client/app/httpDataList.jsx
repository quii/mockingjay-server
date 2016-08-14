import React from 'react';

export const HttpDataList = React.createClass({
    render: function(){
        if(this.props.items) {
            const items = mapKeyVals(this.props.items, (key, val) => <li>{key} -> {val}</li>)
            return <HttpDataView name={this.props.name} items={items} />
        }else{
            return null;
        }
    }
});

export const HttpDataEditor = React.createClass({
    render: function(){
        if(this.props.items) {
            const items = mapKeyVals(this.props.items, (key, val) => <li><input type="text" value={key}/> -> <input type="text" value={val}/></li>)
            return <HttpDataView name={this.props.name} items={items} />
        }else{
            return null;
        }
    }
});

const HttpDataView = React.createClass({
    render: function(){
        return (
            <div className={this.props.name}>
                <h3>{this.props.name}</h3>
                <ul>{this.props.items}</ul>
            </div>
        )
    }
})

function mapKeyVals(items, f){
    return Object.keys(items).map(function (key) {
        let value = items[key];
        return (
            f(key, value)
        )
    });
}