/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type UserCredentials = {
    identifier: string;
    /**
     * Unhashed, password hash comparison is done within EdgeDB
     */
    password: string;
    remember: boolean;
};

