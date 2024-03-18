/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { UserInput } from './UserInput';
export type PendingUserRequestInput = {
    identity?: {
        first_name: string;
        institution?: string;
        last_name: string;
    };
    motive?: string;
    user?: UserInput;
};

