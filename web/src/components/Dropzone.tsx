import { Typography } from "@material-ui/core";
import clsx from "clsx";
import { makeStyles } from "@material-ui/styles";
import { ErrorMessage, FieldProps } from "formik";
import { useMemo } from "react";
import { useDropzone } from "react-dropzone";

const useStyles = makeStyles(() => ({
	error: {
		borderColor: "#f44336 !important"
	}
}));

const baseStyle = {
	flex: 1,
	display: "flex",
	alignItems: "center",
	padding: "20px",
	borderWidth: 2,
	borderRadius: 2,
	height: "175px",
	borderColor: "#eeeeee",
	borderStyle: "dashed",
	backgroundColor: "#fafafa",
	color: "#bdbdbd",
	outline: "none",
	transition: "border .24s ease-in-out"
};

const activeStyle = {
	borderColor: "#2196f3"
};

const acceptStyle = {
	borderColor: "#00e676"
};

type StyledDropzoneProps = FieldProps;

const FILE_SIZE = 64 * 1024;

const StyledDropzone: React.FC<StyledDropzoneProps> = ({ form }) => {
	const classes = useStyles();
	const { getRootProps, getInputProps, acceptedFiles, isDragActive, isDragAccept } = useDropzone({
		maxSize: FILE_SIZE,
		onDrop: (_acceptedFiles: File[]) => {
			const file = _acceptedFiles[0];
			form.setFieldValue("file", file);
			form.setFieldValue("filename", file.name);
			form.setFieldValue("is_base64", true);
		}
	});

	const files = acceptedFiles.map(file => <Typography key={file.name}>{file.name}</Typography>);

	const style = useMemo(
		() => ({
			...baseStyle,
			...(isDragActive ? activeStyle : {}),
			...(isDragAccept ? acceptStyle : {})
		}),
		[isDragActive, isDragAccept]
	);

	return (
		<div className="container">
			<div
				{...getRootProps({ style })}
				className={clsx({ [classes.error]: form.touched["file"] && !!form.errors["file"] })}
			>
				<input {...getInputProps()} />
				{isDragActive ? (
					<p>Drop the files here ...</p>
				) : (
					<>{!files.length ? "Drag'n drop your file here, or click to select the file" : files}</>
				)}
			</div>
			<small className="MuiFormHelperText-root MuiFormHelperText-contained Mui-error Mui-required">
				<ErrorMessage name="file" />
			</small>
		</div>
	);
};

export default StyledDropzone;
