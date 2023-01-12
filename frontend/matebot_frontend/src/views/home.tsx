import React from "react";

type HomeProps = {};

type HomeState = {};

export default class Home extends React.Component<HomeProps, HomeState> {
    constructor(props: HomeProps) {
        super(props);
    }

    render() {
        return (
            <div className="dashboard">
                <div className="dashboard-global-statistics panel">
                    <h1 className="heading">Global statistics</h1>
                </div>
                <div className="dashboard-balance panel">
                    <h1 className="heading">Balance</h1>
                </div>
                <div className="dashboard-polls panel">
                    <h1 className="heading">Polls</h1>
                </div>
                <div className="dashboard-transfer panel">
                    <h1 className="heading">Transfer</h1>
                </div>
                <div className="dashboard-transaction-history panel">
                    <h1 className="heading">Transaction history</h1>
                </div>
                <div className="dashboard-transaction-analysis panel">
                    <h1 className="heading">Transaction analysis</h1>
                </div>
                <div className="dashboard-user-state panel">
                    <h1 className="heading">User state</h1>
                </div>
                <div className="dashboard-consume panel">
                    <h1 className="heading">Consume</h1>
                </div>
            </div>
        );
    }
}
