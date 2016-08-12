import React from 'react';

const httpDataList = React.createClass({
    render: function(){
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

export default httpDataList;