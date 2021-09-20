/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import { Form } from "formik";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
function preventNonNumericalInput(e: React.KeyboardEvent<typeof Form>) {
	e = e || window.event;
	const charCode = typeof e.which == "undefined" ? e.keyCode : e.which;
	const charStr = String.fromCharCode(charCode);

	if (!charStr.match(/^[0-9]+$/)) e.preventDefault();
}

function defaultEndpointPrefix(): string {
	const baseUrl = process.env.REACT_APP_API_BASE_URL;
	if (baseUrl) {
		return baseUrl;
	}

	switch (process.env.NODE_ENV) {
		case "production":
			return "https://api.whisper.rotational.dev/v1";
		case "development":
			return "http://localhost:8318/v1";
		default:
			throw new Error("Could not identify the api prefix");
	}
}

function defaultAbsoluteURL(): string {
	const baseURL = process.env.REACT_APP_UI_BASE_URL;
	if (baseURL) {
		return baseURL;
	}

	switch (process.env.NODE_ENV) {
		case "production":
			return "https://whisper.rotational.dev";
		case "development":
			return "http://localhost:3000";
		default:
			throw new Error("could not identify the ui absolute url");
	}
}

function encodeFileToBase64(file: File): Promise<string | ArrayBuffer | null> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.readAsDataURL(file);
		reader.onload = () => {
			const base64 = typeof reader.result === "string" ? reader.result.split(",")[1] : null;
			resolve(base64);
		};
		reader.onerror = error => reject(error);
	});
}

function dataURLtoFile(dataurl: string, filename?: string): File {
	const fileName = filename || "";
	const bstr = atob(dataurl);

	let n = bstr.length;
	const u8arr = new Uint8Array(n);

	while (n--) {
		u8arr[n] = bstr.charCodeAt(n);
	}

	return new File([u8arr], fileName);
}

function formatBytes(bytes: number, decimals = 2) {
	if (bytes === 0) return "0 Bytes";

	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

	const i = Math.floor(Math.log(bytes) / Math.log(k));

	return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}

function generateSecretLink(token: string): string {
	return `${defaultAbsoluteURL()}/secret/${token}`;
}

function stringToBase64(str: string): string {
	return Buffer.from(str).toString("base64");
}

function selectOnFocus(e: React.FocusEvent<HTMLTextAreaElement>): void {
	e.target.select();
}

export {
	generateSecretLink,
	defaultEndpointPrefix,
	preventNonNumericalInput,
	stringToBase64,
	selectOnFocus,
	dataURLtoFile,
	encodeFileToBase64,
	formatBytes
};
