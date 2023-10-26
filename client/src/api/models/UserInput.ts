/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { PersonInput } from './PersonInput';

export type UserInput = {
    email: string;
    email_public?: boolean;
    identity: PersonInput;
    login: string;
    password: string;
    password_confirmation: string;
};

