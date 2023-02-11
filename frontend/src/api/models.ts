import exp from "constants";

export type Application = {
    id: number;
    name: string;
};

export type Alias = {
    id: number;
    user_id: number;
    application_id: number;
    username: string;
    confirmed: boolean;
};

export type DebtorUser = {
    user_id: number;
    username: string;
    balance: number;
    balance_formatted: string;
    active: boolean;
};

export type User = {
    user_id: number;
    core_id: number;
    username: string;
    balance: number;
    balance_formatted: string;
    permission: boolean;
    active: boolean;
    external: boolean;
    voucher_id: any;
    aliases: Alias[];
    debtors: DebtorUser[];
    created: number;
    modified: number;
};

export type Consumable = {
    name: string;
    description: string;
    price: number;
};

export type SimpleUser = {
    user_id?: number;
    core_id: number;
    username: string;
};

export type Transaction = {
    id: number;
    sender: SimpleUser;
    receiver: SimpleUser;
    amount: number;
    reason?: string;
    timestamp: number;
};

/*
Responses below
 */

export type ConsumeTransactionResponse = {
    message: string;
    transaction: Transaction;
};

export type ConsumablesResponse = {
    message: string;
    consumables: Consumable[];
};

export type ApplicationResponse = {
    message: string;
    applications: Array<Application>;
};

export type BlameResponse = {};

export type BalanceResponse = {
    message: string;
    user_id?: number;
    username?: number;
    balance: number;
    balance_formatted: string;
};

export type StateResponse = {
    message: string;
    user: User;
};

export type TestResponse = {
    message: string;
    authenticated: boolean;
};
