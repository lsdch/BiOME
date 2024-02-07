/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Meta } from './Meta';
import type { Person } from './Person';
export type Institution = {
    code: string;
    description?: string;
    id: string;
    meta: Meta;
    name: string;
    people?: Array<Person>;
};

