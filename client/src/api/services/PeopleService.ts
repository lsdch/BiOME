/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Institution } from '../models/Institution';
import type { InstitutionInput } from '../models/InstitutionInput';
import type { InstitutionUpdate } from '../models/InstitutionUpdate';
import type { Person } from '../models/Person';
import type { PersonInput } from '../models/PersonInput';
import type { PersonUpdate } from '../models/PersonUpdate';
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
    public static listInstitutions(): CancelablePromise<Array<Institution>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/people/institutions',
        });
    }
    /**
     * Create institution
     * Register a new institution that people work in.
     * @param data Institution informations
     * @returns Institution Created
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
     * Delete institution
     * @param code Institution short name
     * @returns Institution Deleted item
     * @throws ApiError
     */
    public static deleteInstitution(
        code: string,
    ): CancelablePromise<Institution> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/people/institutions/{code}',
            path: {
                'code': code,
            },
            errors: {
                404: `Institution does not exist`,
            },
        });
    }
    /**
     * Update institution
     * @param code Institution code
     * @param data Institution informations
     * @returns Institution OK
     * @throws ApiError
     */
    public static updateInstitution(
        code: string,
        data: InstitutionUpdate,
    ): CancelablePromise<Institution> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/people/institutions/{code}',
            path: {
                'code': code,
            },
            body: data,
            errors: {
                400: `Bad Request`,
            },
        });
    }
    /**
     * List persons
     * @returns Person OK
     * @throws ApiError
     */
    public static getPeoplePersons(): CancelablePromise<Array<Person>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/people/persons',
        });
    }
    /**
     * Create person
     * Register a new person
     * @param data Created person
     * @returns Person Created
     * @throws ApiError
     */
    public static createperson(
        data: PersonInput,
    ): CancelablePromise<Person> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/people/persons',
            body: data,
            errors: {
                400: `Bad Request`,
            },
        });
    }
    /**
     * Delete person
     * @param id Item UUID
     * @returns Person OK
     * @throws ApiError
     */
    public static deletePerson(
        id: string,
    ): CancelablePromise<Person> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/people/persons/{id}',
            path: {
                'id': id,
            },
            errors: {
                404: `person does not exist`,
            },
        });
    }
    /**
     * Update person
     * @param id Item UUID
     * @param data Update infos
     * @returns Person OK
     * @throws ApiError
     */
    public static updatePerson(
        id: string,
        data: PersonUpdate,
    ): CancelablePromise<Person> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/people/persons/{id}',
            path: {
                'id': id,
            },
            body: data,
            errors: {
                400: `Bad Request`,
            },
        });
    }
}
