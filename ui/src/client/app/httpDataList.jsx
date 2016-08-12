import React from 'react';

export const HttpDataList = React.createClass({
    render: function(){
        console.log('rendering', this.props.items)
        if(this.props.items) {
            self = this;
            var items = Object.keys(this.props.items).map(function (key) {
                let value = self.props.items[key];
                return (
                    <li>{key} -> {value}</li>
                )
            });
            return (
                <div className={this.props.name}>
                    <h3>{this.props.name}</h3>
                    <ul>{items}</ul>
                </div>
            )
        }else{
            return null;
        }
    }
});

export const HttpDataEditor = React.createClass({
    render: function(){
        if(this.props.items) {
            self = this;
            var items = Object.keys(this.props.items).map(function (key) {
                let value = self.props.items[key];
                return (
                    <li><input type="text" value={key}/> -> <input type="text" value={value}/></li>
                )
            });
            return (
                <div className={this.props.name}>
                    <h3>{this.props.name}</h3>
                    <ul>{items}</ul>
                </div>
            )
        }else{
            return null;
        }
    }
});
