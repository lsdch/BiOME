/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { TaxonRank } from './TaxonRank';
import type { TaxonStatus } from './TaxonStatus';
export type TaxonUpdate = {
    GBIF_ID?: number;
    authorship?: string;
    code: string;
    id: string;
    name: string;
    parent?: string;
    rank: TaxonRank;
    status: TaxonStatus;
};

