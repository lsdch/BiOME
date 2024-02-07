/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Country } from '../models/Country';
import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';
export class LocationService {
    /**
     * List Countries
     * @returns Country OK
     * @throws ApiError
     */
    public static getCountries(): CancelablePromise<Array<Country>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/countries/',
        });
    }
}
