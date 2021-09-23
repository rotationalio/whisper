import axios from "axios";
import { defaultEndpointPrefix } from "../utils/utils";

const CONFIG = {
	API_URL: defaultEndpointPrefix()
};

const client = axios.create({
	baseURL: CONFIG.API_URL
});

export default client;
