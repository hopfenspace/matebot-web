import React from "react";
import { Radar, Line } from "react-chartjs-2";
import {
    Chart as ChartJS,
    Tooltip,
    Legend,
    RadialLinearScale,
    LineElement,
    PointElement,
    Filler,
    CategoryScale,
    LinearScale,
    Title,
} from "chart.js";
import { Api } from "../api/api";
import { toast } from "react-toastify";
import logout from "../icons/logout.svg";
import { Consumable } from "../api/models";

import missing_icon from "../icons/missing-metadata.svg";
import mate_icon from "../icons/mate.svg";
import ice_icon from "../icons/ice.svg";
import water_icon from "../icons/water.svg";
import pizza_icon from "../icons/pizza.svg";
import energy_icon from "../icons/energy-can.svg";
import communism_icon from "../icons/communism.svg";
import send_icon from "../icons/send.svg";
import pay_icon from "../icons/receive.svg";
import { formatAmount } from "../utils/format";

type HomeProps = {};

type HomeState = {
    balance_formatted: string;
    balance: number;

    external: boolean;

    // Global statistics
    community_balance: number;
    community_balance_formatted: string;

    consumables: Consumable[];

    // Controlled state
};

const CONSUMABLE_LOOKUP: { [index: string]: string } = {
    drink: mate_icon,
    ice: ice_icon,
    water: water_icon,
    pizza: pizza_icon,
    energy: energy_icon,
};

export default class Home extends React.Component<HomeProps, HomeState> {
    constructor(props: HomeProps) {
        super(props);

        this.state = {
            balance: 0,
            balance_formatted: "0,00 €",
            external: true,
            consumables: [],
            community_balance: 0,
            community_balance_formatted: "0,00 €",
        };

        this.logout = this.logout.bind(this);
        this.updateState = this.updateState.bind(this);
        this.consume = this.consume.bind(this);
    }

    async logout() {
        (await Api.logout()).match(
            (_) => {
                toast.success("Logged out");
                document.location.hash = "/login";
            },
            (err) => {
                toast.error(err);
            }
        );
    }

    componentDidMount() {
        Api.get_consumables().then((res) =>
            res.match(
                (consumables) => this.setState({ consumables }),
                (err) => {
                    console.log(err);
                }
            )
        );

        this.updateState();
    }

    updateState() {
        Api.me().then((res) =>
            res.match(
                (user) => {
                    this.setState({
                        balance: user.balance,
                        balance_formatted: formatAmount(user.balance),
                        external: user.external,
                    });

                    // Retrieve state for internals
                    if (!user.external) {
                        Api.blame().then((v) => {
                            v.match(
                                (balance_response) => {
                                    console.log(balance_response);
                                },
                                (err) => toast.error(err)
                            );
                        });

                        Api.zwegat().then((v) =>
                            v.match(
                                (balance_response) => {
                                    this.setState({
                                        community_balance: balance_response.balance,
                                        community_balance_formatted: formatAmount(balance_response.balance),
                                    });
                                    console.log(balance_response);
                                },
                                (err) => toast.error(err)
                            )
                        );
                    }
                },
                (err) => {
                    toast.error(err);
                }
            )
        );
    }

    async consume(consumable: Consumable) {
        (await Api.consume(consumable, 1)).match(
            async (_) => {
                toast.success("Success");
                // TODO: Only update necessary state
                await this.updateState();
            },
            (err) => toast.error(err)
        );
    }

