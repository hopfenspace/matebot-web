import React from "react";
import Input from "../components/input";
import { Api } from "../api/api";
import { toast } from "react-toastify";
import logo from "../icons/logo.png";

type LoginState = {
    username: string;
    password: string;
};

type LoginProps = {};

export default class Login extends React.Component<LoginProps, LoginState> {
    constructor(props: LoginProps) {
        super(props);

        this.state = {
            username: "",
            password: "",
        };

        this.login = this.login.bind(this);
    }

    async login(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();

        (await Api.login(this.state.username, this.state.password)).match(
            (_) => {
                toast.success("Logged in");
                window.location.hash = "/";
            },
            (err) => {
                toast.error(err);
            }
        );
    }

    render() {
        return (
            <div className="login">
                <form className="login-form" onSubmit={this.login} method="post">
                    <img src={logo} alt="logo icon" />
                    <h1 className="heading">Login</h1>
                    <label htmlFor="username">Username</label>
                    <Input
                        required={true}
                        id="username"
                        autoFocus={true}
                        className="input"
                        value={this.state.username}
                        onChange={(v: string) => this.setState({ username: v })}
                    />
                    <label htmlFor="password">Password</label>
                    <Input
                        required={true}
                        id="password"
                        className="input"
                        type="password"
                        value={this.state.password}
                        onChange={(v: string) => this.setState({ password: v })}
                    />
                    <button className="button">Login</button>
                    <hr />
                    <button
                        type="button"
                        className="button"
                        onClick={(v) => {
                            v.stopPropagation();

                            document.location.hash = "/register";
                        }}
                    >
                        Register instead
                    </button>
                </form>
            </div>
        );
    }
}
