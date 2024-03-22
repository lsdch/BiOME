/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { OptionalPerson } from './OptionalPerson';
import type { UserRole } from './UserRole';
export type people_OptionalUser = {
    email: string;
    email_confirmed: boolean;
    id: string;
    identity: OptionalPerson;
    is_active: boolean;
    login: string;
    role: UserRole;
};

