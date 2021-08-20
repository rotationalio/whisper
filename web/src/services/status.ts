import { AxiosPromise } from "axios";
import client from "config";

export default function getStatus(config = {}): AxiosPromise {
	return client.get("/status", config);
}
