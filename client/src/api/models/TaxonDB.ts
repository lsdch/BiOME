/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Meta } from './Meta';
import type { TaxonRank } from './TaxonRank';
import type { TaxonStatus } from './TaxonStatus';

export type TaxonDB = {
    GBIF_ID?: number;
    anchor?: boolean;
    authorship?: string;
    code: string;
    id: string;
    meta: Meta;
    name: string;
    rank: TaxonRank;
    status: TaxonStatus;
};

