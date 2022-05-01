import {toast} from "../m_react.js";

export default class SDK {
    constructor() {
        let url = new URL(document.location.origin);
        url.pathname = "/api/frontend/";
        this.baseURL = url.toString();
        this.token = "";
        this.loggedOutCallback = null;
    }

    async test() {
        let req = await fetch(this.baseURL + "test", {
            method: "POST"
        });
        if (req.status !== 200) {
            toast.error("Error while checking logged in status");
            return false
        }
        let res = await req.json();
        return await res.authenticated;
    }

    async logout() {
        let req = await fetch(this.baseURL + "logout", {
            method: "GET",
        });
        if (req.status === 200) {
            return true
        } else {
            let res = await req.json();
            toast.error("Unable to logout: " + res.error);
            return false;
        }
    }

    async login(username, password) {
        let req = await fetch(this.baseURL + "login", {
            method: "POST",
            body: JSON.stringify({
                "username": username,
                "password": password,
            }),
        });
        if (req.status === 401) {
            toast.error("Login failed");
            return false;
        }
        if (req.status !== 200) {
            let res = await req.json();
            toast.error("Unable to login: " + res.error);
            return false;
        }
        return true;
    }

    async register(username, password) {
        let req = await fetch(this.baseURL + "register", {
            method: "POST",
            body: JSON.stringify({
                "username": username,
                "password": password,
            }),
        });
        if (req.status !== 200) {
            let res = await req.json();
            toast.error("Unable to register: " + res.error);
            return false;
        }
        return true;
    }
}