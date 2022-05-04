import {React, toast} from "./m_react.js";
import ctx from "./lib/m_ctx.js";
import TextInput from "./lib/m_input.js";

export default class LoginView extends React.Component {
    static contextType = ctx;

    constructor() {
        super();

        this.state = {
            "username": "",
            "password": "",
        };

        this.checkRequiredParameter = this.checkRequiredParameter.bind(this);
    }

    checkRequiredParameter() {
        if (this.state.username === "") {
            toast.error("Username is required");
            return false;
        }

        if (this.state.password === "") {
            toast.error("Password is required");
            return false;
        }

        return true;
    }

    render() {
        return <div className="loginContent">
            <div className="loginBox">
                <img className="logo"
                     alt="logo"
                     src={this.context.static + "img/matebot_white_1024.png"} />
                <label>Username</label>
                <TextInput className="input"
                           value={this.state.username}
                           onChange={(v) => {
                               this.setState({"username": v});
                           }}
                           style={{marginBottom: "0.5rem"}} />
                <label>Password</label>
                <TextInput className="input"
                           value={this.state.password}
                           onChange={(v) => {
                               this.setState({"password": v});
                           }}
                           style={{marginBottom: "1rem"}}
                           type="password" />
                <button className="button"
                        onClick={(v) => {
                            if (!this.checkRequiredParameter()) {
                                return;
                            }
                            this.context.sdk.login(this.state.username, this.state.password).then(r => {
                                if (r) {
                                    toast.success("Logged in");
                                    this.context.setPath({"path": [""], "logged_in": true});
                                }
                            });
                        }} >
                    Login
                </button>
                <hr style={{color: "var(--level-4)", width: "80%", margin: "1rem"}} />
                <button className="button"
                        onClick={(v) => {
                            if (!this.checkRequiredParameter()) {
                                return;
                            }
                            this.context.sdk.register(this.state.username, this.state.password).then(r => {
                                if (r) {
                                    toast.success("Registered successfully");
                                }
                            });
                        }} >
                    Register instead
                </button>
            </div>
        </div>;
    }
}