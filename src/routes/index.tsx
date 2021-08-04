import React from "react";
import { Route, Switch } from "react-router-dom";
import CreateSecret from "pages/CreateSecret";

const AppRouter: React.FC = () => {
	return (
		<Switch>
			<Route path="/" component={CreateSecret} />
		</Switch>
	);
};

export default AppRouter;
