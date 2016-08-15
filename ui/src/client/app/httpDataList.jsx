import React from 'react';

export const HttpDataList = React.createClass({
    render: function () {
        const items = mapKeyVals(this.props.items, (key, val) => <tr><td className="mdl-data-table__cell--non-numeric">{key}</td><td className="mdl-data-table__cell--non-numeric">{val}</td></tr>);
        if(items.length>0) {
            return <div>
                <div className="mdl-card__title mdl-card--expand">
                    <h6 className="mdl-card__title-text">{this.props.name}</h6>
                </div>
                <table className="mdl-data-table mdl-js-data-table mdl-data-table">
                <thead>
                <tr>
                    <th className="mdl-data-table__cell--non-numeric">Name</th>
                    <th className="mdl-data-table__cell--non-numeric">Value</th>
                </tr>
                </thead>
                <tbody>
                {items}
                </tbody>
            </table>
                </div>
        }else{
            return null;
        }
    }
});

export const HttpDataEditor = React.createClass({
    getInitialState: function(){
        return {
            numberOfItems: this.props.items ? Object.keys(this.props.items).length : 0
        }
    },
    addItem:function () {
        this.setState({
            numberOfItems: this.state.numberOfItems+1
        });
    },
    updateMap: function(ref){
        const newState = {};
        for(let i=0; i < Object.keys(this.refs).length; i+=2){
            var keyName = Object.keys(this.refs)[i];
            var valueName = Object.keys(this.refs)[i+1];

            const k = this.refs[keyName].value;
            const v = this.refs[valueName].value;

            if(k!=="" && v!=="") {
                newState[k] = v;
            }
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
        const remainingItems = this.state.numberOfItems - items.length;

        for(let i=0; i < remainingItems; i++){
            items.push(<li>
                <input onChange={this.updateMap} ref={i+items.length+"key"} type="text"/> ->
                <input onChange={this.updateMap} ref={i+items.length+"value"} type="text"/>

            </li>)
        }

        items.push(<li><button onClick={this.addItem}>+</button></li>);
        return <HttpDataView name={this.props.name} items={items}/>
    }
});

const HttpDataView = React.createClass({
    render: function () {
        if (this.props.items && this.props.items.length > 0) {
            return (
                <div className={this.props.name}>
                    <h5>{this.props.name}</h5>
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