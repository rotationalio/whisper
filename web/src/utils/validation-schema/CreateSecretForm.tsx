import * as Yup from "yup";

const CreateSecretMessageSchema = Yup.object().shape({
	secret: Yup.string().required("The secret message is required"),
	password: Yup.string(),
	lifetime: Yup.object().shape({ value: Yup.string(), label: Yup.string() }).nullable(),
	accessType: Yup.boolean(),
	accesses: Yup.number().max(108, "The max number is 108"),
	file: Yup.mixed(),
	filename: Yup.string()
});

const CreateSecretFileSchema = Yup.object().shape({
	secret: Yup.string(),
	password: Yup.string(),
	lifetime: Yup.object().shape({ value: Yup.string(), label: Yup.string() }).nullable(),
	accessType: Yup.boolean(),
	accesses: Yup.number().max(108, "The max number is 108"),
	file: Yup.mixed().required("The secret file is required"),
	filename: Yup.string()
});

export { CreateSecretFileSchema, CreateSecretMessageSchema };
