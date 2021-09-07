import React from "react";
import AttachmentIcon from "@material-ui/icons/Attachment";
import { Box, makeStyles, Theme, Typography } from "@material-ui/core";
import { formatBytes } from "utils/utils";
import dayjs from "dayjs";
import Button from "./Button";

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		width: "100%",
		padding: theme.spacing(2),
		gap: theme.spacing(4),
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
		height: "300px",
		borderRadius: "5px"
	}
}));

type ShowFileProps = {
	file?: File;
	uploadedAt?: Date;
	onDelete: () => void;
	loading: boolean;
};

const ShowFile: React.FC<ShowFileProps> = ({ file, uploadedAt, onDelete, loading }) => {
	const classes = useStyles();

	const handleDownloadClick = () => {
		if (file) {
			const url = window.URL.createObjectURL(new Blob([file]));
			const link = document.createElement("a");
			link.href = url;
			link.target = "_blank";
			link.setAttribute("download", file.name);
			document.body.appendChild(link);
			link.click();
			link.parentNode && link.parentNode.removeChild(link);
		}
	};

	return (
		<div className={classes.container}>
			<Box display="flex" alignItems="center" justifyContent="center" textAlign="center" height="89%">
				<Box>
					<div>
						<AttachmentIcon fontSize="large" />
					</div>
					<div>
						<Typography variant="h6">{file?.name}</Typography>
					</div>
					<div>
						<Typography variant="caption">
							Uploaded on {dayjs(uploadedAt).format("MMMM D, YYYY")}{" "}
							<span style={{ fontWeight: "bold", fontSize: "large" }}>·</span>{" "}
							{`${file?.size && formatBytes(file?.size)}`}
						</Typography>
					</div>
				</Box>
			</Box>
			<Box display="flex" gridGap="1rem">
				<Button label="Download" variant="contained" color="primary" fullWidth onClick={handleDownloadClick} />
				<Button
					label="Detroy the secret"
					variant="contained"
					style={{ background: "red", color: "#fff" }}
					fullWidth
					onClick={onDelete}
					isLoading={loading}
				/>
			</Box>
		</div>
	);
};

export default ShowFile;
