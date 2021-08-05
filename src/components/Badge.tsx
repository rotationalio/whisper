import { makeStyles, Typography } from "@material-ui/core";

const useStyles = makeStyles({
	root: {
		height: "0.5rem",
		width: "0.5rem",
		display: "inline-block",
		background: (props: any) => props.color,
		borderRadius: "50%",
		marginInline: ".3rem"
	}
});

type BadgeProps = {
	color: string;
	content: string;
};

const Badge: React.FC<BadgeProps> = props => {
	const classes = useStyles(props);

	return (
		<Typography variant="caption">
			{props.content}
			<span color={props.color} className={classes.root}></span>
		</Typography>
	);
};

export default Badge;
