import { AxiosPromise, AxiosRequestConfig } from "axios";
import client from "config";

export default function deleteSecret(token: string, config: AxiosRequestConfig = {}): AxiosPromise {
	return client.delete(`/secrets/${token}`, config);
}
