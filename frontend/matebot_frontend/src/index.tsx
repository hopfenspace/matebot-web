import React from "react";
import ReactDOM from "react-dom/client";
import { ToastContainer } from "react-toastify";

import "./index.css";
import "react-toastify/dist/ReactToastify.css";
import Home from "./views/home";

type RootProps = {};

type RootState = {
    path: Array<String>;
};

class Root extends React.Component<RootProps, RootState> {
    constructor(props: RootProps) {
        super(props);

        this.state = {
            path: [],
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
    <React.StrictMode>
        <Root />
        <ToastContainer theme="dark" autoClose={2500} />
    </React.StrictMode>
);
