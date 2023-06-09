import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { ThemeProvider } from "@material-ui/core";
import GlobalStyles from "./styles/globalStyles";
import theme from "theme";
import initSentry from "./sentry";

initSentry();

ReactDOM.render(
	<React.StrictMode>
		<ThemeProvider theme={theme}>
			<GlobalStyles />
			<App />
		</ThemeProvider>
	</React.StrictMode>,
	document.getElementById("root")
);
