import { AxiosPromise, AxiosRequestConfig, AxiosResponse } from "axios";
import client from "config";
import { Secret } from "utils/interfaces/Secret";
import { stringToBase64 } from "utils/utils";

function createSecret(data: Secret, config: AxiosRequestConfig = {}): AxiosPromise {
	return client.post("/secrets", data, config);
}

function deleteSecret(token: string, password?: string | null, config: AxiosRequestConfig = {}): AxiosPromise {
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
function getSecret(token: string, password?: string, config: AxiosRequestConfig = {}): AxiosPromise {
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

export { getSecret, deleteSecret, createSecret };
