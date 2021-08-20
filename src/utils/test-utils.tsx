import { render as rtlRender } from "@testing-library/react";
// import { BrowserRouter } from "react-router-dom";
import { Route, Router } from "react-router-dom";
// import { render } from "@testing-library/react";
import { createMemoryHistory, MemoryHistory } from "history";

// test utils file
export function render(
	ui: any,
	{ route: any = "/", history = createMemoryHistory({ initialEntries: ["secret/"] }) } = {}
) {
	return {
		...rtlRender(<Router history={history}>{ui}</Router>),
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
): any {
	return {
		...rtlRender(
			<Router history={history}>
				<Route path={path}>{ui}</Route>
			</Router>
		)
	};
}

export * from "@testing-library/react";
