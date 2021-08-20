import axios from "axios";

const CONFIG = {
	API_URL: process.env.REACT_APP_API_BASE_URL
};

const client = axios.create({
	baseURL: CONFIG.API_URL
});

export default client;
