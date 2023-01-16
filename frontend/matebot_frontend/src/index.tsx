import React from "react";
import ReactDOM from "react-dom/client";

import "./index.css";
import "react-toastify/dist/ReactToastify.css";

import Home from "./views/home";
import Login from "./views/login";
import Register from "./views/register";
import { ToastContainer } from "react-toastify";

type RootProps = {};

type RootState = {
    path: Array<String>;
    loggedIn: "logged out" | "logged in";
};

export default class Root extends React.Component<RootProps, RootState> {
    constructor(props: RootProps) {
        super(props);

        this.state = {
            path: [],
            loggedIn: "logged out",
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
    }

    render() {
        const { path } = this.state;

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
        <ToastContainer theme="dark" autoClose={2500} />
    </>
);
