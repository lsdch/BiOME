import { ApiError } from "..";

export function handleCommonErrors(err: ApiError) {
  switch (err.status) {
    case 500:
      return "An unhandled error occurred in the server."
    default:
      break;
  }
}