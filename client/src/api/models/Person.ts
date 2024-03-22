/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Meta } from './Meta';
import type { people_InstitutionInner } from './people_InstitutionInner';
import type { UserRole } from './UserRole';
export type Person = {
    alias: string;
    comment?: string;
    contact?: string;
    first_name: string;
    full_name: string;
    id: string;
    institutions?: Array<people_InstitutionInner>;
    last_name: string;
    meta?: Meta;
    role?: UserRole;
};

