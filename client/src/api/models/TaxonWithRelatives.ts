/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Meta } from './Meta';
import type { TaxonDB } from './TaxonDB';
import type { TaxonRank } from './TaxonRank';
import type { TaxonStatus } from './TaxonStatus';

export type TaxonWithRelatives = {
    GBIF_ID?: number;
    anchor?: boolean;
    authorship?: string;
    children?: Array<TaxonDB>;
    code: string;
    id: string;
    meta: Meta;
    name: string;
    parent?: {
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
    rank: TaxonRank;
    status: TaxonStatus;
};

