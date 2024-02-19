/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Institution } from './Institution';
import type { Meta } from './Meta';
import type { UserRole } from './UserRole';
export type Person = {
    contact?: string;
    first_name: string;
    full_name: string;
    id: string;
    institutions?: Array<Institution>;
    last_name: string;
    meta?: Meta;
    middle_names?: string;
    role?: UserRole;
};

