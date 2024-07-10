import { TaxonStatus } from "@/api";
import { stringUnionToArray } from "../toolkit/enums";

export const taxonStatusOptions = stringUnionToArray<TaxonStatus>()('Accepted', "Unreferenced", "Unclassified")




