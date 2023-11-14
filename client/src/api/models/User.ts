/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Person } from './Person';
import type { UserRole } from './UserRole';

export type User = {
    email: string;
    identity: Person;
    login: string;
    role: UserRole;
    verified: boolean;
};

