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
			return "https://whisper.rotational.dev/v1";
		case "development":
			return "http://localhost:3000";
		default:
			throw new Error("Could not identify the api prefix");
	}
}

function encodeFileToBase64(file: File): Promise<string | ArrayBuffer | null> {
	return new Promise((resolve, reject) => {
		const reader = new FileReader();
		reader.readAsDataURL(file);
		reader.onload = () => {
			resolve(reader.result);
		};
		reader.onerror = error => reject(error);
	});
}

function dataURLtoFile(dataurl: string, filename?: string): File {
	const fileName = filename || "";
	const arr = dataurl.split(",");
	const mime = (arr[0].match(/:(.*?);/) as string[])[1];
	const bstr = atob(arr[1]);

	let n = bstr.length;
	const u8arr = new Uint8Array(n);

	while (n--) {
		u8arr[n] = bstr.charCodeAt(n);
	}

	return new File([u8arr], fileName, { type: mime });
}

function generateSecretLink(token: string): string {
	return process.env.NODE_ENV === "development"
		? `http://localhost:3000/secret/${token}`
		: `${defaultEndpointPrefix()}/secret/${token}`;
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
	encodeFileToBase64
};
