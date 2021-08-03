import React from "react";
import { Switch, Route } from "react-router-dom";
import { Example } from "./pages";

const AppRouter: React.FC = () => {
	return (
		<Switch>
			{/* Main routing component */}
			<Route exact path="/" component={Example} />
		</Switch>
	);
};

export default AppRouter;
