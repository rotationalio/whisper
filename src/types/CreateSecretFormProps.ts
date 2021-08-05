import { FormikHelpers, FormikValues } from "formik";
import { ObjectSchema } from "yup";

export type CreateSecretFormProps = {
	onSubmit: (values: FormikValues, helpers: FormikHelpers<any>) => void;
	validationSchema: ObjectSchema<any>;
	initialValues: FormikValues;
};
