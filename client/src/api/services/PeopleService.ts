/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { PasswordInput } from '../models/PasswordInput';
import type { TokenResponse } from '../models/TokenResponse';
import type { User } from '../models/User';
import type { UserCredentials } from '../models/UserCredentials';
import type { UserInput } from '../models/UserInput';

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
     * Authenticate user
     * Authenticate user with their credentials and set a JWT.
     * @param data User credentials
     * @returns TokenResponse Returns a token and stores it as a session cookie
     * @throws ApiError
     */
    public static login(
        data: UserCredentials,
    ): CancelablePromise<TokenResponse> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/login',
            body: data,
            errors: {
                400: `Invalid credentials`,
            },
        });
    }

    /**
     * Email confirmation
     * Confirms a user email using a token
     * @param token Confirmation token
     * @returns any Email was confirmed and account activated
     * @throws ApiError
     */
    public static emailConfirmation(
        token: string,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/users/confirm',
            query: {
                'token': token,
            },
            errors: {
                400: `Invalid or expired confirmation token`,
                500: `Token parse error`,
            },
        });
    }

    /**
     * Resend confirmation email
     * Send again the confirmation email
     * @param data User informations
     * @returns any Email was sent
     * @throws ApiError
     */
    public static resendConfirmationEmail(
        data: UserCredentials,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users/confirm/resend',
            body: data,
            errors: {
                400: `Invalid parameters`,
            },
        });
    }

    /**
     * Verify a password token is valid
     * @param token Password reset token
     * @returns any Password token is valid
     * @throws ApiError
     */
    public static validatePasswordToken(
        token: string,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/users/password-reset/{token}',
            path: {
                'token': token,
            },
            errors: {
                400: `Invalid or expired confirmation token, or invalid input password`,
            },
        });
    }

    /**
     * Reset account password
     * Resets a user's password using a token sent to their email address.
     * @param token Password reset token
     * @param password New password
     * @returns any Password was reset successfully
     * @throws ApiError
     */
    public static resetPassword(
        token: string,
        password: PasswordInput,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users/password-reset/{token}',
            path: {
                'token': token,
            },
            body: password,
            errors: {
                400: `Invalid or expired confirmation token, or invalid input password`,
                500: `Database error`,
            },
        });
    }

    /**
     * Register user
     * Register a new user account, that is inactive (until email is verified or admin intervention), and has role 'Guest'
     * @param data User informations
     * @returns any User created and waiting for email verification
     * @throws ApiError
     */
    public static registerUser(
        data: UserInput,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users/register',
            body: data,
            errors: {
                400: `Invalid parameters`,
            },
        });
    }

    /**
     * Delete a user
     * Deletes a user
     * @returns any User was deleted successfully
     * @throws ApiError
     */
    public static deleteUsers(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/users/{uuid}',
            errors: {
                401: `Admin privileges required`,
                404: `User does not exist`,
            },
        });
    }

}