    render() {
        ChartJS.register(
            CategoryScale,
            LinearScale,
            PointElement,
            LineElement,
            Title,
            Tooltip,
            Filler,
            Legend,
            RadialLinearScale
        );

        const primary = getComputedStyle(document.body).getPropertyValue("--prim");

        if (!primary.startsWith("#")) {
            console.error("Broken css, colors must start with #");
        }

        const labels = [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December",
        ];
        const balance_data = {
            labels,
            datasets: [
                {
                    fill: true,
                    label: "Balance",
                    data: labels.map(() => 384),
                    borderColor: primary,
                    backgroundColor: primary + 0x14,
                },
            ],
        };

        const balance_options = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: "top" as const,
                },
            },
        };

        const radar_data = {
            labels: ["Thing 1", "Thing 2", "Thing 3", "Thing 4", "Thing 5", "Thing 6"],
            outlineColor: "rgba(255, 255, 255, 50)",
            datasets: [
                {
                    label: "Amount spent",
                    data: [2, 9, 3, 5, 2, 3],
                    backgroundColor: "rgba(255, 99, 132, 0.2)",
                    borderColor: "rgba(255, 99, 132, 1)",
                    borderWidth: 1,
                },
            ],
        };

        const radar_options = {
            scales: {
                r: {
                    ticks: {
                        textStrokeColor: "#bbb",
                        color: "#bbb",
                        backdropColor: "#2b2d3f",
                    },
                    angleLines: {
                        color: "gray",
                    },
                    grid: {
                        color: "gray",
                    },
                },
            },
        };

        // Fill consumables
        let consumables = [];
        for (const consumable of this.state.consumables) {
            if (consumable.name in CONSUMABLE_LOOKUP) {
                consumables.push(
                    <button
                        className="consumables-item"
                        onClick={async (_) => {
                            await this.consume(consumable);
                        }}
                    >
                        <img alt="icon" src={CONSUMABLE_LOOKUP[consumable.name]} />
                        <p>
                            {consumable.name} - {formatAmount(consumable.price)}
                        </p>
                    </button>
                );
            } else {
                consumables.push(
                    <button
                        className="consumables-item"
                        onClick={async (_) => {
                            await this.consume(consumable);
                        }}
                    >
                        <img alt="missing icon" src={missing_icon} />
                        <p>
                            {consumable.name} - {formatAmount(consumable.price)}
                        </p>
                    </button>
                );
            }
        }

        return (
            <div className="dashboard">
                <div className="dashboard-global-statistics panel">
                    <h1 className="heading">Global statistics</h1>
                    <table className="global-statistics-table">
                        <tbody>
                            <tr>
                                <td>Community balance:</td>
                                <td className={this.state.community_balance >= 0 ? "green" : "red"}>
                                    {this.state.community_balance_formatted}
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div className="dashboard-balance panel">
                    <h1 className="heading">Balance</h1>
                    <div className="balance-content">
                        <h2 className={this.state.balance >= 0 ? "heading green" : "heading red"}>
                            {this.state.balance_formatted}
                        </h2>
                        <div className="balance-chart">
                            <Line options={balance_options} data={balance_data} />
                        </div>
                    </div>
                </div>
                <div className="dashboard-polls panel">
                    <h1 className="heading">Polls</h1>
                </div>
                <div className="dashboard-transfer panel">
                    <h1 className="heading">Transfer</h1>
                    <div className="transfer-list">
                        <button className="transfer-item">
                            <img alt="communism icon" src={communism_icon} />
                            <p>Communism</p>
                        </button>
                        <button className="transfer-item">
                            <img alt="send icon" src={send_icon} />
                            <p>Send</p>
                        </button>
                        <button className="transfer-item">
                            <img alt="pay icon" src={pay_icon} />
                            <p>Pay</p>
                        </button>
                    </div>
                </div>
                <div className="dashboard-transaction-history panel">
                    <h1 className="heading">Transaction history</h1>
                </div>
                <div className="dashboard-transaction-analysis panel">
                    <h1 className="heading">Transaction analysis</h1>
                    <div>
                        <Radar data={radar_data} options={radar_options} updateMode="resize" />
                    </div>
                </div>
                <div className="dashboard-user-state panel">
                    <h1 className="heading">User state</h1>
                    <div>
                        <button className="icon-button" onClick={this.logout}>
                            Logout
                            <img alt="logout icon" src={logout} />
                        </button>
                    </div>
                </div>
                <div className="dashboard-consume panel">
                    <h1 className="heading">Consume</h1>
                    <div className="consumables-list">{consumables}</div>
                </div>
            </div>
        );
    }
}
