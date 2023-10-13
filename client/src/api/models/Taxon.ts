/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { TaxonRank } from './TaxonRank';
import type { TaxonStatus } from './TaxonStatus';

export type Taxon = {
    authorship?: string;
    code?: string;
    gbif_ID?: number;
    name: string;
    rank: TaxonRank;
    status: TaxonStatus;
};

