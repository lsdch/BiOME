/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Institution } from '../models/Institution';
import type { InstitutionInput } from '../models/InstitutionInput';
import type { User } from '../models/User';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class PeopleService {

    /**
     * Authenticated user details
     * Get details of currently authenticated user
     * @returns User Authenticated user details
     * @throws ApiError
     */
    public static currentUser(): CancelablePromise<User> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/account',
            errors: {
                400: `User is not authenticated`,
            },
        });
    }

    /**
     * List Institutions
     * @returns Institution OK
     * @throws ApiError
     */
    public static getPeopleInstitutions(): CancelablePromise<Array<Institution>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/people/institutions',
        });
    }

    /**
     * Create institution
     * Register a new institution that people work in.
     * @param data Institution informations
     * @returns Institution Accepted
     * @throws ApiError
     */
    public static createInstitution(
        data: InstitutionInput,
    ): CancelablePromise<Institution> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/people/institutions',
            body: data,
            errors: {
                400: `Bad Request`,
            },
        });
    }

    /**
     * Update institution
     * @param data Institution informations
     * @returns Institution Accepted
     * @throws ApiError
     */
    public static updateInstitution(
        data: Institution,
    ): CancelablePromise<Institution> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/people/institutions/',
            body: data,
            errors: {
                400: `Bad Request`,
            },
        });
    }

    /**
     * Delete institution
     * @param acronym Institution short name
     * @returns any Delete successful
     * @throws ApiError
     */
    public static deleteInstitution(
        acronym: string,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/people/institutions/{acronym}',
            path: {
                'acronym': acronym,
            },
            errors: {
                404: `Institution does not exist`,
            },
        });
    }

}
