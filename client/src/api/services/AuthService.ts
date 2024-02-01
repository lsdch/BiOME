/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { EmailInput } from '../models/EmailInput';
import type { PasswordInput } from '../models/PasswordInput';
import type { TokenResponse } from '../models/TokenResponse';
import type { UserCredentials } from '../models/UserCredentials';
import type { UserInput } from '../models/UserInput';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class AuthService {

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
                400: `Authentication failure`,
                500: `Internal Server Error`,
            },
        });
    }

    /**
     * Logout user
     * Log out currently authenticated user
     * @returns any User logged out
     * @throws ApiError
     */
    public static logout(): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/logout',
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
                400: `Bad Request`,
                500: `Server error`,
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
                400: `Bad Request`,
            },
        });
    }

    /**
     * Request a password reset token
     * A token to reset the password associated to the provided email address is sent, unless the address is not known in the DB.
     * @param email The email address the account was registered with
     * @returns any Email address is valid and a password reset token was sent
     * @throws ApiError
     */
    public static requestPasswordReset(
        email: EmailInput,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users/forgotten-password',
            body: email,
            errors: {
                400: `Invalid email address`,
            },
        });
    }

    /**
     * Set account password
     * Sets a new password for the currently authenticated user
     * @param password New password
     * @returns any New password was set
     * @throws ApiError
     */
    public static setPassword(
        password: PasswordInput,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/users/password',
            body: password,
            errors: {
                400: `Invalid password inputs`,
                403: `Not authenticated`,
                500: `Database or server error`,
            },
        });
    }

    /**
     * Verify that a password token is valid
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
                400: `Invalid or expired password reset token`,
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
                400: `Bad Request`,
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
