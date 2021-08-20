import { AxiosPromise, AxiosRequestConfig, AxiosResponse } from "axios";
import client from "config";
import { stringToBase64 } from "utils/utils";

export default function getSecret(token: string, password?: string, config: AxiosRequestConfig = {}): AxiosPromise {
	const encodedPassword = password ? stringToBase64(password) : "";
	if (password) {
		return client
			.get(`/secrets/${token}`, {
				headers: { Authorization: `Bearer ${encodedPassword}`, "Access-Control-Request-Headers": "Authorization" }
			})
			.then((response: AxiosResponse) => {
				window.sessionStorage.setItem("__KEY__", encodedPassword);
				return response;
			});
	}

	return client.get(`/secrets/${token}`, config);
}
