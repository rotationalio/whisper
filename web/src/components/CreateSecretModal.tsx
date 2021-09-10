import { Box, IconButton, makeStyles, Snackbar, Typography } from "@material-ui/core";
import { AxiosError } from "axios";
import { useModal } from "contexts";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import React, { useState } from "react";
import CopyToClipboard from "react-copy-to-clipboard";
import { deleteSecret } from "services";
import LinkIcon from "@material-ui/icons/Link";
import { ModalType } from "utils/enums/modal";
import Button from "./Button";
import Alert, { Color } from "@material-ui/lab/Alert";
import { generateSecretLink } from "utils/utils";
import Modal from "./Modal";
dayjs.extend(relativeTime);

const useStyles = makeStyles(() => ({
	ellipsis: {
		whiteSpace: "nowrap",
		textOverflow: "ellipsis",
		display: "block",
		overflow: "scroll",
		scrollbarWidth: "none",
		msOverflowStyle: "none",
		"&::-webkit-scrollbar": {
			display: "none"
		}
	}
}));

const CreateSecretModal: React.FC = () => {
	const classes = useStyles();
	const { state, dispatch } = useModal();
	const [alert, setAlert] = useState<{ open: boolean; severity?: Color; message?: string }>({
		open: false,
		severity: "success",
		message: ""
	});
	const [isLoading, setIsLoading] = useState(false);
	const [isExpired, setIsExpired] = useState(false);
	const secretLink = typeof state.modalProps?.token === "string" ? generateSecretLink(state.modalProps?.token) : "";
	const handleCopy = () => setAlert({ open: true, message: "Secret link copied" });
	const handleAlertClose = () => setAlert({ open: false });

	React.useEffect(() => {
		if (state.modalProps?.expires) {
			const timeInMs = new Date(state.modalProps?.expires).getTime() - Date.now();
			setTimeout(() => {
				setIsExpired(true);
			}, timeInMs);
		}
	}, [state.modalProps?.expires]);

	const handleDeleteSecret = async () => {
		if (window.confirm("Do you really want to remove the secret message ?")) {
			const password = sessionStorage.getItem("__KEY__");

			if (state.modalProps?.token) {
				setIsLoading(true);
				deleteSecret(state.modalProps?.token, password).then(
					() => {
						setIsLoading(false);
						setAlert({ open: true, message: "Secret message destroyed" });
						dispatch({ type: ModalType.HIDE_MODAL });
					},
					(error: AxiosError) => {
						setIsLoading(false);
						setAlert({ open: true, message: error.message, severity: "error" });
					}
				);
			}
		}
	};

	const handleClose = () => {
		if (window.confirm("Be sure to copy the link before closing the window")) {
			dispatch({ type: ModalType.HIDE_MODAL });
		}
	};

	return (
		<Modal open={state.modalType === "SHOW_MODAL"} onClose={handleClose}>
			<>
				<Box>
					<Typography variant="h5" align="center" gutterBottom>
						Secret created successfully
					</Typography>
					{!isExpired ? (
						<>
							<Typography align="center" style={{ display: isExpired ? "none" : "visible" }} gutterBottom>
								You can find your secret on this link below
							</Typography>
							<Typography align="center" gutterBottom>
								It expires{" "}
								<span style={{ color: "red", fontWeight: "bold" }}>{dayjs(state.modalProps?.expires).fromNow()}</span>
							</Typography>
						</>
					) : (
						<Typography align="center" gutterBottom>
							<span style={{ color: "red" }}>You can no longer use this link because it has expired</span>
						</Typography>
					)}
				</Box>

				<Box
					display="flex"
					alignItems="center"
					paddingLeft="15px"
					justifyContent="space-between"
					bgcolor="#f1f3f4"
					borderRadius=".2rem"
					border={`${isExpired ? "1px solid red" : undefined}`}
				>
					<div className={classes.ellipsis} style={{ userSelect: isExpired ? "none" : undefined }}>
						<Typography>{secretLink}</Typography>
					</div>
					<div>
						<CopyToClipboard text={secretLink} onCopy={handleCopy}>
							<IconButton aria-label="copy" size="medium" title="Copy this link" disabled={isExpired}>
								<LinkIcon />
							</IconButton>
						</CopyToClipboard>
					</div>
				</Box>

				<Box display="flex" gridGap="2rem" justifyContent="space-around">
					<Button
						isLoading={isLoading}
						label="Destroy this secret"
						onClick={handleDeleteSecret}
						color="secondary"
						variant="contained"
						disabled={isLoading || isExpired}
						style={{ minWidth: "200px" }}
					/>
				</Box>
				<Snackbar open={alert.open} autoHideDuration={5000} onClose={handleAlertClose}>
					<Alert onClose={handleAlertClose} severity={alert.severity}>
						{alert.message}
					</Alert>
				</Snackbar>
			</>
		</Modal>
	);
};

export default CreateSecretModal;
