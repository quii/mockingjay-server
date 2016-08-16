import React from 'react';

const CDC = React.createClass({
    checkCompatability: function () {
        if(this.state.url) {
            $.ajax({
                url: this.state.url,
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
            const url = this.props.url + "?url=" + e.target.value;

            this.setState({
                url
            });

            this.checkCompatability();
        }
    },
    render: function () {
        let checkDetails;
        if(this.state && this.state.data){
            checkDetails = this.state.data.Passed ? passing : failing;
        }else{
            checkDetails = dunno;
        }
        return (
            <header className="mdl-layout__header">
                <div className="mdl-layout__header-row">
                    <div className="mdl-layout-spacer"></div>
                    <div className="mdl-textfield mdl-js-textfield mdl-textfield--expandable
                  mdl-textfield--floating-label mdl-textfield--align-right">
                        {checkDetails}
                        <label className=""
                               htmlFor="fixed-header-drawer-exp">
                            Check endpoints against real URL
                        </label>
                        <div className="mdl-textfield__expandable-holder">
                            <input className="mdl-textfield__input" type="text" name="sample"
                                   id="fixed-header-drawer-exp" onBlur={this.checkCompatability} onKeyPress={this.handleUrlChange} />
                        </div>
                    </div>
                </div>
            </header>

        )
    }
});

const TestIndicator = React.createClass({
    render: function () {
        return <div className="material-icons mdl-badge mdl-badge--overlap md-48" data-badge={this.props.badge}>compare_arrows</div>
    }
});

const passing = <TestIndicator badge="✓"/>
const failing = <TestIndicator badge="✘" />
const dunno = <TestIndicator badge="?" />

function isValidURL(str) {
    var a  = document.createElement('a');
    a.href = str;
    return a.host;
}

export default CDC;