import React from "react";
import AttachmentIcon from "@material-ui/icons/Attachment";
import { Box, Button, makeStyles, Theme, Typography } from "@material-ui/core";
import { formatBytes } from "utils/utils";
import dayjs from "dayjs";

const useStyles = makeStyles((theme: Theme) => ({
	container: {
		width: "100%",
		maxWidth: "500px",
		padding: theme.spacing(2),
		backgroundColor: theme.palette.background.paper,
		gap: theme.spacing(4),
		flexDirection: "column",
		justifyContent: "center",
		alignItems: "center",
		height: "300px",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#f5f7f8",
		borderRadius: "5px"
	}
}));

type ShowFileProps = {
	file?: File;
	uploadedAt?: Date;
};

const ShowFile: React.FC<ShowFileProps> = ({ file, uploadedAt }) => {
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
			<div>
				<Button variant="contained" color="primary" onClick={handleDownloadClick} fullWidth>
					Download
				</Button>
			</div>
		</div>
	);
};

export default ShowFile;
