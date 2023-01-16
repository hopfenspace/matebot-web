import React from "react";
import Input from "../components/input";
import { Api } from "../api/api";
import { toast } from "react-toastify";
import logo from "../icons/logo.png";
import Select from "react-select";
import { apiToOption, Option, SELECT_PROPS } from "../components/select";

type RegisterProps = {};
type RegisterState = {
    create_new: boolean;

    // Lists for the dropdowns
    applications: Option[] | "loading";
    // Controlled state
    username: string;
    password: string;
    application: Option | null;
    existing_username: string;
};
export default class Register extends React.Component<RegisterProps, RegisterState> {
    constructor(props: RegisterProps) {
        super(props);

        this.state = {
            username: "",
            password: "",
            application: null,
            existing_username: "",
            create_new: true,
            applications: [],
        };

        this.register_new = this.register_new.bind(this);
        this.register_existing = this.register_existing.bind(this);
    }

    async register_new(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();

        (await Api.register(this.state.username, this.state.password)).match(
            (_) => {
                toast.success("Registered successfully");
                window.location.hash = "/login";
            },
            (err) => {
                toast.error(err);
            }
        );
    }

    async register_existing(e: React.FormEvent<HTMLFormElement>) {
        if (this.state.application === null) {
            toast.error("Application is required");
            return;
        }

        e.preventDefault();

        (
            await Api.connectAccount(
                this.state.username,
                this.state.password,
                this.state.existing_username,
                this.state.application.label
            )
        ).match(
            (_) => {
                toast.success("Connected account successfully");
                window.location.hash = "/login";
            },
            (err) => {
                toast.error(err);
            }
        );
    }

    componentDidMount() {
        Api.get_applications().then((res) =>
            res.match(
                (api_applications) => {
                    const applications = api_applications.map((v) => {
                        return apiToOption(v);
                    });
                    this.setState({ applications });
                },
                (err) => {
                    toast.error(err);
                }
            )
        );
    }

    render() {
        let form;

        if (this.state.create_new) {
            form = (
                <form className="login-form" onSubmit={this.register_new} method="post">
                    <img src={logo} alt="logo icon" />
                    <h1 className="heading">Register</h1>
                    <div className="register-options">
                        <button className="register-option-button-pressed" type="button" disabled={true}>
                            New account
                        </button>
                        <button
                            className="register-option-button"
                            type="button"
                            onClick={(_) => {
                                this.setState({ create_new: false });
                            }}
                        >
                            Connect existing account
                        </button>
                    </div>
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
                    <button className="button">Register</button>
                    <hr />
                    <button
                        type="button"
                        className="button"
                        onClick={(v) => {
                            v.stopPropagation();

                            document.location.hash = "/login";
                        }}
                    >
                        Login instead
                    </button>
                </form>
            );
        } else {
            form = (
                <form className="login-form" onSubmit={this.register_existing} method="post">
                    <img src={logo} alt="logo icon" />
                    <h1 className="heading">Register</h1>
                    <div className="register-options">
                        <button
                            className="register-option-button"
                            type="button"
                            onClick={(_) => {
                                this.setState({ create_new: true });
                            }}
                        >
                            New account
                        </button>
                        <button className="register-option-button-pressed" type="button" disabled={true}>
                            Connect existing account
                        </button>
                    </div>
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
                    <label htmlFor="application">Application</label>
                    <Select
                        {...SELECT_PROPS}
                        isLoading={this.state.applications === "loading"}
                        options={this.state.applications === "loading" ? [] : this.state.applications}
                        onChange={(v) => this.setState({ application: v })}
                        isClearable={true}
                        value={this.state.application}
                        required={true}
                        id="application"
                    />
                    <label htmlFor="existing_username">Existing username</label>
                    <Input
                        required={true}
                        id="existing_username"
                        className="input"
                        value={this.state.existing_username}
                        onChange={(v: string) => this.setState({ existing_username: v })}
                    />
                    <button className="button">Connect account</button>
                    <hr />
                    <button
                        type="button"
                        className="button"
                        onClick={(v) => {
                            v.stopPropagation();

                            document.location.hash = "/login";
                        }}
                    >
                        Login instead
                    </button>
                </form>
            );
        }

        return <div className="login">{form}</div>;
    }
}
