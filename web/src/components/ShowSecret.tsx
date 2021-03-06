import React from "react";
import { Box, makeStyles, Theme, Typography } from "@material-ui/core";
import { Alert, AlertTitle } from "@material-ui/lab";
import { Secret } from "utils/interfaces/Secret";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { dataURLtoFile, selectOnFocus } from "utils/utils";
import { Link, useHistory } from "react-router-dom";
import { deleteSecret } from "services";
import clsx from "clsx";
import Button from "./Button";
import ShowFile from "./ShowFile";
import { ContentLayout } from "layouts";
dayjs.extend(relativeTime);

type ShowSecretProps = {
	secret?: Secret;
	token: string;
};

const useStyles = makeStyles((theme: Theme) => ({
	section: {
		display: "flex",
		flexDirection: "column",
		margin: "0 auto",
		justifyContent: "center",
		alignItems: "center",
		outline: 0,
		width: "100%"
	},
	box: {
		width: "100%"
	},
	textarea: {
		width: "100%",
		padding: theme.spacing(2),
		minHeight: "15rem",
		border: "1px solid rgba(0,0,0,.2)",
		boxSizing: "border-box",
		lineHeight: 1.75,
		outline: "none",
		borderRadius: 5
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

	const handleDeleteClick = () => {
		const encodedPassword = window.sessionStorage.getItem("__KEY__");

		if (window.confirm("Do you really want to destroy the secret?")) {
			setIsLoading(true);
			deleteSecret(token, encodedPassword).then(
				() => {
					setIsLoading(false);
					history.push("/");
				},
				async () => {
					setIsLoading(false);
				}
			);
		}
	};

	return (
		<ContentLayout>
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
		</ContentLayout>
	);
};

export default ShowSecret;
