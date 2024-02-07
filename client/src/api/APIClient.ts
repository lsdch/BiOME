/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { BaseHttpRequest } from './core/BaseHttpRequest';
import type { OpenAPIConfig } from './core/OpenAPI';
import { FetchHttpRequest } from './core/FetchHttpRequest';

import { AuthService } from './services/AuthService';
import { LocationService } from './services/LocationService';
import { PeopleService } from './services/PeopleService';
import { TaxonomyService } from './services/TaxonomyService';

type HttpRequestConstructor = new (config: OpenAPIConfig) => BaseHttpRequest;

export class APIClient {

    public readonly auth: AuthService;
    public readonly location: LocationService;
    public readonly people: PeopleService;
    public readonly taxonomy: TaxonomyService;

    public readonly request: BaseHttpRequest;

    constructor(config?: Partial<OpenAPIConfig>, HttpRequest: HttpRequestConstructor = FetchHttpRequest) {
        this.request = new HttpRequest({
            BASE: config?.BASE ?? '/api/v1',
            VERSION: config?.VERSION ?? '1.0',
            WITH_CREDENTIALS: config?.WITH_CREDENTIALS ?? false,
            CREDENTIALS: config?.CREDENTIALS ?? 'include',
            TOKEN: config?.TOKEN,
            USERNAME: config?.USERNAME,
            PASSWORD: config?.PASSWORD,
            HEADERS: config?.HEADERS,
            ENCODE_PATH: config?.ENCODE_PATH,
        });

        this.auth = new AuthService(this.request);
        this.location = new LocationService(this.request);
        this.people = new PeopleService(this.request);
        this.taxonomy = new TaxonomyService(this.request);
    }
}

