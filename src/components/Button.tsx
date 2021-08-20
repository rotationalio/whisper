import { Button as MuiButton, CircularProgress, ButtonProps } from "@material-ui/core";

interface _ButtonProps extends ButtonProps {
	isLoading?: boolean;
	label: string;
}

const Button: React.FC<_ButtonProps> = ({ isLoading, children, label, ...props }) => {
	return (
		<>
			<MuiButton disabled={isLoading} {...props} data-testid="custom-button">
				{children ? children : isLoading ? <CircularProgress size={24} /> : label}
			</MuiButton>
		</>
	);
};

export default Button;
