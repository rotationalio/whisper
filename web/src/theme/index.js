import { createTheme, responsiveFontSizes } from "@material-ui/core";
import { green, indigo } from "@material-ui/core/colors";

const theme = createTheme({
	palette: {
		primary: indigo,
		secondary: green
	}
});

export default responsiveFontSizes(theme);
