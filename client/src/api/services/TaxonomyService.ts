/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaxonDB } from '../models/TaxonDB';
import type { TaxonInput } from '../models/TaxonInput';
import type { TaxonWithRelatives } from '../models/TaxonWithRelatives';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class TaxonomyService {

    /**
     * List taxa
     * Lists taxa, optionally filtered by name, rank and status
     * @param pattern Name search pattern
     * @param rank Taxonomic rank
     * @param status Taxonomic status
     * @returns TaxonWithRelatives Get taxon success
     * @throws ApiError
     */
    public static taxonomyList(
        pattern?: string,
        rank?: 'Kingdom' | 'Phylum' | 'Class' | 'Family' | 'Genus' | 'Species' | 'Subspecies',
        status?: 'Accepted' | 'Synonym' | 'Unclassified',
    ): CancelablePromise<Array<TaxonWithRelatives>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/taxonomy/',
            query: {
                'pattern': pattern,
                'rank': rank,
                'status': status,
            },
        });
    }

    /**
     * List anchor taxa
     * Anchors are taxa that were imported as the root of a subtree in the taxonomy.
     * @returns TaxonDB Get anchor taxa list success
     * @throws ApiError
     */
    public static taxonAnchors(): CancelablePromise<Array<TaxonDB>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/taxonomy/anchors',
        });
    }

    /**
     * Import GBIF clade
     * Imports a clade from the GBIF taxonomy, using a its GBIF ID
     * @param code GBIF taxon code
     * @returns TaxonWithRelatives OK
     * @throws ApiError
     */
    public static importGbif(
        code: number,
    ): CancelablePromise<TaxonWithRelatives> {
        return __request(OpenAPI, {
            method: 'PUT',
            url: '/taxonomy/import',
            query: {
                'code': code,
            },
            errors: {
                400: `Bad Request`,
                403: `Forbidden`,
            },
        });
    }

    /**
     * Get a taxon by its code
     * @param code Taxon code
     * @returns TaxonWithRelatives Get taxon success
     * @throws ApiError
     */
    public static getTaxon(
        code: string,
    ): CancelablePromise<TaxonWithRelatives> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/taxonomy/{code}',
            path: {
                'code': code,
            },
            errors: {
                404: `Not Found`,
            },
        });
    }

    /**
     * Delete a taxon by its code
     * @param code Taxon code
     * @returns any OK
     * @throws ApiError
     */
    public static deleteTaxon(
        code: string,
    ): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'DELETE',
            url: '/taxonomy/{code}',
            path: {
                'code': code,
            },
            errors: {
                403: `Forbidden`,
                404: `Not Found`,
            },
        });
    }

    /**
     * Update a taxon by its code
     * @param code Taxon code
     * @param data Taxon
     * @returns TaxonWithRelatives OK
     * @throws ApiError
     */
    public static updateTaxon(
        code: string,
        data: TaxonInput,
    ): CancelablePromise<TaxonWithRelatives> {
        return __request(OpenAPI, {
            method: 'PATCH',
            url: '/taxonomy/{code}',
            path: {
                'code': code,
            },
            body: data,
            errors: {
                403: `Forbidden`,
                404: `Not Found`,
            },
        });
    }

}
