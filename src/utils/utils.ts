import { Form } from "formik";

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
export function preventNonNumericalInput(e: React.KeyboardEvent<typeof Form>) {
	e = e || window.event;
	const charCode = typeof e.which == "undefined" ? e.keyCode : e.which;
	const charStr = String.fromCharCode(charCode);

	if (!charStr.match(/^[0-9]+$/)) e.preventDefault();
}
