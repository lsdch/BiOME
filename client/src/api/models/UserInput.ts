/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { people_PersonInner } from './people_PersonInner';
export type UserInput = {
    email: string;
    /**
     * EmailPublic bool        `edbedb:"email_public" json:"email_public"`
     */
    identity: people_PersonInner;
    login: string;
    password: string;
    password_confirmation: string;
};

