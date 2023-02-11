import React from "react";
import ReactDOM from "react-dom/client";

import "./index.css";
import "react-toastify/dist/ReactToastify.css";

import Home from "./views/home";
import Login from "./views/login";
import Register from "./views/register";
import { ToastContainer } from "react-toastify";
import { Api } from "./api/api";

type RootProps = {};

type RootState = {
    path: Array<String>;
    loggedIn: "logged out" | "logged in";

    state_set: boolean;
};

export default class Root extends React.Component<RootProps, RootState> {
    constructor(props: RootProps) {
        super(props);

        this.state = {
            path: [],
            loggedIn: "logged out",
            state_set: false,
        };
    }

    componentDidMount() {
        let preloader = document.getElementById("preloader");
        if (preloader !== null) {
            preloader.remove();
        }

        const setPath = () => {
            const rawPath = window.location.hash;

            // Ensure well-formed path i.e. always have a #/
            if (!rawPath.startsWith("#/")) {
                window.location.hash = "#/";

                // this method will be immediately triggered again
                return;
            }

            // Split everything after #/
            const path = rawPath.substring(2).split("/");

            // #/ should result in [] not [""]
            if (path.length === 1 && path[0] === "") {
                path.shift();
            }

            this.setState({ path });
        };

        setPath();
        window.addEventListener("hashchange", setPath);

        Api.test().then((v) => {
            if (v === "logged out") {
                document.location.hash = "/login";
            }
            this.setState({ state_set: true });
        });
    }

    render() {
        const { path, state_set, loggedIn } = this.state;

        // If we haven't set our state, don't render yet
        if (!state_set) {
            return <></>;
        }

        let content = (() => {
            switch (path[0]) {
                case "":
                case undefined:
                    return <Home />;
                case "login":
                    return <Login />;
                case "register":
                    return <Register />;
                default:
                    break;
            }
        })();

        if (content === undefined) {
            return <div>Unknown route</div>;
        }

        return <>{content}</>;
    }
}

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);
root.render(
    <>
        <Root />
        <ToastContainer theme="dark" autoClose={3500} pauseOnHover={true} hideProgressBar={true} draggable={true} />
    </>
);
