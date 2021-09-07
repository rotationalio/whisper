import { render as rtlRender } from "@testing-library/react";
import { Route, Router } from "react-router-dom";
import { createMemoryHistory, MemoryHistory } from "history";
import { ServerStatusProvider } from "contexts/serverStatusContext";
import { ThemeProvider } from "@material-ui/core";
import theme from "theme";
import { ModalProvider } from "contexts";

export function render(
	ui: React.ReactNode,
	{ history = createMemoryHistory({ initialEntries: ["secret/"] }) } = {}
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
): Record<string, any> {
	return {
		...rtlRender(
			<ThemeProvider theme={theme}>
				<ModalProvider>
					<Router history={history}>{ui}</Router>
				</ModalProvider>
			</ThemeProvider>
		),
		history
	};
}

interface RenderWithRouterProps {
	route?: string;
	history?: MemoryHistory;
	path?: string;
}

export function renderWithRouterMatch(
	ui: React.ReactNode,
	{ path = "/", route = "/", history = createMemoryHistory({ initialEntries: [route] }) }: RenderWithRouterProps = {}
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
): any {
	return {
		...rtlRender(
			<ServerStatusProvider>
				<ModalProvider>
					<Router history={history}>
						<Route path={path}>{ui}</Route>
					</Router>
				</ModalProvider>
			</ServerStatusProvider>
		)
	};
}

export * from "@testing-library/react";
