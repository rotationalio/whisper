import { AxiosPromise, AxiosRequestConfig } from "axios";
import client from "config";

export default function deleteSecret(
	token: string,
	password?: string | null,
	config: AxiosRequestConfig = {}
): AxiosPromise {
	if (password) {
		return client.delete(`/secrets/${token}`, {
			headers: {
				Authorization: `Bearer ${password}`,
				"Access-Control-Request-Headers": "Authorization"
			}
		});
	}
	return client.delete(`/secrets/${token}`, config);
}
