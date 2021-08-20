import { AxiosPromise, AxiosRequestConfig } from "axios";
import client from "config";
import { Secret } from "utils/interfaces/Secret";

export default function createSecret(data: Secret, config: AxiosRequestConfig = {}): AxiosPromise {
	return client.post("/secrets", data, config);
}
