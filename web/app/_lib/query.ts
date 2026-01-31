import createFetchClient from "openapi-fetch";
import createClient from "openapi-react-query";
import type { paths } from "./api.d";

const apiUrl = process.env.NEXT_PUBLIC_APP_API_URL || "http://localhost:4100";

export const fetchClient = createFetchClient<paths>({ baseUrl: apiUrl });
export const $api = createClient(fetchClient);
