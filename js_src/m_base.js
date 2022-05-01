import {React, ReactDOM, ToastContainer, toast} from "./m_react.js";
import ctx from "./lib/m_ctx.js";
import LoginView from "./m_login.js";
import DashboardView from "./m_dashboard.js";


class Main extends React.Component {
    static contextType = ctx;

    constructor(props) {
        super(props);
        this.state = {
            "path": window.location.hash.substring(2).split("/"),
            "logged_in": false
        };

        this.setPath = this.setPath.bind(this);
        this.componentDidMount = this.componentDidMount.bind(this);

        window.addEventListener('hashchange', () => {
            this.setPath({"path": window.location.hash.substring(2).split("/")});
        });
    }

    setPath(state) {
        if (state.path !== undefined) {
            let path = state.path;
            let newHref = window.location.href;
            if (newHref.includes("#")) {
                newHref = newHref.split("#")[0];
            }
            newHref += "#/" + path.join("/");
            window.location.href = newHref;
        }
        this.setState(state);
    }

    componentDidMount() {
        let preloader = document.getElementById("preloader");
        preloader.style.transition = "all 200ms ease-in-out 0s";
        preloader.style.opacity = "0";
        setTimeout(() => {
            preloader.style.display = "none";
            document.body.id = "";
        }, 200);
        this.context.sdk.test().then((v) => {
            if (!v) {
                this.setPath({"path": ["login"], "logged_in": false});
            } else {
                this.setPath({"path": [""], "logged_in": true});
            }
        });
    }

    render() {
        let content = null;

        if (!this.state.logged_in) {
            content = <LoginView />
        } else {
            // First level router
            if (this.state.path.length === 1) {
                if (this.state.path[0] === "login") {
                    content = <LoginView />
                } else if (this.state.path[0] === "") {
                    content = <DashboardView />
                }
            // Second level router
            } else if (this.state.path.length === 2) {

            }
        }
        
        // Create side menu and embed content
        return <ctx.Provider value={{...this.context, "setPath": this.setPath}}>
            <ToastContainer position="top-right"
                            autoClose={2500}
                            hideProgressBar={false}
                            newestOnTop={false}
                            closeOnClick
                            rtl={false}
                            pauseOnFocusLoss
                            draggable
                            theme="colored"
                            pauseOnHover />
            <div className="content" id="content">
                {content}
            </div>
        </ctx.Provider>;
    }
}

ReactDOM.render(<Main/>, document.getElementById("root"));