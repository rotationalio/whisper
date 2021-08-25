import React from "react";
import AttachmentIcon from "@material-ui/icons/Attachment";
import { Box, Button, makeStyles, Theme, Typography } from "@material-ui/core";

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		position: "absolute",
		top: "50%",
		left: "50%",
		transform: "translate(-50%, -50%)",
		width: "100%",
		maxWidth: "500px",
		background: "#f7f5f5",
		padding: theme.spacing(2)
	}
}));

type ShowFileProps = {
	file?: File;
	uploadedAt?: Date;
};

const ShowFile: React.FC<ShowFileProps> = ({ file }) => {
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
			<Box display="flex" alignItems="center" justifyContent="space-between">
				<div>
					<AttachmentIcon fontSize="large" />
				</div>
				<div>
					<Typography variant="caption">Name</Typography>
					<Typography>{file?.name}</Typography>
				</div>
				<div>
					<Typography variant="caption">Size</Typography>
					<Typography>{`${file?.size && (file?.size / 1024).toFixed(2)} Mb`}</Typography>
				</div>
				<div>
					<Button variant="contained" color="primary" onClick={handleDownloadClick} fullWidth>
						Download
					</Button>
				</div>
			</Box>
		</div>
	);
};

export default ShowFile;
