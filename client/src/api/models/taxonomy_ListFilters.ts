/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { edgedb_OptionalBool } from './edgedb_OptionalBool';
import type { TaxonRank } from './TaxonRank';
import type { TaxonStatus } from './TaxonStatus';
export type taxonomy_ListFilters = {
    anchors_only?: edgedb_OptionalBool;
    pattern?: string;
    rank?: TaxonRank;
    status?: TaxonStatus;
};

