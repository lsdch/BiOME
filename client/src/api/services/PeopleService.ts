/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
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

}
