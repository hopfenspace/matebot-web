import { Err, Ok, Result } from "../utils/result";
import { Application, ApplicationResponse, TestResponse } from "./models";

export const Api = {
    login: async (username: string, password: string): Promise<Result<null, any>> => {
        const res = await fetch("/api/frontend/login", {
            method: "post",
            body: JSON.stringify({
                username,
                password,
            }),
        });

        if (res.status !== 200) {
            const parsed = await res.json();
            return Err(parsed.message);
        }

        return Ok(null);
    },
    logout: async (): Promise<Result<null, string>> => {
        const res = await fetch("/api/frontend/logout", {
            method: "get",
        });

        if (res.status !== 200) {
            const parsed = await res.json();
            return Err(parsed.message);
        }

        return Ok(null);
    },
    test: async (): Promise<"logged out" | "logged in"> => {
        const res = await fetch("/api/frontend/test", {
            method: "get",
        });
        const decoded: TestResponse = await res.json();

        if (res.status === 401) {
            return "logged out";
        }

        if (res.status !== 200 || !decoded.authenticated) {
            return "logged out";
        }

        return "logged in";
    },
    connectAccount: async (
        username: string,
        password: string,
        existing_username: string,
        application: string
    ): Promise<Result<null, string>> => {
        const res = await fetch("/api/frontend/connectAccount", {
            method: "post",
            body: JSON.stringify({
                username,
                password,
                existing_username,
                application,
            }),
        });

        if (res.status !== 201) {
            const parsed = await res.json();
            return Err(parsed.message);
        }

        return Ok(null);
    },
    register: async (username: string, password: string): Promise<Result<null, string>> => {
        const res = await fetch("/api/frontend/register", {
            method: "post",
            body: JSON.stringify({
                username,
                password,
            }),
        });

        if (res.status !== 201) {
            const parsed = await res.json();
            return Err(parsed.message);
        }

        return Ok(null);
    },
    get_applications: async (): Promise<Result<Application[], string>> => {
        const res = await fetch("/api/frontend/applications", {
            method: "get",
        });

        const parsed: ApplicationResponse = await res.json();

        if (res.status !== 200) {
            return Err(parsed.message);
        }

        return Ok(parsed.applications);
    },
};
