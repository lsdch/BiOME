/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Meta } from './Meta';
import type { OptionalPerson } from './OptionalPerson';
import type { UserRole } from './UserRole';
export type User = {
    email: string;
    email_confirmed: boolean;
    id: string;
    identity: OptionalPerson;
    is_active: boolean;
    login: string;
    meta?: Meta;
    role: UserRole;
};

