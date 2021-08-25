import React from "react";
import { Box, makeStyles, Theme, Typography } from "@material-ui/core";
import { Alert, AlertTitle } from "@material-ui/lab";
import { Secret } from "utils/interfaces/Secret";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { dataURLtoFile, selectOnFocus } from "utils/utils";
import { Link, useHistory } from "react-router-dom";
import deleteSecret from "services/deleteSecret";
import clsx from "clsx";
import Button from "./Button";
import ShowFile from "./ShowFile";
dayjs.extend(relativeTime);

type ShowSecretProps = {
	secret?: Secret;
	token: string;
};

const useStyles = makeStyles((theme: Theme) => ({
	section: {
		height: "100vh",
		display: "flex",
		flexDirection: "column",
		margin: "0 auto",
		justifyContent: "center",
		alignItems: "center",
		outline: 0,
		maxWidth: "500px",
		width: "100%"
	},
	box: {
		width: "100%",
		boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
		background: "#f5f7f8",
		padding: theme.spacing(2),
		borderRadius: "5px"
	},
	textarea: {
		width: "100%",
		padding: theme.spacing(2),
		minHeight: "15rem",
		border: "1px solid rgba(0,0,0,.2)",
		boxSizing: "border-box",
		lineHeight: 1.75
	},
	hide: {
		display: "none"
	},
	link: {
		textDecoration: "none",
		maxWidth: "50%",
		width: "100%"
	},
	fullWidth: {
		width: "100%"
	},
	deleteButton: {
		flexGrow: 1,
		background: "red",
		color: "#fff"
	}
}));

const ShowSecret: React.FC<ShowSecretProps> = ({ secret, token }) => {
	const [isLoading, setIsLoading] = React.useState(false);
	const [file, setFile] = React.useState<File | undefined>(undefined);
	const history = useHistory();
	const classes = useStyles();

	React.useEffect(() => {
		if (secret?.is_base64) {
			const _file = dataURLtoFile(secret.secret, secret.filename);
			setFile(_file);
		}
	}, []);

	const deleteWithPassword = (password: string) => {
		deleteSecret(token, {
			headers: {
				Authorization: `Bearer ${password}`,
				"Access-Control-Request-Headers": "Authorization"
			}
		}).then(() => {
			window.sessionStorage.removeItem("__KEY__");
			setIsLoading(false);
			history.push("/");
		});
	};

	const deleteWithoutPassword = () => {
		deleteSecret(token).then(
			() => {
				setIsLoading(false);
				history.push("/");
			},
			async () => {
				setIsLoading(false);
			}
		);
	};

	const handleDeleteClick = () => {
		const encodedPassword = window.sessionStorage.getItem("__KEY__") || null;

		if (window.confirm("Do you really want to destroy the secret?")) {
			setIsLoading(true);
			if (encodedPassword) {
				deleteWithPassword(encodedPassword);
			} else {
				deleteWithoutPassword();
			}
		}
	};

	return (
		<div className={classes.section}>
			<Alert
				severity="warning"
				style={{
					margin: "1rem 0",
					width: "100%",
					display: secret?.destroyed ? undefined : "none"
				}}
			>
				<AlertTitle>Secret Expired</AlertTitle>
				<Typography>
					This is the last time you will be able to access this Secret, it has been destroyed now that you&apos;ve
					retrieved it.
				</Typography>
			</Alert>
			{secret?.is_base64 ? (
				<ShowFile file={file} uploadedAt={secret.created} loading={isLoading} onDelete={handleDeleteClick} />
			) : (
				<div className={classes.box}>
					<Box marginBottom="2rem">
						<Typography variant="h5" gutterBottom>
							Secret
						</Typography>
						<div>
							<textarea
								className={classes.textarea}
								onFocus={selectOnFocus}
								readOnly
								autoFocus
								defaultValue={secret?.secret}
								aria-label="secret-message"
							/>
						</div>
						<Typography variant="caption" gutterBottom>
							Created: {dayjs(secret?.created).fromNow()}
						</Typography>
					</Box>
					<Box display="flex" justifyContent="space-between" gridGap="1rem" flexWrap="wrap">
						<Link to="/" className={clsx({ [classes.fullWidth]: secret?.destroyed }, classes.link)}>
							<Button label="Create another Secret" variant="contained" fullWidth color="primary" />
						</Link>
						<Button
							label="Destroy this secret"
							variant="contained"
							isLoading={isLoading}
							onClick={handleDeleteClick}
							disabled={isLoading}
							style={{ background: "red", color: "#fff", maxWidth: "50%" }}
							className={clsx({ [classes.hide]: secret?.destroyed }, classes.deleteButton)}
						/>
					</Box>
				</div>
			)}
		</div>
	);
};

export default ShowSecret;
