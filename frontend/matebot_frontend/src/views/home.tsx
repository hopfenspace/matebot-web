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
import { log } from "util";

type HomeProps = {};

type HomeState = {};

export default class Home extends React.Component<HomeProps, HomeState> {
    constructor(props: HomeProps) {
        super(props);

        this.logout = this.logout.bind(this);
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

        return (
            <div className="dashboard">
                <div className="dashboard-global-statistics panel">
                    <h1 className="heading">Global statistics</h1>
                </div>
                <div className="dashboard-balance panel">
                    <h1 className="heading">Balance</h1>
                    <div className="balance-content">
                        <h2 className="heading green">84.23€</h2>
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
                </div>
            </div>
        );
    }
}
