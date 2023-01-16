export type TestResponse = {
    message: string;
    authenticated: boolean;
};

export type Application = {
    id: number;
    name: string;
    created: number;
};

export type ApplicationResponse = {
    message: string;
    applications: Array<Application>;
};
